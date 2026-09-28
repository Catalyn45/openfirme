package common

import "sync"

type ConcurentMap[K comparable, V any] struct {
	mutex sync.RWMutex
	cond *sync.Cond
	container  map[K]V
	maxCapacity int
}

func NewConcurentMap[K comparable, V any](maxCapacity int) *ConcurentMap[K, V] {
	concurentMap := &ConcurentMap[K, V] {
		container: make(map[K]V),
		maxCapacity: maxCapacity,
	}

	concurentMap.cond = sync.NewCond(&concurentMap.mutex)

	return concurentMap
}

func (this *ConcurentMap[K, V]) Set(key K, value V) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.container) >= this.maxCapacity {
		return
	}

	this.container[key] = value

	this.cond.Signal()
}

func (this *ConcurentMap[K, V]) Remove(key K) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.container, key)
}

func (this *ConcurentMap[K, V]) RemoveAll(container map[K]V) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for key, _ := range container {
		delete(this.container, key)
	}
}

func (this *ConcurentMap[K, V]) Get() (K, V) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for len(this.container) == 0 {
		this.cond.Wait()
	}

	for key, value := range this.container {
		return key, value
	}

	panic("unreachable")
}

func (this *ConcurentMap[K, V]) IsEmpty() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return len(this.container) == 0
}

func (this *ConcurentMap[K, V]) Clone() map[K]V {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for len(this.container) == 0 {
		this.cond.Wait()
	}

	container := make(map[K]V)
	for k, v := range this.container {
		container[k] = v
	}

	return container
}
