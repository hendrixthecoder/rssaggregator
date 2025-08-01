-- name: CreatePost :one
INSERT INTO posts (
    id, 
    created_at, 
    updated_at, 
    url, 
    title, 
    feed_id,
    description,
    published_at
) 

VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;