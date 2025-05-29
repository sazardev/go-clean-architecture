package usecase

import (
	"context"
	"fmt"
	"strings"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"
)

// PermissionUseCase handles permission-related business logic
type PermissionUseCase struct {
	permissionRepo repository.PermissionRepository
}

// NewPermissionUseCase creates a new permission use case
func NewPermissionUseCase(permissionRepo repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{
		permissionRepo: permissionRepo,
	}
}

// CreatePermission creates a new permission
func (uc *PermissionUseCase) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	// Validate permission data
	if err := uc.validatePermission(permission); err != nil {
		return fmt.Errorf("permission validation failed: %w", err)
	}

	// Check if permission already exists
	_, err := uc.permissionRepo.GetByName(ctx, permission.Name)
	if err == nil {
		return fmt.Errorf("permission with name '%s' already exists", permission.Name)
	}

	// Create permission
	if err := uc.permissionRepo.Create(ctx, permission); err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

// GetPermissionByID retrieves a permission by ID
func (uc *PermissionUseCase) GetPermissionByID(ctx context.Context, id uint) (*entity.Permission, error) {
	permission, err := uc.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return permission, nil
}

// GetPermissionByName retrieves a permission by name
func (uc *PermissionUseCase) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	permission, err := uc.permissionRepo.GetByName(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return permission, nil
}

// GetAllPermissions retrieves all permissions with pagination
func (uc *PermissionUseCase) GetAllPermissions(ctx context.Context, offset, limit int) ([]*entity.Permission, error) {
	permissions, err := uc.permissionRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return permissions, nil
}

// GetPermissionsByResource retrieves permissions by resource
func (uc *PermissionUseCase) GetPermissionsByResource(ctx context.Context, resource string) ([]*entity.Permission, error) {
	permissions, err := uc.permissionRepo.GetByResource(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource: %w", err)
	}

	return permissions, nil
}

// UpdatePermission updates an existing permission
func (uc *PermissionUseCase) UpdatePermission(ctx context.Context, permission *entity.Permission) error {
	// Validate permission data
	if err := uc.validatePermission(permission); err != nil {
		return fmt.Errorf("permission validation failed: %w", err)
	}

	// Check if permission exists
	_, err := uc.permissionRepo.GetByID(ctx, permission.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("permission not found")
		}
		return fmt.Errorf("failed to check existing permission: %w", err)
	}

	// Check if name is already taken by another permission
	nameExists, err := uc.permissionRepo.GetByName(ctx, permission.Name)
	if err == nil && nameExists.ID != permission.ID {
		return fmt.Errorf("permission name '%s' is already taken", permission.Name)
	}

	// Update permission
	if err := uc.permissionRepo.Update(ctx, permission); err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}

// DeletePermission deletes a permission
func (uc *PermissionUseCase) DeletePermission(ctx context.Context, id uint) error {
	// Check if permission exists
	_, err := uc.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("permission not found")
		}
		return fmt.Errorf("failed to check existing permission: %w", err)
	}

	// Delete permission
	if err := uc.permissionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

// ActivatePermission activates a permission
func (uc *PermissionUseCase) ActivatePermission(ctx context.Context, id uint) error {
	// Check if permission exists
	permission, err := uc.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("permission not found")
		}
		return fmt.Errorf("failed to check existing permission: %w", err)
	}

	// Check if already active
	if permission.Active {
		return fmt.Errorf("permission is already active")
	}

	// Activate permission
	if err := uc.permissionRepo.ActivatePermission(ctx, id); err != nil {
		return fmt.Errorf("failed to activate permission: %w", err)
	}

	return nil
}

// DeactivatePermission deactivates a permission
func (uc *PermissionUseCase) DeactivatePermission(ctx context.Context, id uint) error {
	// Check if permission exists
	permission, err := uc.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("permission not found")
		}
		return fmt.Errorf("failed to check existing permission: %w", err)
	}

	// Check if already inactive
	if !permission.Active {
		return fmt.Errorf("permission is already inactive")
	}

	// Deactivate permission
	if err := uc.permissionRepo.DeactivatePermission(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate permission: %w", err)
	}

	return nil
}

// GetActivePermissions retrieves all active permissions
func (uc *PermissionUseCase) GetActivePermissions(ctx context.Context, offset, limit int) ([]*entity.Permission, error) {
	permissions, err := uc.permissionRepo.GetActivePermissions(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get active permissions: %w", err)
	}

	return permissions, nil
}

// BulkCreatePermissions creates multiple permissions
func (uc *PermissionUseCase) BulkCreatePermissions(ctx context.Context, permissions []*entity.Permission) error {
	// Validate all permissions
	for _, permission := range permissions {
		if err := uc.validatePermission(permission); err != nil {
			return fmt.Errorf("validation failed for permission '%s': %w", permission.Name, err)
		}
	}

	// Create permissions
	if err := uc.permissionRepo.BulkCreate(ctx, permissions); err != nil {
		return fmt.Errorf("failed to bulk create permissions: %w", err)
	}

	return nil
}

// CountPermissions returns the total count of permissions
func (uc *PermissionUseCase) CountPermissions(ctx context.Context) (int64, error) {
	count, err := uc.permissionRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count permissions: %w", err)
	}

	return count, nil
}

// validatePermission validates permission data
func (uc *PermissionUseCase) validatePermission(permission *entity.Permission) error {
	if permission == nil {
		return fmt.Errorf("permission cannot be nil")
	}

	if strings.TrimSpace(permission.Name) == "" {
		return fmt.Errorf("permission name is required")
	}

	if strings.TrimSpace(permission.Resource) == "" {
		return fmt.Errorf("permission resource is required")
	}

	if strings.TrimSpace(permission.Action) == "" {
		return fmt.Errorf("permission action is required")
	}

	// Validate resource and action format
	if !isValidResourceAction(permission.Resource) {
		return fmt.Errorf("invalid resource format")
	}

	if !isValidResourceAction(permission.Action) {
		return fmt.Errorf("invalid action format")
	}

	return nil
}

// isValidResourceAction validates resource/action format
func isValidResourceAction(value string) bool {
	// Allow alphanumeric characters, underscores, and hyphens
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	return len(value) > 0
}
