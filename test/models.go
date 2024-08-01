package test

import (
	"sync"
	"sync/atomic"

	"github.com/Yuelioi/vidor/shared"
	"github.com/sirupsen/logrus"
)

type Part struct {
	State  int
	TaskID string
	URL    string
}

type App struct {
	taskQueue *TaskQueue // 任务队列 用于分发任务 同一时刻只会出现一个队列
	tasks     []*Task    // 所有任务 包括不在下载范围内的
	Notice    shared.Notice
	Logger    *logrus.Logger
	Callback  shared.Callback
}
type Task struct {
	part  Part
	state int
}

type TaskQueue struct {
	state          int
	wg             sync.WaitGroup
	mu             sync.Mutex
	tasksRemaining atomic.Int64
	tasks          chan *Task // 任务通道
	queueTasks     []*Task    // 队列任务
	done           chan struct{}
}
