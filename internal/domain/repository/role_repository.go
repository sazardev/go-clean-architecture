package repository

import (
	"context"

	"go-clean-architecture/internal/domain/entity"
)

type RoleRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *entity.Role) error

	// GetByID retrieves a role by ID
	GetByID(ctx context.Context, id uint) (*entity.Role, error)

	// GetByName retrieves a role by name
	GetByName(ctx context.Context, name string) (*entity.Role, error)

	// GetByNameWithPermissions retrieves a role by name with its permissions
	GetByNameWithPermissions(ctx context.Context, name string) (*entity.Role, error)

	// GetByIDWithPermissions retrieves a role by ID with its permissions
	GetByIDWithPermissions(ctx context.Context, id uint) (*entity.Role, error)

	// Update updates an existing role
	Update(ctx context.Context, role *entity.Role) error

	// Delete soft deletes a role
	Delete(ctx context.Context, id uint) error

	// List retrieves all roles with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.Role, error)

	// ListWithPermissions retrieves all roles with their permissions
	ListWithPermissions(ctx context.Context, offset, limit int) ([]*entity.Role, error)

	// Count returns the total count of roles
	Count(ctx context.Context) (int64, error)

	// AssignPermission assigns a permission to a role
	AssignPermission(ctx context.Context, roleID, permissionID uint) error

	// RemovePermission removes a permission from a role
	RemovePermission(ctx context.Context, roleID, permissionID uint) error

	// GetRolePermissions retrieves all permissions for a role
	GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error)

	// ExistsByName checks if a role with the given name exists
	ExistsByName(ctx context.Context, name string) (bool, error)

	// GetActiveRoles retrieves all active roles
	GetActiveRoles(ctx context.Context, offset, limit int) ([]*entity.Role, error)

	// ActivateRole activates a role
	ActivateRole(ctx context.Context, id uint) error

	// DeactivateRole deactivates a role
	DeactivateRole(ctx context.Context, id uint) error

	// GetUsersWithRole retrieves all users that have a specific role
	GetUsersWithRole(ctx context.Context, roleID uint) ([]*entity.User, error)
}
