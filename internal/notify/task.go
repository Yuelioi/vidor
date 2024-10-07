package notify

import (
	"context"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TaskNotification struct {
	ctx context.Context
}

func NewTaskNotification(ctx context.Context) *TaskNotification {
	return &TaskNotification{
		ctx: ctx,
	}
}

func (s *TaskNotification) UpdateTask(task *pb.Task) error {
	runtime.EventsEmit(s.ctx, "system.task", task)
	return nil
}
func (s *TaskNotification) UpdateTasks(tasks []*pb.Task) error {
	runtime.EventsEmit(s.ctx, "system.tasks", tasks)
	return nil
}
