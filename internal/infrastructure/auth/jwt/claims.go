package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims extends the standard JWT claims with additional user information
type CustomClaims struct {
	UserID      uint     `json:"user_id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// HasRole checks if the claims contain a specific role
func (c *CustomClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission checks if the claims contain a specific permission
func (c *CustomClaims) HasPermission(permission string) bool {
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the claims contain any of the specified roles
func (c *CustomClaims) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if the claims contain any of the specified permissions
func (c *CustomClaims) HasAnyPermission(permissions ...string) bool {
	for _, permission := range permissions {
		if c.HasPermission(permission) {
			return true
		}
	}
	return false
}

// GetFullName returns the user's full name
func (c *CustomClaims) GetFullName() string {
	return c.FirstName + " " + c.LastName
}

// IsAdmin checks if the user has admin role
func (c *CustomClaims) IsAdmin() bool {
	return c.HasRole("admin")
}

// IsSuperAdmin checks if the user has super admin role
func (c *CustomClaims) IsSuperAdmin() bool {
	return c.HasRole("super_admin")
}
