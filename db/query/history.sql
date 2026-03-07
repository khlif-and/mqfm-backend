-- name: GetHistoryByUser :many
SELECT h.id, h.user_id, h.audio_id, h.play_count, h.played_at, h.created_at,
       a.title as audio_title, a.artist as audio_artist
FROM histories h
LEFT JOIN audios a ON h.audio_id = a.id
WHERE h.user_id = ?
ORDER BY h.played_at DESC;

-- name: UpsertHistory :exec
INSERT INTO histories (user_id, audio_id, play_count, played_at)
VALUES (?, ?, 1, NOW())
ON DUPLICATE KEY UPDATE play_count = play_count + 1, played_at = NOW();

-- name: DeleteHistory :exec
DELETE FROM histories WHERE user_id = ? AND audio_id = ?;

-- name: ClearHistory :exec
DELETE FROM histories WHERE user_id = ?;
