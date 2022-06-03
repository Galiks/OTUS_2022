package cache

import (
	"errors"
	"sync"

	list "github.com/Galiks/OTUS_2022/hw04_lru_cache/list"
)

var ErrCapacity = errors.New("capacity must be great than 1")

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear() // Правильно понимаю, что это метод для удаления из кэша элемента, а не очистка всего кэша?
}

type lruCache struct {
	capacity int
	queue    list.List
	items    map[Key]*list.Item
	lock     sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity uint32) (Cache, error) {
	if capacity < 1 {
		return nil, ErrCapacity
	}

	return &lruCache{
		capacity: int(capacity),
		queue:    list.NewList(),
		items:    make(map[Key]*list.Item, capacity),
	}, nil
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.items[key]; ok {
		c.queue.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		return ok
	}
	if c.queue.Len() == c.capacity {
		c.Clear()
	}

	elem := c.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	c.items[key] = elem
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	elem, ok := c.items[key]
	if !ok {
		return nil, ok
	}

	c.queue.MoveToFront(elem)
	return elem.Value.(*cacheItem).value, ok
}

func (c *lruCache) Clear() {
	if elem := c.queue.Back(); elem != nil {
		c.queue.Remove(elem)
		delete(c.items, elem.Value.(*cacheItem).key)
	}
}
