package entity

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	Resource    string         `gorm:"not null" json:"resource"` // e.g., "employees", "users", "roles"
	Action      string         `gorm:"not null" json:"action"`   // e.g., "read", "write", "delete"
	Active      bool           `gorm:"default:true" json:"active"`
	Roles       []Role         `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetCasbinFormat returns the permission in Casbin format (resource:action)
func (p *Permission) GetCasbinFormat() string {
	return p.Resource + ":" + p.Action
}

// PermissionType represents common permission patterns
type PermissionType struct {
	Name        string
	Description string
	Resource    string
	Action      string
}

// Common permissions that can be used across the system
var (
	// Employee permissions
	EmployeeRead   = PermissionType{Name: "employee.read", Description: "Read employee data", Resource: "employees", Action: "read"}
	EmployeeWrite  = PermissionType{Name: "employee.write", Description: "Create and update employees", Resource: "employees", Action: "write"}
	EmployeeDelete = PermissionType{Name: "employee.delete", Description: "Delete employees", Resource: "employees", Action: "delete"}

	// User permissions
	UserRead   = PermissionType{Name: "user.read", Description: "Read user data", Resource: "users", Action: "read"}
	UserWrite  = PermissionType{Name: "user.write", Description: "Create and update users", Resource: "users", Action: "write"}
	UserDelete = PermissionType{Name: "user.delete", Description: "Delete users", Resource: "users", Action: "delete"}

	// Role permissions
	RoleRead   = PermissionType{Name: "role.read", Description: "Read role data", Resource: "roles", Action: "read"}
	RoleWrite  = PermissionType{Name: "role.write", Description: "Create and update roles", Resource: "roles", Action: "write"}
	RoleDelete = PermissionType{Name: "role.delete", Description: "Delete roles", Resource: "roles", Action: "delete"}

	// Permission permissions
	PermissionRead   = PermissionType{Name: "permission.read", Description: "Read permission data", Resource: "permissions", Action: "read"}
	PermissionWrite  = PermissionType{Name: "permission.write", Description: "Create and update permissions", Resource: "permissions", Action: "write"}
	PermissionDelete = PermissionType{Name: "permission.delete", Description: "Delete permissions", Resource: "permissions", Action: "delete"}

	// System permissions
	SystemAdmin = PermissionType{Name: "system.admin", Description: "Full system administration", Resource: "system", Action: "admin"}
)

// GetAllPermissionTypes returns all predefined permission types
func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		EmployeeRead, EmployeeWrite, EmployeeDelete,
		UserRead, UserWrite, UserDelete,
		RoleRead, RoleWrite, RoleDelete,
		PermissionRead, PermissionWrite, PermissionDelete,
		SystemAdmin,
	}
}
