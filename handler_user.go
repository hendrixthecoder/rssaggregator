package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hendrixthecoder/rssaggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	respondWithJSON(w, 200, databaseUserToUser(dbUser))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	page, limit := r.URL.Query().Get("page"), r.URL.Query().Get("limit")

	type PostCountResult struct {
		Count int
		Err   error
	}

	postCountChan := make(chan PostCountResult)

	pageInt, ok := parseQueryInt(page, "page", w)
	if !ok {
		return
	}

	limitInt, ok := parseQueryInt(limit, "limit", w)
	if !ok {
		return
	}

	go func() {
		totalPostsCount, err := apiCfg.DB.GetTotalPostCountForUser(r.Context(), dbUser.ID)
		if err != nil {
			log.Println("Error fetching all posts, err: ", err)
			postCountChan <- PostCountResult{Count: 0, Err: err}
			return
		}

		postCountChan <- PostCountResult{Count: int(totalPostsCount), Err: nil}
	}()

	offsetInt := (pageInt - 1) * limitInt

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: dbUser.ID,
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching posts: %v", err))
	}

	type Response struct {
		Posts           []Post `json:"posts"`
		HasNextPage     bool   `json:"hasNextPage"`
		HasPreviousPage bool   `json:"hasPreviousPage"`
	}

	parsed_posts := dtoSliceSerializer(posts, databasePostToPost)
	postCountResult := <-postCountChan
	if postCountResult.Err != nil {
		respondWithError(w, 500, "Error getting post count")
		return
	}

	hasNextPage := offsetInt+limitInt < postCountResult.Count
	hasPreviousPage := offsetInt > 0

	respondWithJSON(w, 200, Response{
		Posts:           parsed_posts,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	})
}

func (apiConfig *apiConfig) handlerSearchUserPosts(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	searchStr := r.URL.Query().Get("q")
	if searchStr == "" {
		respondWithError(w, 400, "Search key can not be empty.")
		return
	}

	posts, err := apiConfig.DB.GetPostsMatchingSearchTerm(r.Context(), database.GetPostsMatchingSearchTermParams{
		UserID: dbUser.ID,
		Title:  "%" + searchStr + "%", // Searches title & description
	})
	if err != nil {
		log.Println("Error fetching searched posts: ", err)
		respondWithError(w, 400, "Error fetching searched posts.")
		return
	}

	respondWithJSON(w, 200, dtoSliceSerializer(posts, databasePostToPost))
}

func (apiCfg *apiConfig) handlerSSEHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Fatal("Streaming unsupported")
	}

	ctx := r.Context()
	newPostChan := make(chan NewPostPayload, 10)

	sseManager := GetSSEManager()
	sseManager.Add(dbUser.ID, newPostChan)

	for {
		select {
		case <-ctx.Done():
			sseManager.mu.Lock()
			defer sseManager.mu.Unlock()

			sseManager.Remove(dbUser.ID, newPostChan)
			return

		case post := <-newPostChan:
			data, err := json.Marshal(post)
			if err != nil {
				log.Println("Failed to marshal post:", err)
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-time.After(30 * time.Second):
			// Send comment to keep connection alive
			fmt.Fprint(w, ": keep-alive\n\n")
			flusher.Flush()
		}
	}
}
