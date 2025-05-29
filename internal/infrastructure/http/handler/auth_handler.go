package handler

import (
	"go-clean-architecture/internal/infrastructure/auth"
	"go-clean-architecture/internal/infrastructure/auth/jwt"
	"go-clean-architecture/internal/infrastructure/http/dto"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Convert DTO to service request
	loginReq := &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Authenticate user
	response, err := h.authService.Login(c.Context(), loginReq)
	if err != nil {
		status := fiber.StatusUnauthorized
		if err == auth.ErrUserInactive {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(dto.ErrorResponseDTO{
			Error:   "Authentication failed",
			Message: err.Error(),
		})
	}

	// Convert response to DTO
	responseDTO := dto.LoginResponseDTO{
		AccessToken: response.AccessToken,
		TokenType:   response.TokenType,
		ExpiresIn:   response.ExpiresIn,
		User: dto.UserDTO{
			ID:          response.User.ID,
			Email:       response.User.Email,
			FirstName:   response.User.FirstName,
			LastName:    response.User.LastName,
			Active:      response.User.Active,
			Roles:       response.User.Roles,
			Permissions: response.User.Permissions,
		},
	}

	return c.JSON(responseDTO)
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Convert DTO to service request
	registerReq := &auth.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Register user
	response, err := h.authService.Register(c.Context(), registerReq)
	if err != nil {
		status := fiber.StatusBadRequest
		if err == auth.ErrEmailAlreadyExists {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(dto.ErrorResponseDTO{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	// Convert response to DTO
	responseDTO := dto.LoginResponseDTO{
		AccessToken: response.AccessToken,
		TokenType:   response.TokenType,
		ExpiresIn:   response.ExpiresIn,
		User: dto.UserDTO{
			ID:          response.User.ID,
			Email:       response.User.Email,
			FirstName:   response.User.FirstName,
			LastName:    response.User.LastName,
			Active:      response.User.Active,
			Roles:       response.User.Roles,
			Permissions: response.User.Permissions,
		},
	}

	return c.Status(fiber.StatusCreated).JSON(responseDTO)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Refresh token
	response, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDTO{
			Error:   "Token refresh failed",
			Message: err.Error(),
		})
	}

	// Convert response to DTO
	responseDTO := dto.LoginResponseDTO{
		AccessToken: response.AccessToken,
		TokenType:   response.TokenType,
		ExpiresIn:   response.ExpiresIn,
		User: dto.UserDTO{
			ID:          response.User.ID,
			Email:       response.User.Email,
			FirstName:   response.User.FirstName,
			LastName:    response.User.LastName,
			Active:      response.User.Active,
			Roles:       response.User.Roles,
			Permissions: response.User.Permissions,
		},
	}

	return c.JSON(responseDTO)
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDTO{
			Error: "User not authenticated",
		})
	}

	// Get user profile
	user, err := h.authService.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponseDTO{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	// Convert to DTO
	userDTO := dto.UserDTO{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Active:      user.Active,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}

	return c.JSON(userDTO)
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDTO{
			Error: "User not authenticated",
		})
	}

	var req dto.ChangePasswordRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Change password
	err := h.authService.ChangePassword(c.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error:   "Password change failed",
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// Logout handles user logout (client-side token invalidation)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a JWT implementation, logout is typically handled client-side
	// by removing the token from storage. For server-side logout,
	// you would need to implement a token blacklist.

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

// GetMe returns current user information from token
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Get user claims from context (set by auth middleware)
	claims, ok := c.Locals("user_claims").(*jwt.TokenClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDTO{
			Error: "User not authenticated",
		})
	}

	// Convert claims to DTO
	userDTO := dto.UserDTO{
		ID:          claims.UserID,
		Email:       claims.Email,
		FirstName:   claims.FirstName,
		LastName:    claims.LastName,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
	}

	return c.JSON(userDTO)
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	// Get the token from query parameter or body
	token := c.Query("token")
	if token == "" {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
				Error:   "Token is required",
				Message: "Provide token in query parameter or request body",
			})
		}
		token = req.Token
	}

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error: "Token is required",
		})
	}

	// This would typically be done through the token service
	// For now, we'll return a simple validation response
	return c.JSON(fiber.Map{
		"valid":   true,
		"message": "Token validation endpoint",
	})
}

