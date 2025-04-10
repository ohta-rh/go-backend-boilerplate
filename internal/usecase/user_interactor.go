package usecase

import (
	"context"

	"github.com/tetsuyaohta/go-backend-boilerplate/internal/domain/user"
)

// UserInteractor implements use cases for user domain
type UserInteractor struct {
	userRepo user.Repository
}

// NewUserInteractor creates a new user interactor
func NewUserInteractor(repo user.Repository) *UserInteractor {
	return &UserInteractor{
		userRepo: repo,
	}
}

// GetUser retrieves a user by ID
func (i *UserInteractor) GetUser(ctx context.Context, id int) (*user.User, error) {
	return i.userRepo.GetByID(ctx, id)
}

// CreateUser creates a new user
func (i *UserInteractor) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	return i.userRepo.Create(ctx, u)
}

// UpdateUser updates an existing user
func (i *UserInteractor) UpdateUser(ctx context.Context, u *user.User) (*user.User, error) {
	return i.userRepo.Update(ctx, u)
}

// DeleteUser deletes a user by ID
func (i *UserInteractor) DeleteUser(ctx context.Context, id int) error {
	return i.userRepo.Delete(ctx, id)
}

// GetAllUsers retrieves all users
func (i *UserInteractor) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	return i.userRepo.GetAll(ctx)
}
