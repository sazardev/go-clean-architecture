package usecase

import (
	"context"
	"errors"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"
	"go-clean-architecture/internal/infrastructure/auth"
	"go-clean-architecture/internal/infrastructure/auth/rbac"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	authService    *auth.AuthService
	policyManager  *rbac.PolicyManager
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	authService *auth.AuthService,
	policyManager *rbac.PolicyManager,
) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		authService:    authService,
		policyManager:  policyManager,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*entity.User, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create user
	user := &entity.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Active:    true,
	}

	// Set password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Assign default role
	defaultRole, err := uc.roleRepo.GetByName(ctx, "employee")
	if err != nil {
		return nil, err
	}

	user.Roles = []entity.Role{*defaultRole}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Sync policies
	if err := uc.policyManager.SyncUserPolicies(user); err != nil {
		// Log error but don't fail
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepo.GetByIDWithRoles(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return uc.userRepo.GetByEmailWithRoles(ctx, email)
}

// GetAllUsers retrieves all users
func (uc *UserUseCase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return uc.userRepo.List(ctx, 0, 1000) // Get first 1000 users
}

// UpdateUser updates a user
func (uc *UserUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	// Update user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Sync policies
	if err := uc.policyManager.SyncUserPolicies(user); err != nil {
		// Log error but don't fail
	}

	return nil
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) error {
	// Get user first
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Remove from RBAC
	if err := uc.policyManager.RemoveRoleFromUser(user.Email, ""); err != nil {
		// Log error but continue
	}

	// Delete user
	return uc.userRepo.Delete(ctx, id)
}

// AssignRoleToUser assigns a role to a user
func (uc *UserUseCase) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	// Get user and role
	user, err := uc.userRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return err
	}

	role, err := uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Check if user already has the role
	for _, userRole := range user.Roles {
		if userRole.ID == roleID {
			return errors.New("user already has this role")
		}
	}

	// Assign role in database
	if err := uc.userRepo.AssignRole(ctx, userID, roleID); err != nil {
		return err
	}

	// Assign role in RBAC
	if err := uc.policyManager.AssignRoleToUser(user.Email, role.Name); err != nil {
		return err
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (uc *UserUseCase) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	// Get user and role
	user, err := uc.userRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return err
	}

	role, err := uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Remove role from database
	if err := uc.userRepo.RemoveRole(ctx, userID, roleID); err != nil {
		return err
	}

	// Remove role from RBAC
	if err := uc.policyManager.RemoveRoleFromUser(user.Email, role.Name); err != nil {
		return err
	}

	return nil
}

// ActivateUser activates a user account
func (uc *UserUseCase) ActivateUser(ctx context.Context, id uint) error {
	return uc.userRepo.ActivateUser(ctx, id)
}

// DeactivateUser deactivates a user account
func (uc *UserUseCase) DeactivateUser(ctx context.Context, id uint) error {
	return uc.userRepo.DeactivateUser(ctx, id)
}

// CheckUserPermission checks if a user has a specific permission
func (uc *UserUseCase) CheckUserPermission(ctx context.Context, userEmail, resource, action string) (bool, error) {
	return uc.policyManager.CheckPermission(userEmail, resource, action)
}
