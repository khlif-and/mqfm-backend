package category

import "time"

type CreateCategoryRequest struct {
	Name string `form:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string `form:"name"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
