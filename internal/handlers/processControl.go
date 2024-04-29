// AccessControl.go
package handlers

import "sync"

var (
	processing = make(map[string]bool)
	mu         sync.Mutex
)

func SetProcessing(sha256 string, status bool) {
	mu.Lock()
	defer mu.Unlock()
	processing[sha256] = status
}

func IsProcessing(sha256 string) bool {
	mu.Lock()
	defer mu.Unlock()
	return processing[sha256]
}

func ClearProcessing(sha256 string) {
	mu.Lock()
	defer mu.Unlock()
	delete(processing, sha256)
}
