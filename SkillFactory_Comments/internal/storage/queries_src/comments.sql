-- name: GetAllComments :many
SELECT id, news_id, parent_comment_id, content, created_at 
FROM comments
ORDER BY id;

-- name: GetCommentById :many
SELECT id, news_id, parent_comment_id, content, created_at 
FROM comments
WHERE news_id = $1 OR parent_comment_id = $1
ORDER BY created_at;

-- name: CreateComments :exec
INSERT INTO comments 
(news_id, parent_comment_id, content) 
VALUES ($1, $2, $3);

-- name: DeleteComment :exec
DELETE FROM comments 
WHERE id = $1;