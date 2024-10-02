package storage

import (
	"sync"
	"time"
)

var (
	blacklist = make(map[string]time.Time)
	mu        sync.Mutex
)

func Add(token string, expiration time.Time) {
	mu.Lock()
	defer mu.Unlock()
	blacklist[token] = expiration
}

func IsBlacklisted(token string) bool {
	mu.Lock()
	defer mu.Unlock()
	expiration, exists := blacklist[token]
	if !exists {
		return false
	}
	if time.Now().After(expiration) {
		delete(blacklist, token)
		return false
	}
	return true
}