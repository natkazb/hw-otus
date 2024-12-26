package hw04lrucache

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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	item, exists := l.items[key]
	if exists {
		item.Value = value
		return true
	} else {
		newL := l.queue.PushBack(value)
		l.items[key] = newL
		return false
	}
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, exists := l.items[key]
	if !exists {
		return nil, false
	}
	return item.Value, true
}

func (l *lruCache) Clear() {
	l.capacity = 0
	l.items = make(map[Key]*ListItem)
	if l.queue.Len() > 0 {
		item := l.queue.Front()
		for item != nil {
			next := item.Next
			l.queue.Remove(item)
			item = next
		}
	}
}
