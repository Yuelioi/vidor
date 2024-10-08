package notify

import (
	"github.com/sirupsen/logrus"
)

type LoggingNotification struct {
	logger  *logrus.Logger
	wrapped Notification
}

func (ln *LoggingNotification) Send(nc Notice) error {
	ln.logger.Infof("[%s]%s from %s", nc.EventName, nc.Provider, nc.Content)
	err := ln.wrapped.Send(nc)
	if err != nil {
		ln.logger.Errorf("发送消息失败: %v", err)
		return err
	}
	return nil
}

func NewLoggingNotification(logger *logrus.Logger, wrapped Notification) *LoggingNotification {
	return &LoggingNotification{
		logger:  logger,
		wrapped: wrapped,
	}
}
