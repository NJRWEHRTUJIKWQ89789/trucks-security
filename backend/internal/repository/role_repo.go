package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoleRepo handles database operations for roles.
type RoleRepo struct {
	db *pgxpool.Pool
}

// NewRoleRepo creates a new RoleRepo instance.
func NewRoleRepo(db *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{db: db}
}

// Create inserts a new role.
func (r *RoleRepo) Create(ctx context.Context, role *models.Role) error {
	role.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO roles (id, tenant_id, name, permissions, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		role.ID, role.TenantID, role.Name, role.Permissions,
	)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// GetByID retrieves a role by ID within a tenant.
func (r *RoleRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Role, error) {
	role := &models.Role{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, name, permissions, created_at, updated_at
		 FROM roles WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&role.ID, &role.TenantID, &role.Name, &role.Permissions, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by id: %w", err)
	}
	return role, nil
}

// List returns a paginated list of roles within a tenant.
func (r *RoleRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Role, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM roles WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, name, permissions, created_at, updated_at
		 FROM roles WHERE tenant_id = $1 ORDER BY name ASC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Permissions, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, total, nil
}

// Update modifies an existing role.
func (r *RoleRepo) Update(ctx context.Context, tenantID, id uuid.UUID, role *models.Role) error {
	_, err := r.db.Exec(ctx,
		`UPDATE roles SET name = $1, permissions = $2, updated_at = NOW() WHERE id = $3 AND tenant_id = $4`,
		role.Name, role.Permissions, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

// Delete removes a role by ID within a tenant.
func (r *RoleRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM roles WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
