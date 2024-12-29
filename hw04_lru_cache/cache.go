package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	item, exists := l.items[key]
	if exists {
		item.Value = value
		l.queue.MoveToFront(item)
		return true
	}
	newL := l.queue.PushFront(value)
	if l.queue.Len() > l.capacity {
		lastElem := l.queue.Back()
		for keyInMap, valInMap := range l.items {
			if valInMap.Value == lastElem.Value {
				delete(l.items, keyInMap)
				break
			}
		}
		l.queue.Remove(lastElem)
	}
	l.items[key] = newL
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	item, exists := l.items[key]
	if !exists {
		return nil, false
	}
	l.queue.MoveToFront(item)
	return item.Value, true
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem) // clear(l.items) Go ver 1.21+
	l.queue = NewList()
}
