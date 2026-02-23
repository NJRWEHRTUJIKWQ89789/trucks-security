package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo handles database operations for users.
type UserRepo struct {
	db *pgxpool.Pool
}

// NewUserRepo creates a new UserRepo instance.
func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user with a bcrypt-hashed password.
func (r *UserRepo) Create(ctx context.Context, u *models.User, plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hash)
	u.ID = uuid.New()

	_, err = r.db.Exec(ctx,
		`INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())`,
		u.ID, u.TenantID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role, u.EmailVerified, u.EmailVerifyToken, u.AvatarURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID within a tenant.
func (r *UserRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at
		 FROM users WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.EmailVerified, &u.EmailVerifyToken, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return u, nil
}

// GetByEmail retrieves a user by email within a tenant.
func (r *UserRepo) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at
		 FROM users WHERE email = $1 AND tenant_id = $2`,
		email, tenantID,
	).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.EmailVerified, &u.EmailVerifyToken, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return u, nil
}

// UpdateEmailVerified sets the email_verified flag for a user.
func (r *UserRepo) UpdateEmailVerified(ctx context.Context, tenantID, id uuid.UUID, verified bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email_verified = $1, email_verify_token = NULL, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`,
		verified, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update email verified: %w", err)
	}
	return nil
}

// SetVerifyToken sets the email verification token for a user.
func (r *UserRepo) SetVerifyToken(ctx context.Context, tenantID, id uuid.UUID, token string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email_verify_token = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`,
		token, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to set verify token: %w", err)
	}
	return nil
}

// GetByVerifyToken retrieves a user by their email verification token within a tenant.
func (r *UserRepo) GetByVerifyToken(ctx context.Context, tenantID uuid.UUID, token string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at
		 FROM users WHERE email_verify_token = $1 AND tenant_id = $2`,
		token, tenantID,
	).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.EmailVerified, &u.EmailVerifyToken, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by verify token: %w", err)
	}
	return u, nil
}

// List returns a paginated list of users within a tenant.
func (r *UserRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.User, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at
		 FROM users WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.EmailVerified, &u.EmailVerifyToken, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, total, nil
}

// GetByEmailGlobal retrieves a user by email across all tenants (used for login).
func (r *UserRepo) GetByEmailGlobal(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.EmailVerified, &u.EmailVerifyToken, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return u, nil
}

// GetByVerifyTokenAnyTenant retrieves a user by verification token across all tenants.
func (r *UserRepo) GetByVerifyTokenAnyTenant(ctx context.Context, token string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, email_verify_token, avatar_url, created_at, updated_at
		 FROM users WHERE email_verify_token = $1`,
		token,
	).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.EmailVerified, &u.EmailVerifyToken, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by verify token: %w", err)
	}
	return u, nil
}
