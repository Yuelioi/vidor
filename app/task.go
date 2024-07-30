package app

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"sync"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Queue = iota
	Working
	Finished
)

// 最小调度单位
type Task struct {
	downloader shared.Downloader
	state      int
	part       *shared.Part
}

// 任务队列
type TaskQueue struct {
	app            *App
	ctx            context.Context // app上下文
	state          int             // Working/Finished
	wg             sync.WaitGroup
	mu             sync.Mutex
	tasksRemaining atomic.Int64
	tasks          chan *Task // 任务通道
	queueTasks     []*Task    // 队列任务
	done           chan struct{}
}

/*
创建新的任务队列
并添加任务
*/
func NewTaskQueue(a *App, tasks []*Task) *TaskQueue {
	limit := utils.ClampInt(a.config.DownloadLimit, 1, 10)

	tq := &TaskQueue{
		app:        a,
		ctx:        a.ctx,
		state:      Working,
		tasks:      make(chan *Task, limit),
		queueTasks: make([]*Task, 0),
		done:       make(chan struct{}),
	}

	tq.AddTasks(tasks)

	for i := 0; i < limit; i++ {
		go tq.worker()
	}

	go func() {
		tq.wg.Wait()
	}()

	return tq
}

// 任务工作函数
// 影响: tq.queueTasks tq.tasks

