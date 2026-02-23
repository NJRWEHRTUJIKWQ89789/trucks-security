package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantRepo handles database operations for tenants.
type TenantRepo struct {
	db *pgxpool.Pool
}

// NewTenantRepo creates a new TenantRepo instance.
func NewTenantRepo(db *pgxpool.Pool) *TenantRepo {
	return &TenantRepo{db: db}
}

// Create inserts a new tenant.
func (r *TenantRepo) Create(ctx context.Context, t *models.Tenant) error {
	t.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO tenants (id, name, domain, plan, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		t.ID, t.Name, t.Domain, t.Plan,
	)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}

// GetByID retrieves a tenant by its ID.
func (r *TenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	t := &models.Tenant{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, domain, plan, created_at, updated_at
		 FROM tenants WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Name, &t.Domain, &t.Plan, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by id: %w", err)
	}
	return t, nil
}

// GetByDomain retrieves a tenant by its domain.
func (r *TenantRepo) GetByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	t := &models.Tenant{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, domain, plan, created_at, updated_at
		 FROM tenants WHERE domain = $1`,
		domain,
	).Scan(&t.ID, &t.Name, &t.Domain, &t.Plan, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by domain: %w", err)
	}
	return t, nil
}
