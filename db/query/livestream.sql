-- name: GetLivestreamStatus :one
SELECT id, is_live, video_id, title, thumbnail, last_checked
FROM live_streams LIMIT 1;
