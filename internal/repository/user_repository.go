package repository

import (
	"context"

	"github.com/tetsuyaohta/go-backend-boilerplate/ent"
	domainUser "github.com/tetsuyaohta/go-backend-boilerplate/internal/domain/user"
)

// UserRepository implements the user.Repository interface using ent.
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository.
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id int) (*domainUser.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapEntUserToDomainUser(u), nil
}

// Create stores a new user.
func (r *UserRepository) Create(ctx context.Context, user *domainUser.User) (*domainUser.User, error) {
	created, err := r.client.User.
		Create().
		SetName(user.Name).
		SetEmail(user.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapEntUserToDomainUser(created), nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *domainUser.User) (*domainUser.User, error) {
	updated, err := r.client.User.
		UpdateOneID(user.ID).
		SetName(user.Name).
		SetEmail(user.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapEntUserToDomainUser(updated), nil
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}

// GetAll retrieves all users.
func (r *UserRepository) GetAll(ctx context.Context) ([]*domainUser.User, error) {
	users, err := r.client.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*domainUser.User, 0, len(users))
	for _, u := range users {
		result = append(result, mapEntUserToDomainUser(u))
	}

	return result, nil
}

// mapEntUserToDomainUser converts an ent User entity to a domain User entity.
func mapEntUserToDomainUser(u *ent.User) *domainUser.User {
	return &domainUser.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
