package main

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

var sseManager *SSEManager
var once sync.Once

type SSEManager struct {
	mu      *sync.RWMutex
	clients map[uuid.UUID][]chan NewPostPayload
}

// GetSSEManager returns the singleton instance of SSEManager.
func GetSSEManager() *SSEManager {
	once.Do(func() {
		sseManager = &SSEManager{
			mu:      &sync.RWMutex{},
			clients: make(map[uuid.UUID][]chan NewPostPayload),
		}
	})

	return sseManager
}

func (sseManager *SSEManager) Add(userId uuid.UUID, ch chan NewPostPayload) {
	sseManager.mu.Lock()
	defer sseManager.mu.Unlock()

	sseManager.clients[userId] = append(sseManager.clients[userId], ch)
}

func (sseManager *SSEManager) Remove(userId uuid.UUID, ch chan NewPostPayload) {
	sseManager.mu.Lock()
	defer sseManager.mu.Unlock()

	clients := sseManager.clients[userId]
	for idx, clientChannel := range clients {
		if clientChannel == ch {
			sseManager.clients[userId] = append(clients[:idx], clients[idx+1:]...)
			log.Printf("Removed channel for user: %s. Rem channels: %v\n", userId, len(sseManager.clients[userId]))
			break
		}
	}

	if len(sseManager.clients[userId]) == 0 {
		delete(sseManager.clients, userId)
	}

	log.Printf("Channel not found for user %s. Nothing to remove.\n", userId)
}

// Broadcasts message to all connected clients for that user.
func (sseManager *SSEManager) BroadcastToUser(userId uuid.UUID, message NewPostPayload) {
	sseManager.mu.RLock()
	defer sseManager.mu.RUnlock()

	for id, clients := range sseManager.clients {
		if id != userId {
			continue
		}

		for _, clientChannel := range clients {
			go func(clientChannel chan NewPostPayload) {
				clientChannel <- message
			}(clientChannel)
		}
	}
}
