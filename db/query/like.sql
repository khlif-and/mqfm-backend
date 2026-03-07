-- name: GetLikesByUser :many
SELECT l.id, l.user_id, l.audio_id, l.created_at,
       a.title as audio_title, a.artist as audio_artist
FROM likes l
LEFT JOIN audios a ON l.audio_id = a.id
WHERE l.user_id = ?
ORDER BY l.created_at DESC;

-- name: CheckLikeExists :one
SELECT COUNT(*) as count FROM likes WHERE user_id = ? AND audio_id = ?;

-- name: CreateLike :execresult
INSERT INTO likes (user_id, audio_id) VALUES (?, ?);

-- name: DeleteLike :exec
DELETE FROM likes WHERE user_id = ? AND audio_id = ?;
