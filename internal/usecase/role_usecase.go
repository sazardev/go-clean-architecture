package usecase

import (
	"context"
	"errors"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"
	"go-clean-architecture/internal/infrastructure/auth/rbac"
)

// RoleUseCase handles role-related business logic
type RoleUseCase struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	userRepo       repository.UserRepository
	policyManager  *rbac.PolicyManager
}

// NewRoleUseCase creates a new role use case
func NewRoleUseCase(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRepo repository.UserRepository,
	policyManager *rbac.PolicyManager,
) *RoleUseCase {
	return &RoleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
		policyManager:  policyManager,
	}
}

// CreateRole creates a new role
func (uc *RoleUseCase) CreateRole(ctx context.Context, name, description string, active bool) (*entity.Role, error) {
	// Check if role already exists
	existingRole, err := uc.roleRepo.GetByName(ctx, name)
	if err == nil && existingRole != nil {
		return nil, errors.New("role already exists")
	}

	// Create role
	role := &entity.Role{
		Name:        name,
		Description: description,
		Active:      active,
	}

	// Save role
	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID retrieves a role by ID
func (uc *RoleUseCase) GetRoleByID(ctx context.Context, id uint) (*entity.Role, error) {
	return uc.roleRepo.GetByIDWithPermissions(ctx, id)
}

// GetRoleByName retrieves a role by name
func (uc *RoleUseCase) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	return uc.roleRepo.GetByNameWithPermissions(ctx, name)
}

// GetAllRoles retrieves all roles
func (uc *RoleUseCase) GetAllRoles(ctx context.Context) ([]*entity.Role, error) {
	return uc.roleRepo.List(ctx, 0, 1000) // Get first 1000 roles
}

// UpdateRole updates a role
func (uc *RoleUseCase) UpdateRole(ctx context.Context, role *entity.Role) error {
	return uc.roleRepo.Update(ctx, role)
}

// DeleteRole deletes a role
func (uc *RoleUseCase) DeleteRole(ctx context.Context, id uint) error { // Get role first
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if role is being used by any users
	users, err := uc.policyManager.GetRoleUsers(role.Name)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		return errors.New("cannot delete role that is assigned to users")
	}

	// Remove from RBAC
	if err := uc.policyManager.RemoveRoleFromUser("", role.Name); err != nil {
		// Log error but continue
	}

	// Delete role
	return uc.roleRepo.Delete(ctx, id)
}

// AssignPermissionToRole assigns a permission to a role
func (uc *RoleUseCase) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
	// Get role and permission
	role, err := uc.roleRepo.GetByIDWithPermissions(ctx, roleID)
	if err != nil {
		return err
	}

	permission, err := uc.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return err
	}

	// Check if role already has the permission
	for _, rolePermission := range role.Permissions {
		if rolePermission.ID == permissionID {
			return errors.New("role already has this permission")
		}
	}

	// Assign permission in database
	if err := uc.roleRepo.AssignPermission(ctx, roleID, permissionID); err != nil {
		return err
	}

	// Grant permission in RBAC
	if err := uc.policyManager.GrantPermissionToRole(role.Name, permission.Resource, permission.Action); err != nil {
		return err
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (uc *RoleUseCase) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	// Get role and permission
	role, err := uc.roleRepo.GetByIDWithPermissions(ctx, roleID)
	if err != nil {
		return err
	}

	permission, err := uc.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return err
	}

	// Remove permission from database
	if err := uc.roleRepo.RemovePermission(ctx, roleID, permissionID); err != nil {
		return err
	}

	// Revoke permission from RBAC
	if err := uc.policyManager.RevokePermissionFromRole(role.Name, permission.Resource, permission.Action); err != nil {
		return err
	}

	return nil
}

// GetRolePermissions gets all permissions for a role
func (uc *RoleUseCase) GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error) {
	role, err := uc.roleRepo.GetByIDWithPermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]*entity.Permission, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissions[i] = &permission
	}

	return permissions, nil
}

// InitializeDefaultRoles creates default roles if they don't exist
func (uc *RoleUseCase) InitializeDefaultRoles(ctx context.Context) error {
	defaultRoles := []struct {
		Name        string
		Description string
	}{
		{"super_admin", "Super Administrator with full access"},
		{"admin", "Administrator with management access"},
		{"hr_manager", "HR Manager with employee management access"},
		{"hr_specialist", "HR Specialist with limited employee access"},
		{"employee", "Regular employee with read-only access"},
	}

	for _, roleData := range defaultRoles { // Check if role exists
		_, err := uc.roleRepo.GetByName(ctx, roleData.Name)
		if err != nil {
			// Role doesn't exist, create it
			role := &entity.Role{
				Name:        roleData.Name,
				Description: roleData.Description,
				Active:      true,
			}

			if err := uc.roleRepo.Create(ctx, role); err != nil {
				return err
			}
		}
	}

	return nil
}
