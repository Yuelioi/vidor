package task

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/Yuelioi/vidor/internal/plugin"
	pb "github.com/Yuelioi/vidor/internal/proto"
)

// TaskQueue 任务队列接口
type TaskQueue struct {
	plugin     *plugin.DownloadPlugin
	stop       chan chan struct{}
	working    atomic.Bool
	onFinished func(*pb.Task, error)
	task       *pb.Task
}

// New 创建一个新的任务队列
func NewTaskQueue(plugin *plugin.DownloadPlugin, task *pb.Task, onFinished func(*pb.Task, error)) *TaskQueue {
	return &TaskQueue{
		plugin:     plugin,
		stop:       make(chan chan struct{}),
		working:    atomic.Bool{},
		onFinished: onFinished,
		task:       task,
	}
}

func (tq *TaskQueue) work() {
	stream, err := tq.plugin.Service.Download(context.Background(), &pb.TaskRequest{
		Task: tq.task,
	})
	if err != nil {
		tq.onFinished(tq.task, nil)
		return
	}

	for {
		progress, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			tq.onFinished(tq.task, nil)
			fmt.Printf("Error receiving progress: %v\n", err)
			break
		}
		// 持续更新: 状态 百分比 速度
		tq.task.Status = progress.Status
		tq.task.Percent = progress.Percent
		tq.task.Speed = progress.Speed
		tq.task.Cover = progress.Cover

		if tq.task.Cover != "" {
			tq.task.Cover = "/files/" + tq.task.Cover
		}
	}
	tq.onFinished(tq.task, nil)
}
