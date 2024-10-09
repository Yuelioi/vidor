package task

import (
	"context"
	"fmt"
	"io"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	pb "github.com/Yuelioi/vidor/internal/proto"
)

// TaskQueue 任务队列接口
type TaskQueue struct {
	manager *plugin.PluginManager

	workingTasks  []*pb.Task
	queueTasks    []*pb.Task
	finishedTasks []*pb.Task
	limit         int
	mu            sync.Mutex
	taskChan      chan *pb.Task
	ctx           context.Context
	queueEnable   atomic.Bool
	notifyEnable  atomic.Bool
	notification  *notify.TaskNotification
}

// New 创建一个新的任务队列
func New(limit int, manager *plugin.PluginManager, ctx context.Context) *TaskQueue {

	taskNotify := notify.NewTaskNotification(ctx)

	return &TaskQueue{
		manager:       manager,
		workingTasks:  make([]*pb.Task, 0),
		queueTasks:    make([]*pb.Task, 0),
		finishedTasks: make([]*pb.Task, 0),
		limit:         limit,
		taskChan:      make(chan *pb.Task, limit),
		ctx:           ctx,
		queueEnable:   atomic.Bool{},
		notifyEnable:  atomic.Bool{},
		notification:  taskNotify,
	}
}
func (tq *TaskQueue) Add(task *pb.Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.queueTasks = append(tq.queueTasks, task)

	if !tq.queueEnable.Load() {
		tq.queueEnable.Store(true)
		go func() {
			tq.Start()
			tq.queueEnable.Store(false)
		}()
	}
	tq.Reload()
}

func (tq *TaskQueue) AddAll(tasks []*pb.Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.queueTasks = append(tq.queueTasks, tasks...)

	fmt.Printf("tq.queueEnable.Load(): %v\n", tq.queueEnable.Load())

	if !tq.queueEnable.Load() {
		tq.queueEnable.Store(true)
		go func() {
			tq.Start()
			tq.queueEnable.Store(false)
		}()
	}
	tq.Reload()
}

func (tq *TaskQueue) StartNotify() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		for _, task := range tq.workingTasks {
			task.State = 1
		}
		for _, task := range tq.queueTasks {
			task.State = 2
		}
		for _, task := range tq.finishedTasks {
			task.State = 3
		}

		tasks := slices.Concat(tq.workingTasks, tq.queueTasks, tq.finishedTasks)

		tq.notification.UpdateTasks(tasks)

		// 没有工作任务 就退出
		if len(tq.queueTasks) == 0 && len(tq.workingTasks) == 0 {
			return
		}
	}
}

func (tq *TaskQueue) Start() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go tq.StartNotify()

	for range ticker.C {
		select {
		case task, ok := <-tq.taskChan:
			if !ok {
				continue
			}

			fmt.Printf("开始下载task.Title: %v\n\n", task.Title)

			tq.mu.Lock()
			tq.workingTasks = append(tq.workingTasks, task)
			tq.mu.Unlock()

			// 下载
			p, err := tq.manager.Select(task.Url)
			if err != nil {
				continue
			}

			stream, err := p.Service.Download(context.Background(), &pb.TaskRequest{
				Task: task,
			})
			if err != nil {
				continue
			}

			go func() {
				for {
					progress, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						fmt.Printf("Error receiving progress: %v\n", err)
						break
					}
					task.Percent = int64(progress.BytesTransferred / progress.TotalBytes)
					task.Speed = progress.Speed
					fmt.Printf("Download progress: %s - %d\n", progress.Id, progress.Speed)
				}
				tq.mu.Lock()
				tq.finishedTasks = append(tq.finishedTasks, task)
				tq.workingTasks = removeTask(task, tq.workingTasks)
				tq.mu.Unlock()
			}()

		default:

			if len(tq.queueTasks) == 0 {
				// 任务队列没任务 直接结束
				fmt.Println("无任务队列, 退出")
				return
			} else if len(tq.workingTasks) < tq.limit {
				//  任务队列有任务 但是正在下载的任务数量小于限制
				tq.Reload()
			}
		}
	}
}

func removeTask(task *pb.Task, tasks []*pb.Task) []*pb.Task {
	for i, t := range tasks {
		if t.Id == task.Id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks
		}
	}
	return tasks
}

// List 列出所有任务
func (tq *TaskQueue) List() {
	for _, task := range tq.queueTasks {
		fmt.Println(task.Id)
	}
}

// 重新装填, 把队列中的补充到任务通道
func (tq *TaskQueue) Reload() {
	free := tq.limit - len(tq.workingTasks)
	tasksToReload := tq.queueTasks[:min(free, len(tq.queueTasks))]
	tq.queueTasks = tq.queueTasks[len(tasksToReload):]

	for _, task := range tasksToReload {
		tq.taskChan <- task
	}
}

// 保存单个任务
func saveTask(srcTask *pb.Task, tasks []*pb.Task, configDir string) error {
	// parts := make([]Part, 0)

	// // 修改/更新
	// for _, task := range tasks {
	// 	if srcTask.part.TaskID == task.part.TaskID {
	// 		parts = append(parts, *srcTask.part)
	// 	} else {
	// 		parts = append(parts, *task.part)
	// 	}
	// }

	// tasksData, err := json.MarshalIndent(parts, "", "  ")
	// if err != nil {
	// 	logger.Error("" + err.Error())
	// }

	// err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }
	return nil
}

// task更新时 保存
func saveTasks(tasks []*pb.Task, configDir string) error {
	// parts := make([]Part, 0)

	// for _, task := range tasks {
	// 	part := *task.part
	// 	parts = append(parts, part)
	// }

	// tasksData, err := json.MarshalIndent(parts, "", "  ")
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	// err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }
	return nil
}

// func (a *App) loadTasks() error {
// 	tasks := make([]*Task, 0)

// 	configFile := filepath.Join(a.appInfo.configDir, "tasks.json")
// 	configData, err := os.ReadFile(configFile)
// 	if err != nil {
// 		logger.Errorf("Cannot read/find task file: %v", err)
// 		a.tasks = tasks
// 		return err
// 	}

// 	parts := make([]Part, 0)
// 	err = json.Unmarshal(configData, &parts)
// 	if err != nil {
// 		logger.Errorf("Cannot convert task data: %v", err)
// 		a.tasks = tasks
// 		return err
// 	}

// 	for _, part := range parts {
// 		// 过滤掉不存在的任务
// 		if _, err = os.Stat(part.DownloadDir); err == nil {
// 			newPart := part
// 			tasks = append(tasks, &Task{
// 				part: &newPart,
// 			})
// 		}
// 	}
// 	a.tasks = tasks
// 	saveTasks(tasks, a.configDir)
// 	return nil
// }
