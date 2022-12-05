package google

import (
	"log"
	"sync"
)

// MutexKV is a simple key/value store for arbitrary mutexes. It can be used to
// serialize changes across arbitrary collaborators that share knowledge of the
// keys they must serialize on.
//
// The initial use case is to let aws_security_group_rule resources serialize
// their access to individual security groups based on SG ID.
type MutexKV struct {
	lock  sync.Mutex
	store map[string]*sync.RWMutex
}

// Locks the mutex for the given key. Caller is responsible for calling Unlock
// for the same key
func (m *MutexKV) Lock(key string) {
	log.Printf("[DEBUG] Locking %q", key)
	m.get(key).Lock()
	log.Printf("[DEBUG] Locked %q", key)
}

// Unlock the mutex for the given key. Caller must have called Lock for the same key first
func (m *MutexKV) Unlock(key string) {
	log.Printf("[DEBUG] Unlocking %q", key)
	m.get(key).Unlock()
	log.Printf("[DEBUG] Unlocked %q", key)
}

// Acquires a read-lock on the mutex for the given key. Caller is responsible for calling RUnlock
// for the same key
func (m *MutexKV) RLock(key string) {
	log.Printf("[DEBUG] RLocking %q", key)
	m.get(key).RLock()
	log.Printf("[DEBUG] RLocked %q", key)
}

// Releases a read-lock on the mutex for the given key. Caller must have called RLock for the same key first
func (m *MutexKV) RUnlock(key string) {
	log.Printf("[DEBUG] RUnlocking %q", key)
	m.get(key).RUnlock()
	log.Printf("[DEBUG] RUnlocked %q", key)
}

// Returns a mutex for the given key, no guarantee of its lock status
func (m *MutexKV) get(key string) *sync.RWMutex {
	m.lock.Lock()
	defer m.lock.Unlock()
	mutex, ok := m.store[key]
	if !ok {
		mutex = &sync.RWMutex{}
		m.store[key] = mutex
	}
	return mutex
}

// Returns a properly initialized MutexKV
func NewMutexKV() *MutexKV {
	return &MutexKV{
		store: make(map[string]*sync.RWMutex),
	}
}
