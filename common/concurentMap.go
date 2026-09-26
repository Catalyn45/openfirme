package common

import "sync"

type ConcurentMap[K comparable, V any] struct {
	mutex sync.RWMutex
	container  map[K]V
	maxCapacity int
}

func NewConcurentMap[K comparable, V any](maxCapacity int) *ConcurentMap[K, V] {
	return &ConcurentMap[K, V] {
		container: make(map[K]V),
		maxCapacity: maxCapacity,
	}
}

func (this *ConcurentMap[K, V]) Set(key K, value V) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.container) >= this.maxCapacity {
		return
	}

	this.container[key] = value
}

func (this *ConcurentMap[K, V]) Remove(key K) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.container, key)
}

func (this *ConcurentMap[K, V]) IsEmpty() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return len(this.container) == 0
}

func (this *ConcurentMap[K, V]) Move() map[K]V {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	container := this.container

	this.container = make(map[K]V)

	return container
}