func (tq *TaskQueue) worker() {
	for {
		select {
		case task, ok := <-tq.tasks:
			if !ok {
				// 队列通道关闭 直接退出
				tq.Close()
				return
			}

			// 创建下载器
			downloader, err := newDownloader(tq.app.downloaders, tq.app.Notice, task.part.Url)
			fmt.Printf("downloader: %v\n", downloader)
			task.downloader = downloader
			if err != nil {
				continue
			}
			tq.state = Working
			task.state = Working
			tq.handleTask(task)

		case <-tq.done:
			tq.Close()
			return

		default:
			// 每秒检测下任务队列
			if len(tq.queueTasks) > 0 {
				tq.reFillTasks()
			} else {
				tq.Close()
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

/*
添加任务
 1. 添加到任务通道
 2. 超额的添加到tq.queueTasks

影响: tq.queueTasks
*/
func (tq *TaskQueue) AddTasks(tasks []*Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.queueTasks = append(tq.queueTasks, tasks...)
}

// 重新填充任务
//
// 影响: tq.queueTasks tq.tasks 已加锁
func (tq *TaskQueue) reFillTasks() {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	availableCapacity := cap(tq.tasks) - len(tq.tasks)
	slice := utils.MinInt(availableCapacity, len(tq.queueTasks))

	for _, task := range tq.queueTasks {
		if availableCapacity <= 0 {
			break
		}

		task.state = Working
		tq.tasks <- task
		// 添加到任务通道 wg+1
		tq.tasksRemaining.Add(1)
		tq.wg.Add(1)
		availableCapacity--
	}
	tq.queueTasks = tq.queueTasks[slice:]
}

// 移除队列中的任务
// 影响: tq.queueTasks 已加锁
func (tq *TaskQueue) removeQueueTasks(tasks []*Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tasksSet := make(map[string]struct{})
	for _, task := range tasks {
		tasksSet[task.part.TaskID] = struct{}{}
	}

	newTasks := make([]*Task, 0)

	// 基于索引重新规划 queueTasks
	for _, task := range tq.queueTasks {
		if _, ok := tasksSet[task.part.TaskID]; ok {
			continue
		}
		newTasks = append(newTasks, task)
	}

	tq.queueTasks = newTasks
}

func (tq *TaskQueue) stopTask(task *Task) {
	fmt.Printf("task.downloader: %v\n", task.downloader)
	task.downloader.StopDownload(task.part, tq.app.Callback)
}

/*
处理任务

	1.正常下载流程
	2.完成需要检测任务队列
*/
func (tq *TaskQueue) handleTask(task *Task) {

	defer func() {
		// if err := recover(); err != nil {
		// 	tq.app.Logger.Errorf("Task handling panic: %v", err)
		// }
	}()

	if task.downloader == nil {
		tq.handleDownloadError(tq.app.Logger, task, fmt.Errorf("downloader is nil"))
		return
	}

	defer func() {
		tq.tasksRemaining.Add(-1)
		// 下载完成 检测任务队列
		if len(tq.queueTasks) > 0 {
			fmt.Printf("任务%s 完成: 准备重新填充\n", task.part.TaskID)
			tq.reFillTasks()
			tq.state = Finished
		} else if tq.tasksRemaining.Load() != 0 {
			fmt.Printf("任务%s 完成, 等待后续任务下载完毕: \n", task.part.TaskID)
			tq.state = Finished

		} else {
			fmt.Printf("任务%s 完成: 关闭下载队列\n\n", task.part.TaskID)
			tq.state = Finished
			close(tq.done)
		}
		task.downloader = nil
		tq.wg.Done()
	}()

	tq.taskStart(tq.app.Logger, task.part)

	// 获取视频元数据
	if task.state == Working {
		if err := task.downloader.GetMeta(*tq.app.config, task.part, tq.app.Callback); err != nil {
			tq.handleDownloadError(tq.app.Logger, task, err)
			updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
			return
		}
	}

	// 默认下封面 又不大!
	if task.state == Working {
		if err := task.downloader.DownloadThumbnail(task.part, tq.app.Callback); err != nil {
			tq.handleDownloadError(tq.app.Logger, task, err)
			return
		}
	}

	if task.state == Working {
		if tq.app.config.DownloadVideo && task.part.State != shared.TaskStatus.Stopped {
			task.part.Status = shared.TaskStatus.DownloadingVideo
			if err := task.downloader.DownloadVideo(task.part, tq.app.Callback); err != nil {
				tq.handleDownloadError(tq.app.Logger, task, err)
				updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
				return
			}
		}
	}
	if task.state == Working {
		if tq.app.config.DownloadAudio && task.part.State != shared.TaskStatus.Stopped {
			task.part.Status = shared.TaskStatus.DownloadingAudio
			if err := task.downloader.DownloadAudio(task.part, tq.app.Callback); err != nil {
				tq.handleDownloadError(tq.app.Logger, task, err)
				updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
				return
			}
		}
	}
	if task.state == Working {
		if tq.app.config.DownloadSubtitle && task.part.State != shared.TaskStatus.Stopped {
			task.part.Status = shared.TaskStatus.DownloadingSubtitle
			if err := task.downloader.DownloadSubtitle(task.part, tq.app.Callback); err != nil {
				tq.handleDownloadError(tq.app.Logger, task, err)
				updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
				return
			}
		}
	}

	// 合并
	if task.state == Working {
		if tq.app.config.DownloadCombine && task.part.State != shared.TaskStatus.Stopped {
			task.part.Status = shared.TaskStatus.Merging
			tq.app.Callback(shared.NoticeData{
				EventName: "updateInfo",
				Message:   task.part,
			})
			if err := task.downloader.Combine(tq.app.config.FFMPEG, task.part); err != nil {
				tq.handleDownloadError(tq.app.Logger, task, err)
				updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
				return
			}
		}
	}

	// 清理工作
	if task.state == Working {
		if err := task.downloader.Clear(task.part, tq.app.Callback); err != nil {
			tq.handleDownloadError(tq.app.Logger, task, err)
			return
		} else {
			tq.app.Logger.Info("清理工作完成")
		}
	}

	tq.taskFinish(task, tq.app.Logger)
	updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)

}

func (tq *TaskQueue) Close() {
	tq.state = Finished
}

func (tq *TaskQueue) taskStart(logger *logrus.Logger, part *shared.Part) {
	logger.Infof(shared.TaskStatus.Downloading)
	part.DownloadPercent = 1
	part.Status = shared.TaskStatus.Downloading
	part.State = shared.TaskStatus.Downloading
}

func updateTaskConfig(logger *logrus.Logger, task *Task, appTasks []*Task, appConfigDir string) error {
	// 更新任务数据
	if err := saveTask(task, appTasks, appConfigDir); err != nil {
		logger.Errorf("保存任务数据失败: %v", err)
		return fmt.Errorf("保存任务数据失败: %v", err)
	} else {
		logger.Infof("%s任务完成", task.part.Title)
	}
	return nil
}

func (tq *TaskQueue) handleDownloadError(logger *logrus.Logger, task *Task, err error) {
	task.state = Finished
	task.part.Status = fmt.Sprintf("%s: %s", shared.TaskStatus.Failed, err.Error())
	logger.Errorf(shared.TaskStatus.Failed, err.Error())
	runtime.EventsEmit(tq.ctx, "updateInfo", task.part)
}

func (tq *TaskQueue) taskFinish(task *Task, logger *logrus.Logger) {
	task.state = Finished
	task.part.DownloadPercent = 100
	task.part.Status = shared.TaskStatus.Finished
	task.part.State = shared.TaskStatus.Finished
	logger.Infof("%s下载完成", task.part.Title)
	runtime.EventsEmit(tq.ctx, "updateInfo", *task.part)
}
