package service

import (
	"context"
	"errors"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// AuthenticationService handles user authentication
type AuthenticationService interface {
	// Login authenticates a user and returns a token
	Login(ctx context.Context, email, password string) (*LoginResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req *RegisterRequest) (*entity.User, error)

	// RefreshToken refreshes an access token
	RefreshToken(ctx context.Context, tokenString string) (*TokenResponse, error)

	// GetUserFromToken extracts and validates user from token
	GetUserFromToken(ctx context.Context, tokenString string) (*entity.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error

	// ResetPassword resets a user's password (admin function)
	ResetPassword(ctx context.Context, userID uint, newPassword string) error
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User        *entity.User `json:"user"`
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// TokenResponse represents a token refresh response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type authenticationService struct {
	userRepo repository.UserRepository
	jwtSvc   JWTService
	authSvc  AuthorizationService
}

// NewAuthenticationService creates a new authentication service
func NewAuthenticationService(
	userRepo repository.UserRepository,
	jwtSvc JWTService,
	authSvc AuthorizationService,
) AuthenticationService {
	return &authenticationService{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
		authSvc:  authSvc,
	}
}

// Login authenticates a user and returns a token
func (a *authenticationService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Get user with roles and permissions
	user, err := a.userRepo.GetByEmailWithRoles(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		return nil, ErrUserNotActive
	}

	// Verify password
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	// Extract roles and permissions
	roles := make([]string, len(user.Roles))
	permissionMap := make(map[string]bool)

	for i, role := range user.Roles {
		roles[i] = role.Name
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}

	permissions := make([]string, 0, len(permissionMap))
	for permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	// Generate JWT token
	token, err := a.jwtSvc.GenerateToken(user.ID, user.Email, roles, permissions)
	if err != nil {
		return nil, err
	}

	// Sync user permissions with Casbin
	if err := a.authSvc.SyncUserPermissions(user); err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1 hour in seconds
	}, nil
}

// Register creates a new user account
func (a *authenticationService) Register(ctx context.Context, req *RegisterRequest) (*entity.User, error) {
	// Check if email already exists
	exists, err := a.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Create new user
	user := &entity.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Active:    true,
	}

	// Hash password
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// RefreshToken refreshes an access token
func (a *authenticationService) RefreshToken(ctx context.Context, tokenString string) (*TokenResponse, error) {
	newToken, err := a.jwtSvc.RefreshToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken: newToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1 hour in seconds
	}, nil
}

// GetUserFromToken extracts and validates user from token
func (a *authenticationService) GetUserFromToken(ctx context.Context, tokenString string) (*entity.User, error) {
	claims, err := a.jwtSvc.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetByIDWithRoles(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.Active {
		return nil, ErrUserNotActive
	}

	return user, nil
}

// ChangePassword changes a user's password
func (a *authenticationService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return ErrInvalidCredentials
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// Update user
	return a.userRepo.Update(ctx, user)
}

// ResetPassword resets a user's password (admin function)
func (a *authenticationService) ResetPassword(ctx context.Context, userID uint, newPassword string) error {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// Update user
	return a.userRepo.Update(ctx, user)
}
