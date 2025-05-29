package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Enforcer wraps Casbin enforcer with additional functionality
type Enforcer struct {
	enforcer *casbin.Enforcer
	adapter  *gormadapter.Adapter
}

// NewEnforcer creates a new RBAC enforcer
func NewEnforcer(db *gorm.DB, modelPath string) (*Enforcer, error) {
	// Create Casbin adapter with GORM
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// Create Casbin enforcer
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Load policy from database
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &Enforcer{
		enforcer: enforcer,
		adapter:  adapter,
	}, nil
}

// Enforce checks if a user has permission to perform an action on a resource
func (e *Enforcer) Enforce(subject, object, action string) (bool, error) {
	return e.enforcer.Enforce(subject, object, action)
}

// EnforceWithRoles checks if any of the user's roles has permission
func (e *Enforcer) EnforceWithRoles(roles []string, object, action string) (bool, error) {
	for _, role := range roles {
		allowed, err := e.enforcer.Enforce(role, object, action)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}

// AddPolicy adds a policy rule
func (e *Enforcer) AddPolicy(subject, object, action string) error {
	added, err := e.enforcer.AddPolicy(subject, object, action)
	if err != nil {
		return err
	}
	if !added {
		return errors.New("policy already exists")
	}
	return e.enforcer.SavePolicy()
}

// RemovePolicy removes a policy rule
func (e *Enforcer) RemovePolicy(subject, object, action string) error {
	removed, err := e.enforcer.RemovePolicy(subject, object, action)
	if err != nil {
		return err
	}
	if !removed {
		return errors.New("policy does not exist")
	}
	return e.enforcer.SavePolicy()
}

// AddRoleForUser assigns a role to a user
func (e *Enforcer) AddRoleForUser(user, role string) error {
	added, err := e.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return err
	}
	if !added {
		return errors.New("role assignment already exists")
	}
	return e.enforcer.SavePolicy()
}

// DeleteRoleForUser removes a role from a user
func (e *Enforcer) DeleteRoleForUser(user, role string) error {
	removed, err := e.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return err
	}
	if !removed {
		return errors.New("role assignment does not exist")
	}
	return e.enforcer.SavePolicy()
}

// GetRolesForUser gets all roles for a user
func (e *Enforcer) GetRolesForUser(user string) ([]string, error) {
	return e.enforcer.GetRolesForUser(user)
}

// GetUsersForRole gets all users with a specific role
func (e *Enforcer) GetUsersForRole(role string) ([]string, error) {
	return e.enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser gets all permissions for a user
func (e *Enforcer) GetPermissionsForUser(user string) ([][]string, error) {
	return e.enforcer.GetPermissionsForUser(user)
}

// HasRoleForUser checks if a user has a specific role
func (e *Enforcer) HasRoleForUser(user, role string) (bool, error) {
	return e.enforcer.HasRoleForUser(user, role)
}

// DeleteUser removes all policies for a user
func (e *Enforcer) DeleteUser(user string) error {
	removed, err := e.enforcer.DeleteUser(user)
	if err != nil {
		return err
	}
	if !removed {
		return errors.New("user does not exist")
	}
	return e.enforcer.SavePolicy()
}

// DeleteRole removes all policies for a role
func (e *Enforcer) DeleteRole(role string) error {
	removed, err := e.enforcer.DeleteRole(role)
	if err != nil {
		return err
	}
	if !removed {
		return errors.New("role does not exist")
	}
	return e.enforcer.SavePolicy()
}

// LoadPolicy reloads the policy from storage
func (e *Enforcer) LoadPolicy() error {
	return e.enforcer.LoadPolicy()
}

// SavePolicy saves the current policy to storage
func (e *Enforcer) SavePolicy() error {
	return e.enforcer.SavePolicy()
}

// BuildContext creates a context with the enforcer
func (e *Enforcer) BuildContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "rbac_enforcer", e)
}

// FromContext extracts the enforcer from context
func FromContext(ctx context.Context) (*Enforcer, bool) {
	enforcer, ok := ctx.Value("rbac_enforcer").(*Enforcer)
	return enforcer, ok
}
