package app

import (
	"context"
	"fmt"

	"sync"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 最小调度单位
type Task struct {
	downloader *shared.Downloader
	part       *shared.Part
}

// 任务队列
type TaskQueue struct {
	app        *App
	ctx        context.Context // app上下文
	wg         sync.WaitGroup
	mu         sync.Mutex
	tasks      chan *Task
	queueTasks chan *Task
	done       chan struct{}
}

// 创建新的任务队列 并添加上下文
func NewTaskQueue(a *App, tasks []*Task) {
	tq := &TaskQueue{
		app:        a,
		ctx:        a.ctx,
		tasks:      make(chan *Task, a.config.DownloadLimit),
		queueTasks: make(chan *Task, 9999),
		done:       make(chan struct{}),
	}

	tq.AddTasks(tasks)

	limit := utils.ClampInt(a.config.DownloadLimit, 1, 10)
	for i := 0; i < limit; i++ {
		go tq.worker()
	}

	go func() {
		tq.wg.Wait()
		close(tq.done)
		print("关闭tq nil")
		tq = nil
	}()
	print("创建队列成功") // 没有到达
	a.taskQueue = tq
}

// 添加所有任务到队列
func (tq *TaskQueue) addTasksWithoutLock(tasks []*Task) {
	for _, task := range tasks {
		tq.wg.Add(1)
		select {
		case tq.tasks <- task:
		default:
			tq.queueTasks <- task
		}
	}
}

func (tq *TaskQueue) AddTasks(tasks []*Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.addTasksWithoutLock(tasks)
}
func (tq *TaskQueue) worker() {

	for {
		select {
		case task, ok := <-tq.tasks:
			if ok {
				println("开始处理", task.part.Title)
				tq.handleTask(task)
			} else {
				if len(tq.queueTasks) > 0 {
					tq.refillTasks()
				} else {
					// if queueTasks is empty, exit the worker
					return
				}
			}
		case <-tq.done:
			println("taskQue 接收到完成信号")
			return
		}
	}
}

func (tq *TaskQueue) refillTasks() {
	for {
		select {
		case task := <-tq.queueTasks:
			tq.tasks <- task
		default:
			return
		}
	}
}

// 移除tq.tasks通道中的任务, 需要先移除所有的 再填充...
func (tq *TaskQueue) RemoveTasks(tasks []*Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	var tempTasks []*Task

	for {
		select {
		case t := <-tq.tasks:
			if !contains(tasks, t) {
				tempTasks = append(tempTasks, t)
			} else {
				tq.wg.Done() // 移除任务时减少计数
			}
		default:
			goto removeDone
		}
	}

removeDone:
	tq.AddTasks(tempTasks)
}

func contains(tasks []*Task, task *Task) bool {
	for _, t := range tasks {
		if t.part.UID == task.part.UID {
			return true
		}
	}
	return false
}

func (tq *TaskQueue) handleTask(task *Task) {

	defer func() {
		tq.wg.Done()
		task.downloader = nil
	}()

	for i, t := range tq.app.tasks {
		if t == task {
			tq.app.tasks[i] = task
			break
		}
	}

	tq.taskStart(tq.app.Logger, task.part)

	// 获取视频元数据
	if err := (*task.downloader).GetMeta(*tq.app.config, task.part, tq.app.Callback); err != nil {
		tq.app.Logger.Errorf("%s失败: %v", shared.TaskStatus.GettingMetadata, err)
		tq.handleDownloadError(tq.app.Logger, task.part, err)
		updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
		return
	}

	// 默认下封面 又不大!
	if err := (*task.downloader).DownloadThumbnail(task.part, tq.app.Callback); err != nil {
		tq.handleDownloadError(tq.app.Logger, task.part, err)
		return
	}

	if tq.app.config.DownloadVideo && task.part.State != shared.TaskStatus.Stopped {
		task.part.Status = shared.TaskStatus.DownloadingVideo
		if err := (*task.downloader).DownloadVideo(task.part, tq.app.Callback); err != nil {
			tq.handleDownloadError(tq.app.Logger, task.part, err)
			updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
			return
		}
	}

	if tq.app.config.DownloadAudio && task.part.State != shared.TaskStatus.Stopped {
		task.part.Status = shared.TaskStatus.DownloadingAudio
		if err := (*task.downloader).DownloadAudio(task.part, tq.app.Callback); err != nil {
			tq.handleDownloadError(tq.app.Logger, task.part, err)
			updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
			return
		}
	}
	if tq.app.config.DownloadSubtitle && task.part.State != shared.TaskStatus.Stopped {
		task.part.Status = shared.TaskStatus.DownloadingSubtitle
		if err := (*task.downloader).DownloadSubtitle(task.part, tq.app.Callback); err != nil {
			tq.handleDownloadError(tq.app.Logger, task.part, err)
			updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
			return
		}
	}

	// 合并
	if tq.app.config.DownloadCombine && task.part.State != shared.TaskStatus.Stopped {
		task.part.Status = shared.TaskStatus.Merging
		tq.app.Callback(shared.NoticeData{
			EventName: "updateInfo",
			Message:   task.part,
		})
		if err := (*task.downloader).Combine(tq.app.config.FFMPEG, task.part); err != nil {
			tq.handleDownloadError(tq.app.Logger, task.part, err)
			updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)
			return
		}
	}

	// 清理工作
	if err := (*task.downloader).Clear(task.part, tq.app.Callback); err != nil {
		tq.handleDownloadError(tq.app.Logger, task.part, err)
		return
	} else {
		tq.app.Logger.Info("清理工作完成")
	}

	tq.taskFinish(task, tq.app.Logger)
	updateTaskConfig(tq.app.Logger, task, tq.app.tasks, tq.app.configDir)

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

func (tq *TaskQueue) handleDownloadError(logger *logrus.Logger, part *shared.Part, err error) {

	part.Status = fmt.Sprintf("%s: %s", shared.TaskStatus.Failed, err.Error())
	logger.Errorf(shared.TaskStatus.Failed)
	runtime.EventsEmit(tq.ctx, "updateInfo", *part)
}

func (tq *TaskQueue) taskFinish(task *Task, logger *logrus.Logger) {

	task.part.DownloadPercent = 100
	task.part.Status = shared.TaskStatus.Finished
	task.part.State = shared.TaskStatus.Finished
	logger.Infof("%s下载完成", task.part.Title)
	runtime.EventsEmit(tq.ctx, "updateInfo", *task.part)
}
