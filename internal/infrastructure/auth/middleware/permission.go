package middleware

import (
	"go-clean-architecture/internal/infrastructure/auth/rbac"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission creates a middleware that checks if the user has a specific permission
func RequirePermission(policyManager *rbac.PolicyManager, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user roles from context (set by auth middleware)
		roles, ok := c.Locals("user_roles").([]string)
		if !ok || len(roles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: No roles assigned",
			})
		}

		// Check if any role has the required permission
		hasPermission, err := policyManager.CheckPermissionWithRoles(roles, resource, action)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions",
			})
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user roles from context (set by auth middleware)
		roles, ok := c.Locals("user_roles").([]string)
		if !ok || len(roles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: No roles assigned",
			})
		}

		// Check if user has the required role
		hasRole := false
		for _, role := range roles {
			if role == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: Required role not found",
			})
		}

		return c.Next()
	}
}

// RequireAnyRole creates a middleware that checks if the user has any of the specified roles
func RequireAnyRole(roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user roles from context (set by auth middleware)
		userRoles, ok := c.Locals("user_roles").([]string)
		if !ok || len(userRoles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: No roles assigned",
			})
		}

		// Check if user has any of the required roles
		hasAnyRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range roleNames {
				if userRole == requiredRole {
					hasAnyRole = true
					break
				}
			}
			if hasAnyRole {
				break
			}
		}

		if !hasAnyRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: None of the required roles found",
			})
		}

		return c.Next()
	}
}

// RequireAnyPermission creates a middleware that checks if the user has any of the specified permissions
func RequireAnyPermission(policyManager *rbac.PolicyManager, permissions ...Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user roles from context (set by auth middleware)
		roles, ok := c.Locals("user_roles").([]string)
		if !ok || len(roles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: No roles assigned",
			})
		}

		// Check if any role has any of the required permissions
		hasAnyPermission := false
		for _, perm := range permissions {
			hasPermission, err := policyManager.CheckPermissionWithRoles(roles, perm.Resource, perm.Action)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to check permissions",
				})
			}
			if hasPermission {
				hasAnyPermission = true
				break
			}
		}

		if !hasAnyPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// AdminOnly creates a middleware that only allows admin users
func AdminOnly() fiber.Handler {
	return RequireAnyRole("admin", "super_admin")
}

// SuperAdminOnly creates a middleware that only allows super admin users
func SuperAdminOnly() fiber.Handler {
	return RequireRole("super_admin")
}

// HROnly creates a middleware that only allows HR personnel
func HROnly() fiber.Handler {
	return RequireAnyRole("admin", "super_admin", "hr_manager", "hr_specialist")
}

// Permission represents a resource-action pair for permission checking
type Permission struct {
	Resource string
	Action   string
}

// Common permission constants
var (
	ReadEmployees   = Permission{Resource: "employees", Action: "read"}
	CreateEmployees = Permission{Resource: "employees", Action: "create"}
	UpdateEmployees = Permission{Resource: "employees", Action: "update"}
	DeleteEmployees = Permission{Resource: "employees", Action: "delete"}

	ReadUsers   = Permission{Resource: "users", Action: "read"}
	CreateUsers = Permission{Resource: "users", Action: "create"}
	UpdateUsers = Permission{Resource: "users", Action: "update"}
	DeleteUsers = Permission{Resource: "users", Action: "delete"}

	ReadRoles   = Permission{Resource: "roles", Action: "read"}
	CreateRoles = Permission{Resource: "roles", Action: "create"}
	UpdateRoles = Permission{Resource: "roles", Action: "update"}
	DeleteRoles = Permission{Resource: "roles", Action: "delete"}
)
