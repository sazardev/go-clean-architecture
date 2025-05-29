package repository

import (
	"context"

	"go-clean-architecture/internal/domain/entity"
)

type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByEmailWithRoles retrieves a user by email with their roles and permissions
	GetByEmailWithRoles(ctx context.Context, email string) (*entity.User, error)

	// GetByIDWithRoles retrieves a user by ID with their roles and permissions
	GetByIDWithRoles(ctx context.Context, id uint) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uint) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// ListWithRoles retrieves all users with their roles
	ListWithRoles(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// Count returns the total count of users
	Count(ctx context.Context) (int64, error)

	// AssignRole assigns a role to a user
	AssignRole(ctx context.Context, userID, roleID uint) error

	// RemoveRole removes a role from a user
	RemoveRole(ctx context.Context, userID, roleID uint) error

	// GetUserRoles retrieves all roles for a user
	GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// GetActiveUsers retrieves all active users
	GetActiveUsers(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// ActivateUser activates a user
	ActivateUser(ctx context.Context, id uint) error

	// DeactivateUser deactivates a user
	DeactivateUser(ctx context.Context, id uint) error
}
