-- name: GetAdminByEmail :one
SELECT id, username, email, password, role, created_at, updated_at
FROM admins WHERE email = ? AND deleted_at IS NULL LIMIT 1;

-- name: GetAdminByID :one
SELECT id, username, email, password, role, created_at, updated_at
FROM admins WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: CreateAdmin :execresult
INSERT INTO admins (username, email, password, role) VALUES (?, ?, ?, ?);
