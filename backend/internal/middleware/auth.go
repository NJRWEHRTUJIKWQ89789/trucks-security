package middleware

import (
	"context"
	"crypto/ed25519"
	"net/http"

	"cargomax-api/internal/auth"
	"cargomax-api/internal/models"
)

// AuthMiddleware returns an HTTP middleware that validates the EdDSA JWT from
// the "cargomax_token" cookie and injects the authenticated user's claims
// (tenant ID, user ID, role, email) into the request context.
//
// Requests without a valid token receive a 401 Unauthorized response.
func AuthMiddleware(pubKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := auth.GetAccessToken(r)
			if tokenString == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateToken(pubKey, tokenString)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, models.CtxTenantID, claims.TenantID)
			ctx = context.WithValue(ctx, models.CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, models.CtxUserRole, claims.Role)
			ctx = context.WithValue(ctx, models.CtxUserEmail, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
