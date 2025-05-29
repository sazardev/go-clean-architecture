package auth

import (
	"context"
	"errors"
	"time"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"
	"go-clean-architecture/internal/infrastructure/auth/jwt"
	"go-clean-architecture/internal/infrastructure/auth/rbac"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// AuthService provides authentication functionality
type AuthService struct {
	userRepo      repository.UserRepository
	roleRepo      repository.RoleRepository
	tokenService  *jwt.TokenService
	policyManager *rbac.PolicyManager
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tokenService *jwt.TokenService,
	policyManager *rbac.PolicyManager,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		tokenService:  tokenService,
		policyManager: policyManager,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int64     `json:"expires_in"`
	User        *UserInfo `json:"user"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID          uint     `json:"id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Active      bool     `json:"active"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by email with roles
	user, err := s.userRepo.GetByEmailWithRoles(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		return nil, ErrUserInactive
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Sync user policies with Casbin
	if err := s.policyManager.SyncUserPolicies(user); err != nil {
		// Log error but don't fail login
		// logger.Error("Failed to sync user policies", "error", err)
	}

	// Prepare response
	userInfo := s.buildUserInfo(user)

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * time.Hour / time.Second), // 24 hours in seconds
		User:        userInfo,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) { // Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create new user
	user := &entity.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Active:    true,
	}

	// Set password
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}
	// Assign default role (employee)
	defaultRole, err := s.roleRepo.GetByName(ctx, "employee")
	if err != nil {
		return nil, err
	}

	user.Roles = []entity.Role{*defaultRole}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	// Reload user with roles
	user, err = s.userRepo.GetByIDWithRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Sync user policies with Casbin
	if err := s.policyManager.SyncUserPolicies(user); err != nil {
		// Log error but don't fail registration
		// logger.Error("Failed to sync user policies", "error", err)
	}

	// Prepare response
	userInfo := s.buildUserInfo(user)

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * time.Hour / time.Second),
		User:        userInfo,
	}, nil
}

// RefreshToken generates a new token from a valid refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate the refresh token
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil && err != jwt.ErrExpiredToken {
		return nil, errors.New("invalid refresh token")
	}
	// Get fresh user data
	user, err := s.userRepo.GetByIDWithRoles(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is still active
	if !user.Active {
		return nil, ErrUserInactive
	}

	// Generate new token
	newToken, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Prepare response
	userInfo := s.buildUserInfo(user)

	return &LoginResponse{
		AccessToken: newToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * time.Hour / time.Second),
		User:        userInfo,
	}, nil
}

// GetProfile returns the current user's profile
func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*UserInfo, error) {
	user, err := s.userRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.buildUserInfo(user), nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return errors.New("invalid current password")
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// Update user
	return s.userRepo.Update(ctx, user)
}

// buildUserInfo creates a UserInfo from an entity.User
func (s *AuthService) buildUserInfo(user *entity.User) *UserInfo {
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	permissions := make([]string, 0)
	permissionMap := make(map[string]bool)

	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			if !permissionMap[permission.Name] {
				permissions = append(permissions, permission.Name)
				permissionMap[permission.Name] = true
			}
		}
	}

	return &UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Active:      user.Active,
		Roles:       roles,
		Permissions: permissions,
	}
}
