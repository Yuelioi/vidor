package app

import (
	"sync"

	"github.com/Yuelioi/vidor/internal/plugin"
	pb "github.com/Yuelioi/vidor/internal/proto"
)

type Cache struct {
	tasks      sync.Map
	downloader *plugin.DownloadPlugin
}

func NewCache() *Cache {
	return &Cache{
		tasks: sync.Map{},
	}
}

// 任务缓存
func (c *Cache) Task(id string) (*pb.Task, bool) {
	value, exists := c.tasks.Load(id)
	if !exists {
		return nil, false
	}
	return value.(*pb.Task), true
}

func (c *Cache) AddTask(id string, info *pb.Task) {
	c.tasks.Store(id, info)
}

func (c *Cache) DeleteTask(id string) {
	c.tasks.Delete(id)
}

func (c *Cache) ClearTasks() {
	c.tasks = sync.Map{}
}

func (c *Cache) AddTasks(tasks []*pb.Task) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(task *pb.Task) {
			defer wg.Done()
			c.AddTask(task.Id, task)
		}(task)
	}
	wg.Wait()
}

func (c *Cache) Tasks(ids []string) ([]*pb.Task, error) {
	var wg sync.WaitGroup
	tasks := make([]*pb.Task, len(ids))
	for i, id := range ids {
		wg.Add(1)
		go func(i int, id string) {
			defer wg.Done()
			task, ok := c.Task(id)
			if ok {
				tasks[i] = task
			}
		}(i, id)
	}
	wg.Wait()
	return tasks, nil
}

// 下载器缓存
func (c *Cache) Downloader() *plugin.DownloadPlugin {
	return c.downloader
}
func (c *Cache) SetDownloader(p *plugin.DownloadPlugin) {
	c.downloader = p
}
