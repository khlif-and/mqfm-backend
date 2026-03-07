-- name: ListAudios :many
SELECT id, title, artist, description, file_path, duration, status, thumbnail, category_id, created_at, updated_at
FROM audios WHERE deleted_at IS NULL ORDER BY id;

-- name: GetAudioByID :one
SELECT id, title, artist, description, file_path, duration, status, thumbnail, category_id, created_at, updated_at
FROM audios WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: SearchAudios :many
SELECT id, title, artist, description, file_path, duration, status, thumbnail, category_id, created_at, updated_at
FROM audios WHERE (title LIKE CONCAT('%', ?, '%') OR artist LIKE CONCAT('%', ?, '%')) AND deleted_at IS NULL;

-- name: CreateAudio :execresult
INSERT INTO audios (title, artist, description, file_path, duration, status, thumbnail, category_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);
