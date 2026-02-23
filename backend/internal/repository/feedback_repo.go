package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FeedbackRepo handles database operations for client feedback.
type FeedbackRepo struct {
	db *pgxpool.Pool
}

// NewFeedbackRepo creates a new FeedbackRepo instance.
func NewFeedbackRepo(db *pgxpool.Pool) *FeedbackRepo {
	return &FeedbackRepo{db: db}
}

// Create inserts a new feedback entry.
func (r *FeedbackRepo) Create(ctx context.Context, f *models.Feedback) error {
	f.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO client_feedback (id, tenant_id, client_id, rating, comment, category, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		f.ID, f.TenantID, f.ClientID, f.Rating, f.Comment, f.Category,
	)
	if err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}
	return nil
}

// ListByClient returns a paginated list of feedback for a specific client within a tenant.
func (r *FeedbackRepo) ListByClient(ctx context.Context, tenantID, clientID uuid.UUID, page, perPage int) ([]models.Feedback, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM client_feedback WHERE tenant_id = $1 AND client_id = $2`,
		tenantID, clientID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count feedback: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, client_id, rating, comment, category, created_at
		 FROM client_feedback WHERE tenant_id = $1 AND client_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		tenantID, clientID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list feedback: %w", err)
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var f models.Feedback
		if err := rows.Scan(&f.ID, &f.TenantID, &f.ClientID, &f.Rating, &f.Comment, &f.Category, &f.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan feedback: %w", err)
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, total, nil
}

// ListAll returns a paginated list of all feedback for a tenant, including the client company name.
func (r *FeedbackRepo) ListAll(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Feedback, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM client_feedback WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count feedback: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT f.id, f.tenant_id, f.client_id, f.rating, f.comment, f.category, f.created_at, c.company_name
		 FROM client_feedback f
		 JOIN clients c ON c.id = f.client_id AND c.tenant_id = f.tenant_id
		 WHERE f.tenant_id = $1
		 ORDER BY f.created_at DESC
		 LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list all feedback: %w", err)
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var f models.Feedback
		if err := rows.Scan(&f.ID, &f.TenantID, &f.ClientID, &f.Rating, &f.Comment, &f.Category, &f.CreatedAt, &f.ClientName); err != nil {
			return nil, 0, fmt.Errorf("failed to scan feedback: %w", err)
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, total, nil
}

// GetAvgRating returns the average rating for a specific client within a tenant.
func (r *FeedbackRepo) GetAvgRating(ctx context.Context, tenantID, clientID uuid.UUID) (float64, error) {
	var avg *float64
	err := r.db.QueryRow(ctx,
		`SELECT AVG(rating) FROM client_feedback WHERE tenant_id = $1 AND client_id = $2`,
		tenantID, clientID,
	).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}
	if avg == nil {
		return 0, nil
	}
	return *avg, nil
}
