package service

import (
	"errors"
	"time"

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
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService interface {
	// GenerateToken generates a JWT token for a user
	GenerateToken(userID uint, email string, roles []string, permissions []string) (string, error)

	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(tokenString string) (*TokenClaims, error)

	// RefreshToken generates a new token with updated expiration
	RefreshToken(tokenString string) (string, error)

	// ExtractClaims extracts claims from a token without validating expiration
	ExtractClaims(tokenString string) (*TokenClaims, error)
}

type jwtService struct {
	secretKey       string
	issuer          string
	expirationTime  time.Duration
	refreshDuration time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey, issuer string, expirationTime, refreshDuration time.Duration) JWTService {
	return &jwtService{
		secretKey:       secretKey,
		issuer:          issuer,
		expirationTime:  expirationTime,
		refreshDuration: refreshDuration,
	}
}

// GenerateToken generates a JWT token for a user
func (j *jwtService) GenerateToken(userID uint, email string, roles []string, permissions []string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(j.expirationTime)

	claims := &TokenClaims{
		UserID:      userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (j *jwtService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenClaims
}

// RefreshToken generates a new token with updated expiration
func (j *jwtService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ExtractClaims(tokenString)
	if err != nil {
		return "", err
	}

	// Check if the token is within refresh window
	if time.Until(claims.ExpiresAt.Time) > j.refreshDuration {
		return "", errors.New("token is still valid and doesn't need refresh")
	}

	// Generate new token with same claims but updated expiration
	return j.GenerateToken(claims.UserID, claims.Email, claims.Roles, claims.Permissions)
}

// ExtractClaims extracts claims from a token without validating expiration
func (j *jwtService) ExtractClaims(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims, nil
	}

	return nil, ErrTokenClaims
}
