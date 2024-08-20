package app

import (
	"sync"

	pb "github.com/Yuelioi/vidor/internal/proto"
)

type Cache struct {
	cache sync.Map
}

func NewCache() *Cache {
	return &Cache{
		cache: sync.Map{},
	}
}

func (c *Cache) Get(id string) (*pb.Task, bool) {
	value, exists := c.cache.Load(id)
	if !exists {
		return nil, false
	}
	return value.(*pb.Task), true
}

func (c *Cache) Set(id string, info *pb.Task) {
	c.cache.Store(id, info)
}

func (c *Cache) Delete(id string) {
	c.cache.Delete(id)
}

func (c *Cache) Clear() {
	c.cache = sync.Map{}
}
