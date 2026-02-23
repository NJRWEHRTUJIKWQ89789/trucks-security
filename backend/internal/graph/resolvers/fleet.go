package resolvers

import (
	"fmt"
	"time"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// FleetQueries returns the GraphQL query fields for vehicles, drivers, and maintenance.
func (r *Resolver) FleetQueries() graphql.Fields {
	return graphql.Fields{
		// -----------------------------------------------------------------
		// vehicles (paginated)
		// -----------------------------------------------------------------
		"vehicles": &graphql.Field{
			Type:        types.VehicleConnectionType,
			Description: "Returns a paginated list of fleet vehicles.",
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
				"status":  &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				var status string
				if v, ok := p.Args["status"].(string); ok {
					status = v
				}

				items, total, err := r.VehicleRepo.List(p.Context, tenantID, status, page, perPage)
				if err != nil {
					return nil, fmt.Errorf("failed to list vehicles: %w", err)
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}

				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},

		// -----------------------------------------------------------------
		// vehicle (by ID)
		// -----------------------------------------------------------------
		"vehicle": &graphql.Field{
			Type:        types.VehicleType,
			Description: "Returns a single vehicle by its UUID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid vehicle id: %w", err)
				}

				vehicle, err := r.VehicleRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("vehicle not found: %w", err)
				}
				return vehicle, nil
			},
		},

		// -----------------------------------------------------------------
		// drivers (paginated)
		// -----------------------------------------------------------------
		"drivers": &graphql.Field{
			Type:        types.DriverConnectionType,
			Description: "Returns a paginated list of drivers.",
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
				"status":  &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				items, total, err := r.DriverRepo.List(p.Context, tenantID, page, perPage)
				if err != nil {
					return nil, fmt.Errorf("failed to list drivers: %w", err)
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}

				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},

		// -----------------------------------------------------------------
		// driver (by ID)
		// -----------------------------------------------------------------
		"driver": &graphql.Field{
			Type:        types.DriverType,
			Description: "Returns a single driver by its UUID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid driver id: %w", err)
				}

				driver, err := r.DriverRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("driver not found: %w", err)
				}
				return driver, nil
			},
		},

		// -----------------------------------------------------------------
		// maintenanceRecords (paginated, optional vehicleId filter)
		// -----------------------------------------------------------------
		"maintenanceRecords": &graphql.Field{
			Type:        types.MaintenanceConnectionType,
			Description: "Returns a paginated list of maintenance records, optionally filtered by vehicle.",
			Args: graphql.FieldConfigArgument{
				"page":      &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage":   &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
				"vehicleId": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				var vehicleID *uuid.UUID
				if vidStr, ok := p.Args["vehicleId"].(string); ok && vidStr != "" {
					parsed, parseErr := uuid.Parse(vidStr)
					if parseErr != nil {
						return nil, fmt.Errorf("invalid vehicle id: %w", parseErr)
					}
					vehicleID = &parsed
				}

				items, total, err := r.MaintenanceRepo.List(p.Context, tenantID, vehicleID, page, perPage)
				if err != nil {
					return nil, fmt.Errorf("failed to list maintenance records: %w", err)
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}

				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},
	}
}