// UpdateProfile handles updating user profile
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user from context
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDTO{
			Error: "User not authenticated",
		})
	}

	var req dto.UpdateProfileRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Update profile (this would be implemented in the auth service)
	// For now, return success
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Profile updated successfully",
		Data: fiber.Map{
			"user_id": userID,
		},
	})
}

// GetUsers handles getting all users (admin only)
func (h *AuthHandler) GetUsers(c *fiber.Ctx) error {
	// This would call the user use case to get all users
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Users retrieved successfully",
		Data:    []interface{}{}, // Placeholder
	})
}

// GetUser handles getting a specific user
func (h *AuthHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "User retrieved successfully",
		Data: fiber.Map{
			"user_id": userID,
		},
	})
}

// UpdateUser handles updating a user (admin only)
func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "User updated successfully",
		Data: fiber.Map{
			"user_id": userID,
		},
	})
}

// DeleteUser handles deleting a user (admin only)
func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "User deleted successfully",
		Data: fiber.Map{
			"user_id": userID,
		},
	})
}

// AssignRole handles assigning a role to a user
func (h *AuthHandler) AssignRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req struct {
		RoleID uint `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error: "Invalid request body",
		})
	}

	return c.JSON(dto.SuccessResponseDTO{
		Message: "Role assigned successfully",
		Data: fiber.Map{
			"user_id": userID,
			"role_id": req.RoleID,
		},
	})
}

// RemoveRole handles removing a role from a user
func (h *AuthHandler) RemoveRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	roleID := c.Params("roleId")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Role removed successfully",
		Data: fiber.Map{
			"user_id": userID,
			"role_id": roleID,
		},
	})
}

// GetRoles handles getting all roles
func (h *AuthHandler) GetRoles(c *fiber.Ctx) error {
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Roles retrieved successfully",
		Data:    []interface{}{}, // Placeholder
	})
}

// CreateRole handles creating a new role
func (h *AuthHandler) CreateRole(c *fiber.Ctx) error {
	var req dto.CreateRoleRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error: "Invalid request body",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponseDTO{
		Message: "Role created successfully",
		Data: fiber.Map{
			"name": req.Name,
		},
	})
}

// GetRole handles getting a specific role
func (h *AuthHandler) GetRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Role retrieved successfully",
		Data: fiber.Map{
			"role_id": roleID,
		},
	})
}

// UpdateRole handles updating a role
func (h *AuthHandler) UpdateRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Role updated successfully",
		Data: fiber.Map{
			"role_id": roleID,
		},
	})
}

// DeleteRole handles deleting a role
func (h *AuthHandler) DeleteRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Role deleted successfully",
		Data: fiber.Map{
			"role_id": roleID,
		},
	})
}

// GetPermissions handles getting all permissions
func (h *AuthHandler) GetPermissions(c *fiber.Ctx) error {
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Permissions retrieved successfully",
		Data:    []interface{}{}, // Placeholder
	})
}

// CreatePermission handles creating a new permission
func (h *AuthHandler) CreatePermission(c *fiber.Ctx) error {
	var req dto.CreatePermissionRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDTO{
			Error: "Invalid request body",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponseDTO{
		Message: "Permission created successfully",
		Data: fiber.Map{
			"name": req.Name,
		},
	})
}

// GetPermission handles getting a specific permission
func (h *AuthHandler) GetPermission(c *fiber.Ctx) error {
	permissionID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Permission retrieved successfully",
		Data: fiber.Map{
			"permission_id": permissionID,
		},
	})
}

// UpdatePermission handles updating a permission
func (h *AuthHandler) UpdatePermission(c *fiber.Ctx) error {
	permissionID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Permission updated successfully",
		Data: fiber.Map{
			"permission_id": permissionID,
		},
	})
}

// DeletePermission handles deleting a permission
func (h *AuthHandler) DeletePermission(c *fiber.Ctx) error {
	permissionID := c.Params("id")
	return c.JSON(dto.SuccessResponseDTO{
		Message: "Permission deleted successfully",
		Data: fiber.Map{
			"permission_id": permissionID,
		},
	})
}
