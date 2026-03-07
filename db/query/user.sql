-- name: GetUserByEmail :one
SELECT id, username, email, password, profile_picture, role, provider, provider_id, created_at, updated_at
FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, email, password, profile_picture, role, provider, provider_id, created_at, updated_at
FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByProviderID :one
SELECT id, username, email, password, profile_picture, role, provider, provider_id, created_at, updated_at
FROM users WHERE provider = ? AND provider_id = ? AND deleted_at IS NULL LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, email, password, profile_picture, role, provider, provider_id) VALUES (?, ?, ?, ?, ?, ?, ?);
