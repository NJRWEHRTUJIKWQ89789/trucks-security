package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClientRepo handles database operations for clients.
type ClientRepo struct {
	db *pgxpool.Pool
}

// NewClientRepo creates a new ClientRepo instance.
func NewClientRepo(db *pgxpool.Pool) *ClientRepo {
	return &ClientRepo{db: db}
}

// Create inserts a new client.
func (r *ClientRepo) Create(ctx context.Context, c *models.Client) error {
	c.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO clients (id, tenant_id, company_name, contact_person, email, phone, address, industry, total_shipments, total_spent, satisfaction_rating, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())`,
		c.ID, c.TenantID, c.CompanyName, c.ContactPerson, c.Email, c.Phone, c.Address, c.Industry, c.TotalShipments, c.TotalSpent, c.SatisfactionRating, c.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	return nil
}

// GetByID retrieves a client by ID within a tenant.
func (r *ClientRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Client, error) {
	c := &models.Client{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, company_name, contact_person, email, phone, address, industry, total_shipments, total_spent, satisfaction_rating, status, created_at, updated_at
		 FROM clients WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&c.ID, &c.TenantID, &c.CompanyName, &c.ContactPerson, &c.Email, &c.Phone, &c.Address, &c.Industry, &c.TotalShipments, &c.TotalSpent, &c.SatisfactionRating, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get client by id: %w", err)
	}
	return c, nil
}

// List returns a paginated list of clients within a tenant.
func (r *ClientRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Client, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM clients WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count clients: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, company_name, contact_person, email, phone, address, industry, total_shipments, total_spent, satisfaction_rating, status, created_at, updated_at
		 FROM clients WHERE tenant_id = $1 ORDER BY company_name ASC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list clients: %w", err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.TenantID, &c.CompanyName, &c.ContactPerson, &c.Email, &c.Phone, &c.Address, &c.Industry, &c.TotalShipments, &c.TotalSpent, &c.SatisfactionRating, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan client: %w", err)
		}
		clients = append(clients, c)
	}
	return clients, total, nil
}

// Update modifies an existing client.
func (r *ClientRepo) Update(ctx context.Context, tenantID, id uuid.UUID, c *models.Client) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE clients SET company_name = $1, contact_person = $2, email = $3, phone = $4, address = $5, industry = $6, total_shipments = $7, total_spent = $8, satisfaction_rating = $9, status = $10, updated_at = NOW()
		 WHERE id = $11 AND tenant_id = $12`,
		c.CompanyName, c.ContactPerson, c.Email, c.Phone, c.Address, c.Industry, c.TotalShipments, c.TotalSpent, c.SatisfactionRating, c.Status, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("client not found")
	}
	return nil
}

// Delete removes a client by ID within a tenant.
func (r *ClientRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM clients WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("client not found")
	}
	return nil
}

// GetTopBySpent returns the top N clients by total_spent within a tenant.
func (r *ClientRepo) GetTopBySpent(ctx context.Context, tenantID uuid.UUID, limit int) ([]models.Client, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, company_name, contact_person, email, phone, address, industry, total_shipments, total_spent, satisfaction_rating, status, created_at, updated_at
		 FROM clients WHERE tenant_id = $1 AND status = 'active' ORDER BY total_spent DESC NULLS LAST LIMIT $2`,
		tenantID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get top clients by spent: %w", err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.TenantID, &c.CompanyName, &c.ContactPerson, &c.Email, &c.Phone, &c.Address, &c.Industry, &c.TotalShipments, &c.TotalSpent, &c.SatisfactionRating, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan top client: %w", err)
		}
		clients = append(clients, c)
	}
	return clients, nil
}
