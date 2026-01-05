package dto

import "time"

// RegisterRequest defines the input for user registration.
type RegisterRequest struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
	// ProfilePicture is handled via c.FormFile, so it's not strictly in the struct binding 
	// unless we use specific multipart binding, but Gin's binding for file is tricky.
	// We will handle file manually in controller but keep data validation here.
}

// LoginRequest defines the input for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest defines the input for updating user profile.
type UpdateUserRequest struct {
	Username string `form:"username"`
	// ProfilePicture handled manually
}

// UserResponse defines the standard output for user data.
type UserResponse struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	ProfilePicture string    `json:"profile_picture"`
	Initials       string    `json:"initials"`
	AvatarColor    string    `json:"avatar_color"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Token          string    `json:"token,omitempty"`
}
