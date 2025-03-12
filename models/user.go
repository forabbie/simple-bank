package models

import "time"

// loginUserRequest defines the request structure for user login.
type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	ID                int64      `uri:"id" binding:"required,min=1"`
	Password          *string    `json:"password" binding:"omitempty,min=6"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	FullName          *string    `json:"full_name"`
	Email             *string    `json:"email,omitempty" binding:"omitempty,email"`
}
