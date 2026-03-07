-- name: ListCategories :many
SELECT id, name, description, image, created_at, updated_at
FROM categories WHERE deleted_at IS NULL ORDER BY id;

-- name: GetCategoryByID :one
SELECT id, name, description, image, created_at, updated_at
FROM categories WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: SearchCategories :many
SELECT id, name, description, image, created_at, updated_at
FROM categories WHERE name LIKE CONCAT('%', ?, '%') AND deleted_at IS NULL;

-- name: CreateCategory :execresult
INSERT INTO categories (name, description, image) VALUES (?, ?, ?);
