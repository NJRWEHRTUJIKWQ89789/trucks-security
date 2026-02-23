package resolvers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"cargomax-api/internal/auth"
	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
)

// AuthQueries returns the GraphQL query fields related to authentication.
func (r *Resolver) AuthQueries() graphql.Fields {
	return graphql.Fields{
		"me": &graphql.Field{
			Type:        types.UserType,
			Description: "Returns the currently authenticated user.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID, ok := p.Context.Value(models.CtxUserID).(uuid.UUID)
				if !ok {
					return nil, fmt.Errorf("authentication required")
				}
				tenantID := p.Context.Value(models.CtxTenantID).(uuid.UUID)

				user, err := r.UserRepo.GetByID(p.Context, tenantID, userID)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch user: %w", err)
				}
				return user, nil
			},
		},
	}
}

// AuthMutations returns the GraphQL mutation fields related to authentication.
func (r *Resolver) AuthMutations() graphql.Fields {
	return graphql.Fields{
		// -----------------------------------------------------------------
		// register
		// -----------------------------------------------------------------
		"register": &graphql.Field{
			Type:        types.AuthPayloadType,
			Description: "Create a new tenant and user account, returning the user and setting auth cookies.",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.RegisterInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})

				email := input["email"].(string)
				password := input["password"].(string)
				firstName := input["firstName"].(string)
				lastName := input["lastName"].(string)
				tenantName := input["tenantName"].(string)

				// Create tenant.
				tenant := &models.Tenant{
					ID:        uuid.New(),
					Name:      tenantName,
					Plan:      "starter",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := r.TenantRepo.Create(p.Context, tenant); err != nil {
					return nil, fmt.Errorf("failed to create tenant: %w", err)
				}

				// Generate email verification token.
				verifyBytes := make([]byte, 32)
				if _, err := rand.Read(verifyBytes); err != nil {
					return nil, fmt.Errorf("failed to generate verify token: %w", err)
				}
				verifyToken := hex.EncodeToString(verifyBytes)

				// Create user (UserRepo.Create hashes the plain password internally).
				user := &models.User{
					TenantID:         tenant.ID,
					Email:            email,
					FirstName:        &firstName,
					LastName:         &lastName,
					Role:             "admin",
					EmailVerified:    false,
					EmailVerifyToken: &verifyToken,
				}
				if err := r.UserRepo.Create(p.Context, user, password); err != nil {
					return nil, fmt.Errorf("failed to create user: %w", err)
				}

				// Issue tokens.
				accessToken, err := auth.CreateAccessToken(r.Config.JWTPrivateKey, user.ID, tenant.ID, user.Email, user.Role)
				if err != nil {
					return nil, fmt.Errorf("failed to create access token: %w", err)
				}
				refreshToken, err := auth.CreateRefreshToken(r.Config.JWTPrivateKey, user.ID, tenant.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to create refresh token: %w", err)
				}

				// Set cookies via the response writer stored in context.
				w := p.Context.Value(models.CtxResponseWriter).(http.ResponseWriter)
				auth.SetAuthCookies(w, accessToken, refreshToken, r.Config.CookieDomain, r.Config.CookieSecure)

				return map[string]interface{}{
					"user":  user,
					"token": accessToken,
				}, nil
			},
		},

		// -----------------------------------------------------------------
		// login
		// -----------------------------------------------------------------
		"login": &graphql.Field{
			Type:        types.AuthPayloadType,
			Description: "Authenticate with email and password, returning the user and setting auth cookies.",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.LoginInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})
				email := input["email"].(string)
				password := input["password"].(string)

				// Look up user by email across all tenants (login is pre-auth,
				// so there is no tenant in the request context yet).
				user, err := r.UserRepo.GetByEmailGlobal(p.Context, email)
				if err != nil {
					return nil, fmt.Errorf("invalid email or password")
				}

				// Verify password.
				if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
					return nil, fmt.Errorf("invalid email or password")
				}

				// Issue tokens.
				accessToken, err := auth.CreateAccessToken(r.Config.JWTPrivateKey, user.ID, user.TenantID, user.Email, user.Role)
				if err != nil {
					return nil, fmt.Errorf("failed to create access token: %w", err)
				}
				refreshToken, err := auth.CreateRefreshToken(r.Config.JWTPrivateKey, user.ID, user.TenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to create refresh token: %w", err)
				}

				// Set cookies.
				w := p.Context.Value(models.CtxResponseWriter).(http.ResponseWriter)
				auth.SetAuthCookies(w, accessToken, refreshToken, r.Config.CookieDomain, r.Config.CookieSecure)

				return map[string]interface{}{
					"user":  user,
					"token": accessToken,
				}, nil
			},
		},

		// -----------------------------------------------------------------
		// logout
		// -----------------------------------------------------------------
		"logout": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Clear auth cookies and log out the current user.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				w := p.Context.Value(models.CtxResponseWriter).(http.ResponseWriter)
				auth.ClearAuthCookies(w, r.Config.CookieDomain, r.Config.CookieSecure)
				return true, nil
			},
		},

		// -----------------------------------------------------------------
		// verifyEmail
		// -----------------------------------------------------------------
		"verifyEmail": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Verify a user's email address using the token sent during registration.",
			Args: graphql.FieldConfigArgument{
				"token": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				token := p.Args["token"].(string)
				if token == "" {
					return false, fmt.Errorf("verification token is required")
				}

				// The verification token is globally unique, so we can look it
				// up without a tenant filter (user is not authenticated yet).
				user, err := r.UserRepo.GetByVerifyTokenAnyTenant(p.Context, token)
				if err != nil {
					return false, fmt.Errorf("invalid or expired verification token")
				}

				if err := r.UserRepo.UpdateEmailVerified(p.Context, user.TenantID, user.ID, true); err != nil {
					return false, fmt.Errorf("failed to verify email: %w", err)
				}

				return true, nil
			},
		},

		// -----------------------------------------------------------------
		// refreshToken
		// -----------------------------------------------------------------
		"refreshToken": &graphql.Field{
			Type:        types.AuthPayloadType,
			Description: "Exchange a valid refresh token cookie for a new access/refresh token pair.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Read the HTTP request from context so we can access the
				// refresh cookie directly, independent of the optional auth
				// middleware (which only validates the access token).
				req, ok := p.Context.Value(models.CtxHTTPRequest).(*http.Request)
				if !ok {
					return nil, fmt.Errorf("internal error: HTTP request not found in context")
				}

				// Extract the refresh token from the cookie.
				refreshTokenStr := auth.GetRefreshToken(req)
				if refreshTokenStr == "" {
					return nil, fmt.Errorf("refresh token cookie is missing")
				}

				// Validate the refresh token and extract claims.
				claims, err := auth.ValidateToken(r.Config.JWTPublicKey, refreshTokenStr)
				if err != nil {
					return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
				}

				// Ensure this is actually a refresh token, not an access token.
				if claims.TokenType != "refresh" {
					return nil, fmt.Errorf("invalid token type: expected refresh token")
				}

				// Fetch fresh user to get current role/email.
				user, err := r.UserRepo.GetByID(p.Context, claims.TenantID, claims.UserID)
				if err != nil {
					return nil, fmt.Errorf("user not found")
				}

				// Issue new token pair.
				accessToken, err := auth.CreateAccessToken(r.Config.JWTPrivateKey, user.ID, user.TenantID, user.Email, user.Role)
				if err != nil {
					return nil, fmt.Errorf("failed to create access token: %w", err)
				}
				refreshToken, err := auth.CreateRefreshToken(r.Config.JWTPrivateKey, user.ID, user.TenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to create refresh token: %w", err)
				}

				// Set new cookies.
				w := p.Context.Value(models.CtxResponseWriter).(http.ResponseWriter)
				auth.SetAuthCookies(w, accessToken, refreshToken, r.Config.CookieDomain, r.Config.CookieSecure)

				return map[string]interface{}{
					"user":  user,
					"token": accessToken,
				}, nil
			},
		},
	}
}