// FleetMutations returns the GraphQL mutation fields for vehicles, drivers, and maintenance.
func (r *Resolver) FleetMutations() graphql.Fields {
	return graphql.Fields{
		// =================================================================
		// Vehicle mutations
		// =================================================================

		// -----------------------------------------------------------------
		// createVehicle
		// -----------------------------------------------------------------
		"createVehicle": &graphql.Field{
			Type:        types.VehicleType,
			Description: "Create a new vehicle for the current tenant.",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.VehicleInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				now := time.Now()
				vehicle := &models.Vehicle{
					ID:        uuid.New(),
					TenantID:  tenantID,
					VehicleID: input["vehicleId"].(string),
					Status:    "available",
					FuelLevel: 100,
					Mileage:   0,
					CreatedAt: now,
					UpdatedAt: now,
				}

				if v, ok := input["name"].(string); ok {
					vehicle.Name = &v
				}
				if v, ok := input["type"].(string); ok {
					vehicle.Type = &v
				}
				if v, ok := input["status"].(string); ok && v != "" {
					vehicle.Status = v
				}
				if v, ok := input["fuelLevel"].(int); ok {
					vehicle.FuelLevel = v
				}
				if v, ok := input["mileage"].(int); ok {
					vehicle.Mileage = v
				}
				if v, ok := input["lastService"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						vehicle.LastService = &t
					}
				}
				if v, ok := input["nextService"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						vehicle.NextService = &t
					}
				}
				if v, ok := input["licensePlate"].(string); ok {
					vehicle.LicensePlate = &v
				}
				if v, ok := input["year"].(int); ok {
					vehicle.Year = &v
				}

				if err := r.VehicleRepo.Create(p.Context, vehicle); err != nil {
					return nil, fmt.Errorf("failed to create vehicle: %w", err)
				}
				return vehicle, nil
			},
		},

		// -----------------------------------------------------------------
		// updateVehicle
		// -----------------------------------------------------------------
		"updateVehicle": &graphql.Field{
			Type:        types.VehicleType,
			Description: "Update an existing vehicle.",
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.VehicleInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid vehicle id: %w", err)
				}

				vehicle, err := r.VehicleRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("vehicle not found: %w", err)
				}

				input := p.Args["input"].(map[string]interface{})

				if v, ok := input["vehicleId"].(string); ok {
					vehicle.VehicleID = v
				}
				if v, ok := input["name"].(string); ok {
					vehicle.Name = &v
				}
				if v, ok := input["type"].(string); ok {
					vehicle.Type = &v
				}
				if v, ok := input["status"].(string); ok {
					vehicle.Status = v
				}
				if v, ok := input["fuelLevel"].(int); ok {
					vehicle.FuelLevel = v
				}
				if v, ok := input["mileage"].(int); ok {
					vehicle.Mileage = v
				}
				if v, ok := input["lastService"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						vehicle.LastService = &t
					}
				}
				if v, ok := input["nextService"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						vehicle.NextService = &t
					}
				}
				if v, ok := input["licensePlate"].(string); ok {
					vehicle.LicensePlate = &v
				}
				if v, ok := input["year"].(int); ok {
					vehicle.Year = &v
				}

				vehicle.UpdatedAt = time.Now()

				if err := r.VehicleRepo.Update(p.Context, tenantID, id, vehicle); err != nil {
					return nil, fmt.Errorf("failed to update vehicle: %w", err)
				}
				return vehicle, nil
			},
		},

		// -----------------------------------------------------------------
		// deleteVehicle
		// -----------------------------------------------------------------
		"deleteVehicle": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a vehicle by ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return false, fmt.Errorf("invalid vehicle id: %w", err)
				}

				if err := r.VehicleRepo.Delete(p.Context, tenantID, id); err != nil {
					return false, fmt.Errorf("failed to delete vehicle: %w", err)
				}
				return true, nil
			},
		},

		// =================================================================
		// Driver mutations
		// =================================================================

		// -----------------------------------------------------------------
		// createDriver
		// -----------------------------------------------------------------
		"createDriver": &graphql.Field{
			Type:        types.DriverType,
			Description: "Create a new driver for the current tenant.",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.DriverInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				now := time.Now()
				driver := &models.Driver{
					ID:              uuid.New(),
					TenantID:        tenantID,
					EmployeeID:      input["employeeId"].(string),
					Status:          "available",
					TotalDeliveries: 0,
					CreatedAt:       now,
					UpdatedAt:       now,
				}

				if v, ok := input["firstName"].(string); ok {
					driver.FirstName = &v
				}
				if v, ok := input["lastName"].(string); ok {
					driver.LastName = &v
				}
				if v, ok := input["email"].(string); ok {
					driver.Email = &v
				}
				if v, ok := input["phone"].(string); ok {
					driver.Phone = &v
				}
				if v, ok := input["licenseNumber"].(string); ok {
					driver.LicenseNumber = &v
				}
				if v, ok := input["licenseExpiry"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						driver.LicenseExpiry = &t
					}
				}
				if v, ok := input["status"].(string); ok && v != "" {
					driver.Status = v
				}
				if v, ok := input["rating"].(float64); ok {
					driver.Rating = &v
				}
				if v, ok := input["totalDeliveries"].(int); ok {
					driver.TotalDeliveries = v
				}
				if v, ok := input["vehicleId"].(string); ok && v != "" {
					vid, parseErr := uuid.Parse(v)
					if parseErr != nil {
						return nil, fmt.Errorf("invalid vehicle id for driver: %w", parseErr)
					}
					// Validate the vehicle belongs to the same tenant (prevent IDOR).
					if _, err := r.VehicleRepo.GetByID(p.Context, tenantID, vid); err != nil {
						return nil, fmt.Errorf("vehicle not found in tenant")
					}
					driver.VehicleID = &vid
				}

				if err := r.DriverRepo.Create(p.Context, driver); err != nil {
					return nil, fmt.Errorf("failed to create driver: %w", err)
				}
				return driver, nil
			},
		},

		// -----------------------------------------------------------------
		// updateDriver
		// -----------------------------------------------------------------
		"updateDriver": &graphql.Field{
			Type:        types.DriverType,
			Description: "Update an existing driver.",
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.DriverInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid driver id: %w", err)
				}

				driver, err := r.DriverRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("driver not found: %w", err)
				}

				input := p.Args["input"].(map[string]interface{})

				if v, ok := input["employeeId"].(string); ok {
					driver.EmployeeID = v
				}
				if v, ok := input["firstName"].(string); ok {
					driver.FirstName = &v
				}
				if v, ok := input["lastName"].(string); ok {
					driver.LastName = &v
				}
				if v, ok := input["email"].(string); ok {
					driver.Email = &v
				}
				if v, ok := input["phone"].(string); ok {
					driver.Phone = &v
				}
				if v, ok := input["licenseNumber"].(string); ok {
					driver.LicenseNumber = &v
				}
				if v, ok := input["licenseExpiry"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						driver.LicenseExpiry = &t
					}
				}
				if v, ok := input["status"].(string); ok {
					driver.Status = v
				}
				if v, ok := input["rating"].(float64); ok {
					driver.Rating = &v
				}
				if v, ok := input["totalDeliveries"].(int); ok {
					driver.TotalDeliveries = v
				}
				if v, ok := input["vehicleId"].(string); ok && v != "" {
					vid, parseErr := uuid.Parse(v)
					if parseErr != nil {
						return nil, fmt.Errorf("invalid vehicle id for driver: %w", parseErr)
					}
					// Validate the vehicle belongs to the same tenant (prevent IDOR).
					if _, err := r.VehicleRepo.GetByID(p.Context, tenantID, vid); err != nil {
						return nil, fmt.Errorf("vehicle not found in tenant")
					}
					driver.VehicleID = &vid
				}

				driver.UpdatedAt = time.Now()

				if err := r.DriverRepo.Update(p.Context, tenantID, id, driver); err != nil {
					return nil, fmt.Errorf("failed to update driver: %w", err)
				}
				return driver, nil
			},
		},

		// -----------------------------------------------------------------
		// deleteDriver
		// -----------------------------------------------------------------
		"deleteDriver": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a driver by ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return false, fmt.Errorf("invalid driver id: %w", err)
				}

				if err := r.DriverRepo.Delete(p.Context, tenantID, id); err != nil {
					return false, fmt.Errorf("failed to delete driver: %w", err)
				}
				return true, nil
			},
		},

		// =================================================================
		// Maintenance mutations
		// =================================================================

		// -----------------------------------------------------------------
		// createMaintenance
		// -----------------------------------------------------------------
		"createMaintenance": &graphql.Field{
			Type:        types.MaintenanceType,
			Description: "Create a new maintenance record.",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.MaintenanceInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				vehicleID, err := uuid.Parse(input["vehicleId"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid vehicle id: %w", err)
				}

				// Validate the vehicle belongs to the same tenant (prevent IDOR).
				if _, err := r.VehicleRepo.GetByID(p.Context, tenantID, vehicleID); err != nil {
					return nil, fmt.Errorf("vehicle not found in tenant")
				}

				now := time.Now()
				record := &models.MaintenanceRecord{
					ID:        uuid.New(),
					TenantID:  tenantID,
					VehicleID: vehicleID,
					Status:    "scheduled",
					CreatedAt: now,
					UpdatedAt: now,
				}

				if v, ok := input["type"].(string); ok {
					record.Type = &v
				}
				if v, ok := input["description"].(string); ok {
					record.Description = &v
				}
				if v, ok := input["status"].(string); ok && v != "" {
					record.Status = v
				}
				if v, ok := input["scheduledDate"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						record.ScheduledDate = &t
					}
				}
				if v, ok := input["completedDate"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						record.CompletedDate = &t
					}
				}
				if v, ok := input["cost"].(float64); ok {
					record.Cost = &v
				}
				if v, ok := input["mechanic"].(string); ok {
					record.Mechanic = &v
				}

				if err := r.MaintenanceRepo.Create(p.Context, record); err != nil {
					return nil, fmt.Errorf("failed to create maintenance record: %w", err)
				}
				return record, nil
			},
		},

		// -----------------------------------------------------------------
		// updateMaintenance
		// -----------------------------------------------------------------
		"updateMaintenance": &graphql.Field{
			Type:        types.MaintenanceType,
			Description: "Update an existing maintenance record.",
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.MaintenanceInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid maintenance record id: %w", err)
				}

				record, err := r.MaintenanceRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("maintenance record not found: %w", err)
				}

				input := p.Args["input"].(map[string]interface{})

				if v, ok := input["vehicleId"].(string); ok && v != "" {
					vid, parseErr := uuid.Parse(v)
					if parseErr != nil {
						return nil, fmt.Errorf("invalid vehicle id: %w", parseErr)
					}
					// Validate the vehicle belongs to the same tenant (prevent IDOR).
					if _, err := r.VehicleRepo.GetByID(p.Context, tenantID, vid); err != nil {
						return nil, fmt.Errorf("vehicle not found in tenant")
					}
					record.VehicleID = vid
				}
				if v, ok := input["type"].(string); ok {
					record.Type = &v
				}
				if v, ok := input["description"].(string); ok {
					record.Description = &v
				}
				if v, ok := input["status"].(string); ok {
					record.Status = v
				}
				if v, ok := input["scheduledDate"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						record.ScheduledDate = &t
					}
				}
				if v, ok := input["completedDate"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						record.CompletedDate = &t
					}
				}
				if v, ok := input["cost"].(float64); ok {
					record.Cost = &v
				}
				if v, ok := input["mechanic"].(string); ok {
					record.Mechanic = &v
				}

				record.UpdatedAt = time.Now()

				if err := r.MaintenanceRepo.Update(p.Context, tenantID, id, record); err != nil {
					return nil, fmt.Errorf("failed to update maintenance record: %w", err)
				}
				return record, nil
			},
		},
	}
}
