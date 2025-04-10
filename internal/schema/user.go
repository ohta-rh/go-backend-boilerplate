package schema

import (
	"time"

	domainUser "easy-go-backend/internal/domain/user"
)

// UserRequest represents the request schema for user operations.
type UserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UserResponse represents the response schema for user operations.
type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToEntity converts a UserRequest to a domain User entity.
func (r *UserRequest) ToEntity() *domainUser.User {
	return &domainUser.User{
		Name:  r.Name,
		Email: r.Email,
	}
}

// FromEntity creates a UserResponse from a domain User entity.
func FromEntity(u *domainUser.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// FromEntities creates a slice of UserResponse from a slice of domain User entities.
func FromEntities(users []*domainUser.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = FromEntity(u)
	}

	return responses
}
