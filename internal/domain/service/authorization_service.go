package service

import (
	"fmt"

	"go-clean-architecture/internal/domain/entity"

	"github.com/casbin/casbin/v2"
)

// AuthorizationService handles authorization logic using Casbin
type AuthorizationService interface {
	// Enforce checks if a user has permission to perform an action on a resource
	Enforce(userID uint, resource, action string) (bool, error)

	// EnforceWithRole checks if a role has permission to perform an action on a resource
	EnforceWithRole(role, resource, action string) (bool, error)

	// AddPolicy adds a policy to Casbin
	AddPolicy(role, resource, action string) (bool, error)

	// RemovePolicy removes a policy from Casbin
	RemovePolicy(role, resource, action string) (bool, error)

	// AddRoleForUser adds a role for a user
	AddRoleForUser(userID uint, role string) (bool, error)

	// DeleteRoleForUser deletes a role for a user
	DeleteRoleForUser(userID uint, role string) (bool, error)

	// GetRolesForUser gets all roles for a user
	GetRolesForUser(userID uint) ([]string, error)

	// GetUsersForRole gets all users for a role
	GetUsersForRole(role string) ([]string, error)

	// GetPermissionsForUser gets all permissions for a user
	GetPermissionsForUser(userID uint) ([][]string, error)

	// GetPermissionsForRole gets all permissions for a role
	GetPermissionsForRole(role string) ([][]string, error)

	// LoadPolicy loads policy from storage
	LoadPolicy() error
	// SavePolicy saves policy to storage
	SavePolicy() error

	// SyncUserPermissions synchronizes user permissions with database
	SyncUserPermissions(user *entity.User) error
}

type authorizationService struct {
	enforcer *casbin.Enforcer
}

// NewAuthorizationService creates a new authorization service
func NewAuthorizationService(enforcer *casbin.Enforcer) AuthorizationService {
	return &authorizationService{
		enforcer: enforcer,
	}
}

// Enforce checks if a user has permission to perform an action on a resource
func (a *authorizationService) Enforce(userID uint, resource, action string) (bool, error) {
	return a.enforcer.Enforce(fmt.Sprintf("user:%d", userID), resource, action)
}

// EnforceWithRole checks if a role has permission to perform an action on a resource
func (a *authorizationService) EnforceWithRole(role, resource, action string) (bool, error) {
	return a.enforcer.Enforce(role, resource, action)
}

// AddPolicy adds a policy to Casbin
func (a *authorizationService) AddPolicy(role, resource, action string) (bool, error) {
	return a.enforcer.AddPolicy(role, resource, action)
}

// RemovePolicy removes a policy from Casbin
func (a *authorizationService) RemovePolicy(role, resource, action string) (bool, error) {
	return a.enforcer.RemovePolicy(role, resource, action)
}

// AddRoleForUser adds a role for a user
func (a *authorizationService) AddRoleForUser(userID uint, role string) (bool, error) {
	return a.enforcer.AddRoleForUser(fmt.Sprintf("user:%d", userID), role)
}

// DeleteRoleForUser deletes a role for a user
func (a *authorizationService) DeleteRoleForUser(userID uint, role string) (bool, error) {
	return a.enforcer.DeleteRoleForUser(fmt.Sprintf("user:%d", userID), role)
}

// GetRolesForUser gets all roles for a user
func (a *authorizationService) GetRolesForUser(userID uint) ([]string, error) {
	return a.enforcer.GetRolesForUser(fmt.Sprintf("user:%d", userID))
}

// GetUsersForRole gets all users for a role
func (a *authorizationService) GetUsersForRole(role string) ([]string, error) {
	return a.enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser gets all permissions for a user
func (a *authorizationService) GetPermissionsForUser(userID uint) ([][]string, error) {
	return a.enforcer.GetPermissionsForUser(fmt.Sprintf("user:%d", userID))
}

// GetPermissionsForRole gets all permissions for a role
func (a *authorizationService) GetPermissionsForRole(role string) ([][]string, error) {
	return a.enforcer.GetImplicitPermissionsForUser(role)
}

// LoadPolicy loads policy from storage
func (a *authorizationService) LoadPolicy() error {
	return a.enforcer.LoadPolicy()
}

// SavePolicy saves policy to storage
func (a *authorizationService) SavePolicy() error {
	return a.enforcer.SavePolicy()
}

// SyncUserPermissions synchronizes user permissions with database
func (a *authorizationService) SyncUserPermissions(user *entity.User) error {
	// Remove all existing roles for the user
	existingRoles, err := a.GetRolesForUser(user.ID)
	if err != nil {
		return err
	}

	for _, role := range existingRoles {
		if _, err := a.DeleteRoleForUser(user.ID, role); err != nil {
			return err
		}
	}

	// Add current roles
	for _, role := range user.Roles {
		if _, err := a.AddRoleForUser(user.ID, role.Name); err != nil {
			return err
		}
	}

	return nil
}

// SyncRolePermissions synchronizes role permissions with database
func (a *authorizationService) SyncRolePermissions(role *entity.Role) error {
	// Remove all existing permissions for the role
	existingPermissions, err := a.GetPermissionsForRole(role.Name)
	if err != nil {
		return err
	}

	for _, perm := range existingPermissions {
		if len(perm) >= 3 {
			if _, err := a.RemovePolicy(perm[0], perm[1], perm[2]); err != nil {
				return err
			}
		}
	}

	// Add current permissions
	for _, permission := range role.Permissions {
		if _, err := a.AddPolicy(role.Name, permission.Resource, permission.Action); err != nil {
			return err
		}
	}

	return nil
}
