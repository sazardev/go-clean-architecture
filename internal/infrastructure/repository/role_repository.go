package repository

import (
	"context"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

// Create creates a new role
func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name
func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByNameWithPermissions retrieves a role by name with its permissions
func (r *roleRepository) GetByNameWithPermissions(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name = ?", name).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByIDWithPermissions retrieves a role by ID with its permissions
func (r *roleRepository) GetByIDWithPermissions(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update updates an existing role
func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete soft deletes a role
func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Role{}, id).Error
}

// List retrieves all roles with pagination
func (r *roleRepository) List(ctx context.Context, offset, limit int) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&roles).Error
	return roles, err
}

// ListWithPermissions retrieves all roles with their permissions
func (r *roleRepository) ListWithPermissions(ctx context.Context, offset, limit int) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Offset(offset).
		Limit(limit).
		Find(&roles).Error
	return roles, err
}

// Count returns the total count of roles
func (r *roleRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Role{}).Count(&count).Error
	return count, err
}

// AssignPermission assigns a permission to a role
func (r *roleRepository) AssignPermission(ctx context.Context, roleID, permissionID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		roleID, permissionID,
	).Error
}

// RemovePermission removes a permission from a role
func (r *roleRepository) RemovePermission(ctx context.Context, roleID, permissionID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?",
		roleID, permissionID,
	).Error
}

// GetRolePermissions retrieves all permissions for a role
func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// ExistsByName checks if a role with the given name exists
func (r *roleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}

// GetActiveRoles retrieves all active roles
func (r *roleRepository) GetActiveRoles(ctx context.Context, offset, limit int) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&roles).Error
	return roles, err
}

// ActivateRole activates a role
func (r *roleRepository) ActivateRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("id = ?", id).
		Update("active", true).Error
}

// DeactivateRole deactivates a role
func (r *roleRepository) DeactivateRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("id = ?", id).
		Update("active", false).Error
}

// GetUsersWithRole retrieves all users that have a specific role
func (r *roleRepository) GetUsersWithRole(ctx context.Context, roleID uint) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users).Error
	return users, err
}
