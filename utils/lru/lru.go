package lru

import (
	"sync"

	lowlevel "github.com/NikoMalik/low-level-functions"
)

type LRU struct {
	items   map[string]*Element[string]
	list    *List[string]
	mu      sync.Mutex
	maxSize int
}

func New(maxSize int) *LRU {
	if maxSize < 1 {
		panic("assertion error: maxSize < 1")
	}
	return &LRU{
		maxSize: maxSize,
		items:   make(map[string]*Element[string], maxSize),
		list:    NewList[string](),
	}
}

// GetOrAdd fetch item from lru and increase eviction order or create
func (l *LRU) GetOrAdd(keyB []byte) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	element, ok := l.items[lowlevel.String(keyB)]
	if ok {
		l.list.MoveToFront(element)
		return element.Value
	}

	if len(l.items) >= l.maxSize {
		element = l.list.Back()
		l.list.Remove(element)
		delete(l.items, element.Value)
	}

	keyS := lowlevel.String(keyB)
	element = l.list.PushFront(keyS)
	l.items[keyS] = element
	return keyS
}
