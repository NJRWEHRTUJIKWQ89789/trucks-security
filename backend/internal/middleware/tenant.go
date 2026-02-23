package middleware

import (
	"net/http"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
)

// TenantMiddleware verifies that a valid tenant_id exists in the request
// context. It should be placed after AuthMiddleware in the middleware chain.
//
// Requests missing a tenant ID receive a 403 Forbidden response.
func TenantMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID, ok := r.Context().Value(models.CtxTenantID).(uuid.UUID)
			if !ok || tenantID == uuid.Nil {
				http.Error(w, `{"error":"missing tenant context"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
