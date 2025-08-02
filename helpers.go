package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"
)

func parseQueryInt(value, fieldName string, w http.ResponseWriter) (int, bool) {
	i, err := strconv.Atoi(value)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid %s type provided: %s", fieldName, value))
		return 0, false
	}
	return i, true
}

func startNewPostListener(dsn string) {
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(event pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("Listener error:", err)
		}
	})

	defer listener.Close()

	err := listener.Listen("new_post")
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	log.Println("Listening for notifications on channel 'new_post'")

	for {
		select {
		case notification := <-listener.Notify:

			var newPostPayload NewPostPayload
			err := json.Unmarshal([]byte(notification.Extra), &newPostPayload)
			if err != nil {
				log.Println("Unable to parse new post payload:", err)
				continue
			}

			sseManager := GetSSEManager()
			sseManager.BroadcastToUser(newPostPayload.UserID, newPostPayload)

		case <-time.After(90 * time.Second):
			// Periodic ping to keep connection alive
			log.Println("No events received for 90 seconds, checking connection")
			go func() {
				if err := listener.Ping(); err != nil {
					log.Println("Ping error:", err)
				}
			}()
		}
	}
}
