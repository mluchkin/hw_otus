package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    sync.RWMutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lr *lruCache) Set(key Key, value interface{}) bool {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	ci := cacheItem{key: string(key), value: value}
	if v, ok := lr.items[key]; ok {
		v.Value = ci
		lr.queue.MoveToFront(v)
		return true
	}

	if lr.capacity == lr.queue.Len() {
		v := lr.queue.Back()
		lr.queue.Remove(v)
		delete(lr.items, Key(v.Value.(cacheItem).key))
	}

	l := lr.queue.PushFront(ci)
	lr.items[key] = l

	return false
}

func (lr *lruCache) Get(key Key) (interface{}, bool) {
	lr.mutex.RLock()
	defer lr.mutex.RUnlock()
	if v, ok := lr.items[key]; ok {
		lr.queue.MoveToFront(v)
		return v.Value.(cacheItem).value, true
	}
	return nil, false
}

func (lr *lruCache) Clear() {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	for i := range lr.items {
		delete(lr.items, i)
	}

	for i := lr.queue.Back(); i != nil; i = i.Prev {
		lr.queue.Remove(i)
	}
}
