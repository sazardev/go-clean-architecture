package repository

import (
	"context"

	"go-clean-architecture/internal/domain/entity"
)

type PermissionRepository interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *entity.Permission) error

	// GetByID retrieves a permission by ID
	GetByID(ctx context.Context, id uint) (*entity.Permission, error)

	// GetByName retrieves a permission by name
	GetByName(ctx context.Context, name string) (*entity.Permission, error)

	// Update updates an existing permission
	Update(ctx context.Context, permission *entity.Permission) error

	// Delete soft deletes a permission
	Delete(ctx context.Context, id uint) error

	// List retrieves all permissions with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.Permission, error)

	// Count returns the total count of permissions
	Count(ctx context.Context) (int64, error)

	// ExistsByName checks if a permission with the given name exists
	ExistsByName(ctx context.Context, name string) (bool, error)

	// GetActivePermissions retrieves all active permissions
	GetActivePermissions(ctx context.Context, offset, limit int) ([]*entity.Permission, error)

	// ActivatePermission activates a permission
	ActivatePermission(ctx context.Context, id uint) error

	// DeactivatePermission deactivates a permission
	DeactivatePermission(ctx context.Context, id uint) error

	// GetByResource retrieves permissions by resource
	GetByResource(ctx context.Context, resource string) ([]*entity.Permission, error)

	// GetByResourceAndAction retrieves a permission by resource and action
	GetByResourceAndAction(ctx context.Context, resource, action string) (*entity.Permission, error)

	// BulkCreate creates multiple permissions
	BulkCreate(ctx context.Context, permissions []*entity.Permission) error

	// GetRolesWithPermission retrieves all roles that have a specific permission
	GetRolesWithPermission(ctx context.Context, permissionID uint) ([]*entity.Role, error)
}
