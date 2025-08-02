-- +goose Up
ALTER TABLE posts 
ADD COLUMN user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;
-- +goose Down
ALTER TABLE posts
DROP COLUMN user_id;