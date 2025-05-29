package rbac

import (
	"context"
	"go-clean-architecture/internal/domain/entity"
)

// PolicyManager handles RBAC policy management
type PolicyManager struct {
	enforcer *Enforcer
}

// NewPolicyManager creates a new policy manager
func NewPolicyManager(enforcer *Enforcer) *PolicyManager {
	return &PolicyManager{
		enforcer: enforcer,
	}
}

// InitializeDefaultPolicies sets up default RBAC policies
func (pm *PolicyManager) InitializeDefaultPolicies(ctx context.Context) error {
	// Default permissions for employees resource
	employeePermissions := []Permission{
		{Resource: "employees", Action: "read"},
		{Resource: "employees", Action: "create"},
		{Resource: "employees", Action: "update"},
		{Resource: "employees", Action: "delete"},
	}

	// Default permissions for users resource
	userPermissions := []Permission{
		{Resource: "users", Action: "read"},
		{Resource: "users", Action: "create"},
		{Resource: "users", Action: "update"},
		{Resource: "users", Action: "delete"},
		{Resource: "users", Action: "assign_role"},
		{Resource: "users", Action: "remove_role"},
	}

	// Default permissions for roles resource
	rolePermissions := []Permission{
		{Resource: "roles", Action: "read"},
		{Resource: "roles", Action: "create"},
		{Resource: "roles", Action: "update"},
		{Resource: "roles", Action: "delete"},
	}

	// Super Admin - full access
	for _, perm := range append(append(employeePermissions, userPermissions...), rolePermissions...) {
		if err := pm.enforcer.AddPolicy("super_admin", perm.Resource, perm.Action); err != nil {
			// Policy might already exist, continue
		}
	}

	// Admin - most access except user role assignment
	adminPermissions := append(employeePermissions, userPermissions[:4]...)
	adminPermissions = append(adminPermissions, rolePermissions[:3]...) // No role deletion
	for _, perm := range adminPermissions {
		if err := pm.enforcer.AddPolicy("admin", perm.Resource, perm.Action); err != nil {
			// Policy might already exist, continue
		}
	}

	// HR Manager - employee management + limited user management
	hrManagerPermissions := append(employeePermissions, userPermissions[:3]...) // No user deletion
	hrManagerPermissions = append(hrManagerPermissions, Permission{Resource: "roles", Action: "read"})
	for _, perm := range hrManagerPermissions {
		if err := pm.enforcer.AddPolicy("hr_manager", perm.Resource, perm.Action); err != nil {
			// Policy might already exist, continue
		}
	}

	// HR Specialist - employee management only
	for _, perm := range employeePermissions {
		if err := pm.enforcer.AddPolicy("hr_specialist", perm.Resource, perm.Action); err != nil {
			// Policy might already exist, continue
		}
	}

	// Employee - read only access to employees
	if err := pm.enforcer.AddPolicy("employee", "employees", "read"); err != nil {
		// Policy might already exist, continue
	}

	return nil
}

// Permission represents a resource-action pair
type Permission struct {
	Resource string
	Action   string
}

// AssignRoleToUser assigns a role to a user
func (pm *PolicyManager) AssignRoleToUser(userEmail, roleName string) error {
	return pm.enforcer.AddRoleForUser(userEmail, roleName)
}

// RemoveRoleFromUser removes a role from a user
func (pm *PolicyManager) RemoveRoleFromUser(userEmail, roleName string) error {
	return pm.enforcer.DeleteRoleForUser(userEmail, roleName)
}

// GrantPermissionToRole grants a permission to a role
func (pm *PolicyManager) GrantPermissionToRole(roleName, resource, action string) error {
	return pm.enforcer.AddPolicy(roleName, resource, action)
}

// RevokePermissionFromRole revokes a permission from a role
func (pm *PolicyManager) RevokePermissionFromRole(roleName, resource, action string) error {
	return pm.enforcer.RemovePolicy(roleName, resource, action)
}

// GetUserRoles returns all roles for a user
func (pm *PolicyManager) GetUserRoles(userEmail string) ([]string, error) {
	return pm.enforcer.GetRolesForUser(userEmail)
}

// GetRoleUsers returns all users with a specific role
func (pm *PolicyManager) GetRoleUsers(roleName string) ([]string, error) {
	return pm.enforcer.GetUsersForRole(roleName)
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (pm *PolicyManager) CheckPermission(userEmail, resource, action string) (bool, error) {
	return pm.enforcer.Enforce(userEmail, resource, action)
}

// CheckPermissionWithRoles checks if any of the user's roles has permission
func (pm *PolicyManager) CheckPermissionWithRoles(roles []string, resource, action string) (bool, error) {
	return pm.enforcer.EnforceWithRoles(roles, resource, action)
}

// SyncUserPolicies synchronizes user policies with database entities
func (pm *PolicyManager) SyncUserPolicies(user *entity.User) error {
	userEmail := user.Email

	// Remove all existing roles for the user
	existingRoles, err := pm.enforcer.GetRolesForUser(userEmail)
	if err != nil {
		return err
	}

	for _, role := range existingRoles {
		if err := pm.enforcer.DeleteRoleForUser(userEmail, role); err != nil {
			return err
		}
	}

	// Add current roles
	for _, role := range user.Roles {
		if err := pm.enforcer.AddRoleForUser(userEmail, role.Name); err != nil {
			return err
		}
	}

	return nil
}
