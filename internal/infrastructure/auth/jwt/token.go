package jwt

import (
	"errors"
	"time"

	"go-clean-architecture/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrTokenClaims  = errors.New("invalid token claims")
)

// TokenClaims represents the claims stored in JWT tokens
type TokenClaims struct {
	UserID      uint     `json:"user_id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// TokenService handles JWT token operations
type TokenService struct {
	secretKey       []byte
	tokenExpiration time.Duration
	issuer          string
}

// NewTokenService creates a new JWT token service
func NewTokenService(secretKey string, tokenExpiration time.Duration, issuer string) *TokenService {
	return &TokenService{
		secretKey:       []byte(secretKey),
		tokenExpiration: tokenExpiration,
		issuer:          issuer,
	}
}

// GenerateToken generates a JWT token for a user
func (t *TokenService) GenerateToken(user *entity.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	// Extract role names and permissions
	roles := make([]string, len(user.Roles))
	permissionMap := make(map[string]bool)
	var permissions []string

	for i, role := range user.Roles {
		roles[i] = role.Name

		// Collect unique permissions
		for _, permission := range role.Permissions {
			if !permissionMap[permission.Name] {
				permissions = append(permissions, permission.Name)
				permissionMap[permission.Name] = true
			}
		}
	}

	// Create claims
	claims := &TokenClaims{
		UserID:      user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create and sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (t *TokenService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return t.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenClaims
	}

	return claims, nil
}

// RefreshToken generates a new token from valid existing claims
func (t *TokenService) RefreshToken(claims *TokenClaims) (string, error) {
	if claims == nil {
		return "", errors.New("claims cannot be nil")
	}

	// Create new claims with extended expiration
	newClaims := &TokenClaims{
		UserID:      claims.UserID,
		Email:       claims.Email,
		FirstName:   claims.FirstName,
		LastName:    claims.LastName,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   claims.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create and sign new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString(t.secretKey)
}

// ExtractTokenFromBearer extracts JWT token from Bearer authorization header
func ExtractTokenFromBearer(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
