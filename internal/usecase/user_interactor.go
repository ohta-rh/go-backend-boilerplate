package usecase

import (
	"context"
	"errors"

	domainUser "easy-go-backend/internal/domain/user"
)

// UserInteractor implements use cases for user domain.
type UserInteractor struct {
	userRepo domainUser.Repository
}

// NewUserInteractor creates a new user interactor.
func NewUserInteractor(repo domainUser.Repository) *UserInteractor {
	return &UserInteractor{
		userRepo: repo,
	}
}

// GetUser retrieves a user by ID.
func (i *UserInteractor) GetUser(ctx context.Context, id int) (*domainUser.User, error) {
	u, err := i.userRepo.GetByID(ctx, id)
	if err != nil && errors.Is(err, domainUser.ErrUserNotFound) {
		// ユーザーが見つからない場合は、ErrUserNotFoundを返す
		return nil, domainUser.ErrUserNotFound
	}

	return u, err
}

// CreateUser creates a new user.
func (i *UserInteractor) CreateUser(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	return i.userRepo.Create(ctx, u)
}

// UpdateUser updates an existing user.
func (i *UserInteractor) UpdateUser(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	updatedUser, err := i.userRepo.Update(ctx, u)
	if err != nil && errors.Is(err, domainUser.ErrUserNotFound) {
		return nil, domainUser.ErrUserNotFound
	}

	return updatedUser, err
}

// DeleteUser deletes a user by ID.
func (i *UserInteractor) DeleteUser(ctx context.Context, id int) error {
	err := i.userRepo.Delete(ctx, id)
	if err != nil && errors.Is(err, domainUser.ErrUserNotFound) {
		return domainUser.ErrUserNotFound
	}

	return err
}

// GetAllUsers retrieves all users.
func (i *UserInteractor) GetAllUsers(ctx context.Context) ([]*domainUser.User, error) {
	return i.userRepo.GetAll(ctx)
}
