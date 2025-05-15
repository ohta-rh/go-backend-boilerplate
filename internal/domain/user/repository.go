package user

import (
	"context"
)

// Repository defines the methods to interact with user storage.
type Repository interface {
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int) (*User, error)

	// Create stores a new user
	Create(ctx context.Context, user *User) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) (*User, error)

	// Delete removes a user by ID
	Delete(ctx context.Context, id int) error

	// GetAll retrieves all users
	GetAll(ctx context.Context) ([]*User, error)
}
