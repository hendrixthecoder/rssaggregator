package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/hendrixthecoder/rssaggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	FeedID      uuid.UUID `json:"feed_id"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
}

type NewPostPayload struct {
	PostID      uuid.UUID `json:"post_id"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"decription"`
	PublishedAt string    `json:"published_at"`
	Url         string    `json:"url"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		UserID:    dbFeed.UserID,
		Url:       dbFeed.Url,
	}
}

func databaseFeedFollowToFeed(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func databasePostToPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Url:         dbPost.Url,
		Title:       dbPost.Title,
		FeedID:      dbPost.FeedID,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
	}
}

// Serializer func to serialize slice of data from Goose types to JSON-normalized type.
func dtoSliceSerializer[T any, DTO any](data []T, serializer func(T) DTO) []DTO {
	serialized := make([]DTO, len(data))

	for idx, feed := range data {
		serialized[idx] = serializer(feed)
	}

	return serialized
}
