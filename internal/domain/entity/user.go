package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FirstName string         `gorm:"not null" json:"first_name"`
	LastName  string         `gorm:"not null" json:"last_name"`
	Active    bool           `gorm:"default:true" json:"active"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SetPassword encrypts and sets the user password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the given password with the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// GetPermissions returns all permissions assigned to the user through their roles
func (u *User) GetPermissions() []Permission {
	var permissions []Permission
	permissionMap := make(map[uint]bool)

	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if !permissionMap[permission.ID] {
				permissions = append(permissions, permission)
				permissionMap[permission.ID] = true
			}
		}
	}

	return permissions
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permissionName string) bool {
	permissions := u.GetPermissions()
	for _, permission := range permissions {
		if permission.Name == permissionName {
			return true
		}
	}
	return false
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
