-- name: CreatePost :one
INSERT INTO posts (
    id, 
    created_at, 
    updated_at, 
    url, 
    title, 
    feed_id,
    description,
    published_at,
    user_id
) 

VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* 
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTotalPostCountForUser :one
SELECT COUNT(*)
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: GetPostsMatchingSearchTerm :many
SELECT posts.*
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1 
    AND (posts.title ILIKE $2 OR posts.description ILIKE $2)
LIMIT 10;