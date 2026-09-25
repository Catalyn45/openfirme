package common

import "sync"

type Empty struct{}

type ConcurentSet[T comparable] struct {
	mutex sync.RWMutex
	container  map[T]Empty
	maxCapacity int
}

func NewSet[T comparable](maxCapacity int) *ConcurentSet[T] {
	return &ConcurentSet[T]{
		container: make(map[T]Empty),
		maxCapacity: maxCapacity,
	}
}

func (this *ConcurentSet[T]) Add(v T) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.container) >= this.maxCapacity {
		return
	}

	this.container[v] = Empty{}
}

func (this *ConcurentSet[T]) Remove(v T) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.container, v)
}

func (this *ConcurentSet[T]) Has(v T) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	_, ok := this.container[v]
	return ok
}

func (this *ConcurentSet[T]) IsEmpty() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return len(this.container) == 0
}

func (this *ConcurentSet[T]) IteratePop(callback func (param T)) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for el, _ := range this.container {
		callback(el)
	}

	clear(this.container)
}
