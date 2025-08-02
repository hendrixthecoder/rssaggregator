-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_new_post() RETURNS trigger AS $$
DECLARE
    payload JSON;
BEGIN
    payload := json_build_object(
        'post_id', NEW.id,
        'feed_id', NEW.feed_id,
        'title', NEW.title,
        'user_id', NEW.user_id,
        'description', NEW.description,
        'url', NEW.url,
        'published_at', NEW.published_at
    );
    PERFORM pg_notify('new_post', payload::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER new_post_notify
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION notify_new_post();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS new_post_notify ON posts;
DROP FUNCTION IF EXISTS notify_new_post();
-- +goose StatementEnd
