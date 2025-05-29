package repository

import (
	"context"
	"errors"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"

	"gorm.io/gorm"
)

// permissionRepository implements repository.PermissionRepository
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

// Create creates a new permission
func (r *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	result := r.db.WithContext(ctx).Create(permission)
	return result.Error
}

// GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*entity.Permission, error) {
	var permission entity.Permission
	result := r.db.WithContext(ctx).First(&permission, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, result.Error
	}
	return &permission, nil
}

// GetByName retrieves a permission by name
func (r *permissionRepository) GetByName(ctx context.Context, name string) (*entity.Permission, error) {
	var permission entity.Permission
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, result.Error
	}
	return &permission, nil
}

// List retrieves all permissions with pagination
func (r *permissionRepository) List(ctx context.Context, offset, limit int) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

// Count returns the total count of permissions
func (r *permissionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entity.Permission{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// GetByResource retrieves permissions by resource
func (r *permissionRepository) GetByResource(ctx context.Context, resource string) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	result := r.db.WithContext(ctx).Where("resource = ?", resource).Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

// GetByResourceAndAction retrieves a permission by resource and action
func (r *permissionRepository) GetByResourceAndAction(ctx context.Context, resource, action string) (*entity.Permission, error) {
	var permission entity.Permission
	result := r.db.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, result.Error
	}
	return &permission, nil
}

// Update updates an existing permission
func (r *permissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	result := r.db.WithContext(ctx).Save(permission)
	return result.Error
}

// Delete soft deletes a permission
func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Permission{}, id)
	return result.Error
}

// ExistsByName checks if a permission with the given name exists
func (r *permissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entity.Permission{}).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// GetActivePermissions retrieves all active permissions with pagination
func (r *permissionRepository) GetActivePermissions(ctx context.Context, offset, limit int) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	result := r.db.WithContext(ctx).Where("active = ?", true).Offset(offset).Limit(limit).Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

// ActivatePermission activates a permission
func (r *permissionRepository) ActivatePermission(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&entity.Permission{}).Where("id = ?", id).Update("active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

// DeactivatePermission deactivates a permission
func (r *permissionRepository) DeactivatePermission(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&entity.Permission{}).Where("id = ?", id).Update("active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

// BulkCreate creates multiple permissions in a transaction
func (r *permissionRepository) BulkCreate(ctx context.Context, permissions []*entity.Permission) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, permission := range permissions {
			if err := tx.Create(permission).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetRolesWithPermission retrieves all roles that have a specific permission
func (r *permissionRepository) GetRolesWithPermission(ctx context.Context, permissionID uint) ([]*entity.Role, error) {
	var roles []*entity.Role
	result := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Where("role_permissions.permission_id = ?", permissionID).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
