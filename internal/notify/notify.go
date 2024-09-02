package notify

import (
	"context"

	"github.com/sirupsen/logrus"
)

type NotificationService interface {
	Send(ctx context.Context, provider, eventName, message, messageType string)
}

type Notification struct {
	ctx    context.Context
	logger *logrus.Logger
}

type Notice struct {
	Message     string `json:"message"`     // 消息内容
	MessageType string `json:"messageType"` // 消息类型 ['success', 'info', 'warn', 'error', 'secondary', 'contrast']
	Provider    string `json:"provider"`    // 消息提供者
}

func New(ctx context.Context, logger *logrus.Logger) *Notification {
	return &Notification{}
}

func (n Notification) Send(ns NotificationService, provider, eventName, message, messageType string) {
	n.logger.Infof("%s发送了%s消息[消息名:%s][消息主体:%s]", provider, messageType, eventName, message)
	ns.Send(n.ctx, provider, eventName, message, messageType)
}
