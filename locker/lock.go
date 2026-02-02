package locker

import "sync"

// locks (in memory) by "id".
// example: there can be a Locker object that prevents
// simultaneous modifications to a customer
// identified by "id"
type Locker struct {
	globalMutex *sync.Mutex
	mutexes     map[string]*sync.Mutex
}

func NewLocker() *Locker {
	return &Locker{
		globalMutex: &sync.Mutex{},
		mutexes:     map[string]*sync.Mutex{},
	}
}

// locks a lock identified by "id"
func (h *Locker) Lock(id string) {
	h.globalMutex.Lock()

	bM, ok := h.mutexes[id]
	if !ok {
		bM = &sync.Mutex{}
		h.mutexes[id] = bM
	}
	h.globalMutex.Unlock()

	bM.Lock()
}

// unlocks a lock identified by "id"
func (h *Locker) Unlock(id string) {

	h.globalMutex.Lock()
	defer h.globalMutex.Unlock()

	bM, ok := h.mutexes[id]
	if !ok {
		// if not found, do nothing
		return
	}

	bM.Unlock()
}
