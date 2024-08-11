package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func CreateLogger(appTempDir string) (*logrus.Logger, error) {
	logFilePath := filepath.Join(appTempDir, fmt.Sprintf("logger_%s.txt", time.Now().Format("20060102_150405")))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	logger.SetOutput(multiWriter)
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logger.SetFormatter(formatter)
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("初始化完毕")
	return logger, nil
}
