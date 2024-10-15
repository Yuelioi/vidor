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
}

// New 创建一个新的任务队列
func NewTaskQueue(plugin *plugin.DownloadPlugin, onFinished func(*pb.Task, error)) *TaskQueue {
	return &TaskQueue{
		plugin:     plugin,
		stop:       make(chan chan struct{}),
		working:    atomic.Bool{},
		onFinished: onFinished,
	}
}

func (tq *TaskQueue) work(task *pb.Task) {
	stream, err := tq.plugin.Service.Download(context.Background(), &pb.TaskRequest{
		Task: task,
	})
	if err != nil {
		tq.onFinished(task, nil)
		return
	}

	for {
		progress, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			tq.onFinished(task, nil)
			fmt.Printf("Error receiving progress: %v\n", err)
			break
		}
		// 持续更新: 状态 百分比 速度
		task.Status = progress.Status
		task.Percent = progress.Percent
		task.Speed = progress.Speed
		task.Cover = progress.Cover

		if task.Cover != "" {
			task.Cover = "/files/" + task.Cover
		}
	}
	tq.onFinished(task, nil)
}
