-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at , title , description, published_at, url, feed_id)
VALUES($1, $2, $3, $4, $5, $6) 