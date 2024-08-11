package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Yuelioi/vidor/shared"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
)

func saveConfig(configDir string, config shared.Config) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configDir, "config.json"), configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// task更新时 保存
func saveTasks(tasks []*Task, configDir string) error {
	parts := make([]shared.Part, 0)

	for _, task := range tasks {
		part := *task.part
		parts = append(parts, part)
	}

	tasksData, err := json.MarshalIndent(parts, "", "  ")
	if err != nil {
		logger.Error(err)
		return err
	}

	err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func loadLocalPlugin(pluginPath string) (shared.Downloader, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"downloader": &shared.DownloaderRPCPlugin{},
		},
		Cmd: exec.Command(pluginPath),
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("error creating client for plugin %s: %v", pluginPath, err)
	}

	raw, err := rpcClient.Dispense("downloader")
	if err != nil {
		return nil, fmt.Errorf("error dispensing plugin %s: %v", pluginPath, err)
	}

	downloader, ok := raw.(shared.Downloader)
	if !ok {
		return nil, fmt.Errorf("plugin %s does not implement the expected interface", pluginPath)
	}

	return downloader, nil
}

// 保存单个任务
func saveTask(srcTask *Task, tasks []*Task, configDir string) error {
	parts := make([]shared.Part, 0)

	// 修改/更新
	for _, task := range tasks {
		if srcTask.part.TaskID == task.part.TaskID {
			parts = append(parts, *srcTask.part)
		} else {
			parts = append(parts, *task.part)
		}
	}

	tasksData, err := json.MarshalIndent(parts, "", "  ")
	if err != nil {
		logger.Error("" + err.Error())
	}

	err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// 创建下载器
func newDownloader(plugins []shared.PluginMeta, config shared.Config, notice shared.Notice, link string) (shared.Downloader, error) {
	for _, plugin := range plugins {
		for _, regex := range plugin.Regexs {
			if regex.MatchString(link) {
				return plugin.New(context.Background(), config), nil
			}
		}
	}
	return nil, errors.New("没有对应的下载器")
}

/*
创建任务

 1. 获取下载器
 2. 填充uid等初始化数据
*/
func createNewTask(part shared.Part, downloadDir, workName string) (*Task, error) {
	return &Task{
		state: Queue,
		part: &shared.Part{
			TaskID:      uuid.New().String(),
			DownloadDir: filepath.Join(downloadDir, workName),
			URL:         part.URL,
			Title:       part.Title,
			Thumbnail:   part.Thumbnail,
			Video:       part.Video,
			Audio:       part.Audio,
			Status:      shared.TaskStatus.Queue,
			CreatedAt:   time.Now(),
			State:       shared.TaskStatus.Queue,
		},
	}, nil
}

// 基于app任务的url判断任务是否存在
func taskExists(tasks []*Task, url string) bool {
	for _, existingTask := range tasks {
		if existingTask.part.URL == url {
			return true
		}
	}
	return false
}
