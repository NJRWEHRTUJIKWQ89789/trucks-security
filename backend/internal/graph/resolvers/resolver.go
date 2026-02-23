package resolvers

import (
	"context"
	"fmt"

	"cargomax-api/internal/config"
	"cargomax-api/internal/models"
	"cargomax-api/internal/repository"

	"github.com/google/uuid"
)

// Resolver holds references to every repository and the application configuration.
// It is the single root object shared across all GraphQL resolver methods.
type Resolver struct {
	UserRepo         *repository.UserRepo
	TenantRepo       *repository.TenantRepo
	ShipmentRepo     *repository.ShipmentRepo
	VehicleRepo      *repository.VehicleRepo
	DriverRepo       *repository.DriverRepo
	MaintenanceRepo  *repository.MaintenanceRepo
	WarehouseRepo    *repository.WarehouseRepo
	InventoryRepo    *repository.InventoryRepo
	OrderRepo        *repository.OrderRepo
	VendorRepo       *repository.VendorRepo
	ClientRepo       *repository.ClientRepo
	FeedbackRepo     *repository.FeedbackRepo
	DashboardRepo    *repository.DashboardRepo
	ReportRepo       *repository.ReportRepo
	NotificationRepo *repository.NotificationRepo
	SettingRepo      *repository.SettingRepo
	RoleRepo         *repository.RoleRepo
	ActivityRepo     *repository.ActivityRepo
	Config           *config.Config
}

// NewResolver constructs a Resolver with all required dependencies.
func NewResolver(
	userRepo *repository.UserRepo,
	tenantRepo *repository.TenantRepo,
	shipmentRepo *repository.ShipmentRepo,
	vehicleRepo *repository.VehicleRepo,
	driverRepo *repository.DriverRepo,
	maintenanceRepo *repository.MaintenanceRepo,
	warehouseRepo *repository.WarehouseRepo,
	inventoryRepo *repository.InventoryRepo,
	orderRepo *repository.OrderRepo,
	vendorRepo *repository.VendorRepo,
	clientRepo *repository.ClientRepo,
	feedbackRepo *repository.FeedbackRepo,
	dashboardRepo *repository.DashboardRepo,
	reportRepo *repository.ReportRepo,
	notificationRepo *repository.NotificationRepo,
	settingRepo *repository.SettingRepo,
	roleRepo *repository.RoleRepo,
	activityRepo *repository.ActivityRepo,
	cfg *config.Config,
) *Resolver {
	return &Resolver{
		UserRepo:         userRepo,
		TenantRepo:       tenantRepo,
		ShipmentRepo:     shipmentRepo,
		VehicleRepo:      vehicleRepo,
		DriverRepo:       driverRepo,
		MaintenanceRepo:  maintenanceRepo,
		WarehouseRepo:    warehouseRepo,
		InventoryRepo:    inventoryRepo,
		OrderRepo:        orderRepo,
		VendorRepo:       vendorRepo,
		ClientRepo:       clientRepo,
		FeedbackRepo:     feedbackRepo,
		DashboardRepo:    dashboardRepo,
		ReportRepo:       reportRepo,
		NotificationRepo: notificationRepo,
		SettingRepo:      settingRepo,
		RoleRepo:         roleRepo,
		ActivityRepo:     activityRepo,
		Config:           cfg,
	}
}

// requireAuth extracts the tenant ID and user ID from the resolver context.
// It returns a clear "authentication required" error if either value is absent,
// preventing a server panic from an unchecked type assertion.
func requireAuth(ctx context.Context) (tenantID uuid.UUID, userID uuid.UUID, err error) {
	tid, ok := ctx.Value(models.CtxTenantID).(uuid.UUID)
	if !ok || tid == uuid.Nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("authentication required")
	}
	uid, ok := ctx.Value(models.CtxUserID).(uuid.UUID)
	if !ok || uid == uuid.Nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("authentication required")
	}
	return tid, uid, nil
}

// requireTenant extracts only the tenant ID from the resolver context.
// Use this for resolvers that need tenant scoping but not necessarily a user ID.
func requireTenant(ctx context.Context) (uuid.UUID, error) {
	tid, ok := ctx.Value(models.CtxTenantID).(uuid.UUID)
	if !ok || tid == uuid.Nil {
		return uuid.Nil, fmt.Errorf("authentication required")
	}
	return tid, nil
}
