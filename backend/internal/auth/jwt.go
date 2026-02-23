package auth

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT payload for CargoMax tokens.
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Email    string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	TokenType string    `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

// CreateAccessToken creates a short-lived (15 min) EdDSA-signed JWT containing
// full user identity claims.
func CreateAccessToken(privKey ed25519.PrivateKey, userID, tenantID uuid.UUID, email, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signed, err := token.SignedString(privKey)
	if err != nil {
		return "", fmt.Errorf("auth: sign access token: %w", err)
	}
	return signed, nil
}

// CreateRefreshToken creates a long-lived (7 day) EdDSA-signed JWT containing
// only the user and tenant identifiers.
func CreateRefreshToken(privKey ed25519.PrivateKey, userID, tenantID uuid.UUID) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TenantID:  tenantID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signed, err := token.SignedString(privKey)
	if err != nil {
		return "", fmt.Errorf("auth: sign refresh token: %w", err)
	}
	return signed, nil
}

// ValidateToken parses and validates an EdDSA-signed JWT, returning the
// embedded claims on success.
func ValidateToken(pubKey ed25519.PublicKey, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("auth: unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth: validate token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("auth: invalid token claims")
	}
	return claims, nil
}
