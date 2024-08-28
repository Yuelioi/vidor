package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Notice struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
	Duration    int    `json:"duration"`
}

// 消息分发系统
type N struct {
	ctx context.Context
	log *logrus.Logger
}

func (n *N) Notify(eventName string, message string, messageType string, duration int) {
	n.log.Infof(message)
	runtime.EventsEmit(n.ctx, eventName, &Notice{
		Message:     message,
		MessageType: messageType,
		Duration:    duration,
	})
}
