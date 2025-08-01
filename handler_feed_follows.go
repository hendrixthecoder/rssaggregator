package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hendrixthecoder/rssaggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    dbUser.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeed(feed_follow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	feed_follows, err := apiCfg.DB.GetFeedFollows(r.Context(), dbUser.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching feed follows: %v", err))
		return
	}

	respondWithJSON(w, 200, dtoSliceSerializer(feed_follows, databaseFeedFollowToFeed))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowId")

	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid UUID: %v", err))
		return
	}

	deleteErr := apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: dbUser.ID,
	})
	if deleteErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow: %v", err))
	}

	respondWithJSON(w, 204, fmt.Sprintln("Record delete successfully!"))
}
