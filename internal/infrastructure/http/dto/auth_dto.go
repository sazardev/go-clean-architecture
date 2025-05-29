package dto

// LoginRequestDTO represents a login request
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponseDTO represents a login response
type LoginResponseDTO struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   int64   `json:"expires_in"`
	User        UserDTO `json:"user"`
}

// RegisterRequestDTO represents a registration request
type RegisterRequestDTO struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}

// RefreshTokenRequestDTO represents a token refresh request
type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequestDTO represents a password change request
type ChangePasswordRequestDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// UserDTO represents user information in responses
type UserDTO struct {
	ID          uint     `json:"id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Active      bool     `json:"active"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// RoleDTO represents role information
type RoleDTO struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Active      bool            `json:"active"`
	Permissions []PermissionDTO `json:"permissions,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// PermissionDTO represents permission information
type PermissionDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateRoleRequestDTO represents a role creation request
type CreateRoleRequestDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Active      *bool  `json:"active"`
}

// UpdateRoleRequestDTO represents a role update request
type UpdateRoleRequestDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Active      *bool  `json:"active"`
}

// CreatePermissionRequestDTO represents a permission creation request
type CreatePermissionRequestDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Active      *bool  `json:"active"`
}

// UpdatePermissionRequestDTO represents a permission update request
type UpdatePermissionRequestDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Active      *bool  `json:"active"`
}

// AssignRoleRequestDTO represents a role assignment request
type AssignRoleRequestDTO struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

// AssignPermissionRequestDTO represents a permission assignment request
type AssignPermissionRequestDTO struct {
	RoleID       uint `json:"role_id" validate:"required"`
	PermissionID uint `json:"permission_id" validate:"required"`
}

// ErrorResponseDTO represents an error response
type ErrorResponseDTO struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// UpdateProfileRequestDTO represents a profile update request
type UpdateProfileRequestDTO struct {
	FirstName string `json:"first_name" validate:"min=2"`
	LastName  string `json:"last_name" validate:"min=2"`
}

// SuccessResponseDTO represents a success response
type SuccessResponseDTO struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
