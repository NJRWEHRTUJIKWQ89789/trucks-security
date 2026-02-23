package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VendorRepo handles database operations for vendors.
type VendorRepo struct {
	db *pgxpool.Pool
}

// NewVendorRepo creates a new VendorRepo instance.
func NewVendorRepo(db *pgxpool.Pool) *VendorRepo {
	return &VendorRepo{db: db}
}

// Create inserts a new vendor.
func (r *VendorRepo) Create(ctx context.Context, v *models.Vendor) error {
	v.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO vendors (id, tenant_id, name, contact_person, email, phone, address, category, rating, contract_start, contract_end, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())`,
		v.ID, v.TenantID, v.Name, v.ContactPerson, v.Email, v.Phone, v.Address, v.Category, v.Rating, v.ContractStart, v.ContractEnd, v.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create vendor: %w", err)
	}
	return nil
}

// GetByID retrieves a vendor by ID within a tenant.
func (r *VendorRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Vendor, error) {
	v := &models.Vendor{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, name, contact_person, email, phone, address, category, rating, contract_start, contract_end, status, created_at, updated_at
		 FROM vendors WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&v.ID, &v.TenantID, &v.Name, &v.ContactPerson, &v.Email, &v.Phone, &v.Address, &v.Category, &v.Rating, &v.ContractStart, &v.ContractEnd, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor by id: %w", err)
	}
	return v, nil
}

// List returns a paginated list of vendors within a tenant.
func (r *VendorRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Vendor, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM vendors WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendors: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, name, contact_person, email, phone, address, category, rating, contract_start, contract_end, status, created_at, updated_at
		 FROM vendors WHERE tenant_id = $1 ORDER BY name ASC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vendors: %w", err)
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var v models.Vendor
		if err := rows.Scan(&v.ID, &v.TenantID, &v.Name, &v.ContactPerson, &v.Email, &v.Phone, &v.Address, &v.Category, &v.Rating, &v.ContractStart, &v.ContractEnd, &v.Status, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor: %w", err)
		}
		vendors = append(vendors, v)
	}
	return vendors, total, nil
}

// Update modifies an existing vendor.
func (r *VendorRepo) Update(ctx context.Context, tenantID, id uuid.UUID, v *models.Vendor) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE vendors SET name = $1, contact_person = $2, email = $3, phone = $4, address = $5, category = $6, rating = $7, contract_start = $8, contract_end = $9, status = $10, updated_at = NOW()
		 WHERE id = $11 AND tenant_id = $12`,
		v.Name, v.ContactPerson, v.Email, v.Phone, v.Address, v.Category, v.Rating, v.ContractStart, v.ContractEnd, v.Status, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update vendor: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("vendor not found")
	}
	return nil
}

// Delete removes a vendor by ID within a tenant.
func (r *VendorRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM vendors WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete vendor: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("vendor not found")
	}
	return nil
}
