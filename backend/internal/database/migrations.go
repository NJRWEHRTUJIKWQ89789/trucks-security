package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations creates all tables and indexes if they do not already exist.
func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	migrations := []string{
		// 1. tenants
		`CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			domain VARCHAR(255) UNIQUE,
			plan VARCHAR(50) DEFAULT 'starter',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 2. users
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			email VARCHAR(255) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			role VARCHAR(50) DEFAULT 'viewer',
			email_verified BOOLEAN DEFAULT FALSE,
			email_verify_token VARCHAR(255),
			avatar_url VARCHAR(500),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, email)
		)`,

		// 3. shipments
		`CREATE TABLE IF NOT EXISTS shipments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			tracking_number VARCHAR(50) NOT NULL,
			origin VARCHAR(255),
			destination VARCHAR(255),
			status VARCHAR(50) DEFAULT 'pending',
			carrier VARCHAR(100),
			weight DECIMAL(10,2),
			dimensions VARCHAR(100),
			estimated_delivery TIMESTAMPTZ,
			actual_delivery TIMESTAMPTZ,
			customer_name VARCHAR(255),
			customer_email VARCHAR(255),
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, tracking_number)
		)`,

		// 4. vehicles
		`CREATE TABLE IF NOT EXISTS vehicles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			vehicle_id VARCHAR(50) NOT NULL,
			name VARCHAR(255),
			type VARCHAR(50),
			status VARCHAR(50) DEFAULT 'available',
			fuel_level INTEGER DEFAULT 100,
			mileage INTEGER DEFAULT 0,
			last_service TIMESTAMPTZ,
			next_service TIMESTAMPTZ,
			license_plate VARCHAR(50),
			year INTEGER,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, vehicle_id)
		)`,

		// 5. drivers
		`CREATE TABLE IF NOT EXISTS drivers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			employee_id VARCHAR(50) NOT NULL,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			email VARCHAR(255),
			phone VARCHAR(50),
			license_number VARCHAR(100),
			license_expiry TIMESTAMPTZ,
			status VARCHAR(50) DEFAULT 'available',
			rating DECIMAL(3,2) DEFAULT 0,
			total_deliveries INTEGER DEFAULT 0,
			vehicle_id UUID REFERENCES vehicles(id),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, employee_id)
		)`,

		// 6. maintenance_records
		`CREATE TABLE IF NOT EXISTS maintenance_records (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			vehicle_id UUID NOT NULL REFERENCES vehicles(id),
			type VARCHAR(100),
			description TEXT,
			status VARCHAR(50) DEFAULT 'scheduled',
			scheduled_date TIMESTAMPTZ,
			completed_date TIMESTAMPTZ,
			cost DECIMAL(10,2),
			mechanic VARCHAR(255),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 7. warehouses
		`CREATE TABLE IF NOT EXISTS warehouses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			location VARCHAR(255),
			address TEXT,
			capacity INTEGER DEFAULT 0,
			used_capacity INTEGER DEFAULT 0,
			manager VARCHAR(255),
			phone VARCHAR(50),
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 8. inventory_items
		`CREATE TABLE IF NOT EXISTS inventory_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			warehouse_id UUID NOT NULL REFERENCES warehouses(id),
			sku VARCHAR(100) NOT NULL,
			name VARCHAR(255),
			category VARCHAR(100),
			quantity INTEGER DEFAULT 0,
			min_quantity INTEGER DEFAULT 0,
			unit_price DECIMAL(10,2),
			weight DECIMAL(10,2),
			status VARCHAR(50) DEFAULT 'in_stock',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, sku)
		)`,

		// 9. orders
		`CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			order_number VARCHAR(50) NOT NULL,
			customer_name VARCHAR(255),
			customer_email VARCHAR(255),
			status VARCHAR(50) DEFAULT 'pending',
			type VARCHAR(50) DEFAULT 'standard',
			total_amount DECIMAL(10,2),
			shipment_id UUID REFERENCES shipments(id),
			scheduled_date TIMESTAMPTZ,
			return_reason TEXT,
			cancellation_reason TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, order_number)
		)`,

		// 10. vendors
		`CREATE TABLE IF NOT EXISTS vendors (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			contact_person VARCHAR(255),
			email VARCHAR(255),
			phone VARCHAR(50),
			address TEXT,
			category VARCHAR(100),
			rating DECIMAL(3,2) DEFAULT 0,
			contract_start TIMESTAMPTZ,
			contract_end TIMESTAMPTZ,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 11. clients
		`CREATE TABLE IF NOT EXISTS clients (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			company_name VARCHAR(255) NOT NULL,
			contact_person VARCHAR(255),
			email VARCHAR(255),
			phone VARCHAR(50),
			address TEXT,
			industry VARCHAR(100),
			total_shipments INTEGER DEFAULT 0,
			total_spent DECIMAL(12,2) DEFAULT 0,
			satisfaction_rating DECIMAL(3,2) DEFAULT 0,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 12. client_feedback
		`CREATE TABLE IF NOT EXISTS client_feedback (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			client_id UUID NOT NULL REFERENCES clients(id),
			rating INTEGER CHECK (rating >= 1 AND rating <= 5),
			comment TEXT,
			category VARCHAR(100),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 13. notifications
		`CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id),
			title VARCHAR(255),
			message TEXT,
			type VARCHAR(50),
			read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 14. roles
		`CREATE TABLE IF NOT EXISTS roles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			permissions JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, name)
		)`,

		// 15. settings
		`CREATE TABLE IF NOT EXISTS settings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			category VARCHAR(100),
			updated_by UUID REFERENCES users(id),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tenant_id, key)
		)`,

		// 16. activity_log
		`CREATE TABLE IF NOT EXISTS activity_log (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id),
			action VARCHAR(255),
			entity_type VARCHAR(100),
			entity_id UUID,
			details JSONB,
			ip_address VARCHAR(45),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// 17. notification_preferences
		`CREATE TABLE IF NOT EXISTS notification_preferences (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id),
			event_type VARCHAR(100),
			email_enabled BOOLEAN DEFAULT TRUE,
			sms_enabled BOOLEAN DEFAULT FALSE,
			push_enabled BOOLEAN DEFAULT FALSE,
			UNIQUE(tenant_id, user_id, event_type)
		)`,

		// Indexes on tenant_id for all tables
		`CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_shipments_tenant_id ON shipments(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_vehicles_tenant_id ON vehicles(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_drivers_tenant_id ON drivers(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_maintenance_records_tenant_id ON maintenance_records(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_warehouses_tenant_id ON warehouses(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_items_tenant_id ON inventory_items(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_vendors_tenant_id ON vendors(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_clients_tenant_id ON clients(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_client_feedback_tenant_id ON client_feedback(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notification_preferences_tenant_id ON notification_preferences(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_roles_tenant_id ON roles(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_settings_tenant_id ON settings(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activity_log_tenant_id ON activity_log(tenant_id)`,

		// 18. Add pin_hash column to drivers table
		`ALTER TABLE drivers ADD COLUMN IF NOT EXISTS pin_hash VARCHAR(255)`,

		// 19. gps_pings (high-volume)
		`CREATE TABLE IF NOT EXISTS gps_pings (
			id BIGSERIAL PRIMARY KEY,
			tenant_id UUID NOT NULL REFERENCES tenants(id),
			driver_id UUID NOT NULL REFERENCES drivers(id),
			truck_id UUID NOT NULL,
			shift_id UUID NOT NULL,
			latitude DECIMAL(10,7) NOT NULL,
			longitude DECIMAL(10,7) NOT NULL,
			speed_kmh DECIMAL(5,1) DEFAULT 0,
			heading SMALLINT DEFAULT 0 CHECK (heading >= 0 AND heading <= 360),
			accuracy DECIMAL(5,1) DEFAULT 0,
			battery_level SMALLINT DEFAULT 0 CHECK (battery_level >= 0 AND battery_level <= 100),
			is_moving BOOLEAN DEFAULT false,
			recorded_at TIMESTAMPTZ NOT NULL,
			received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			is_delayed BOOLEAN DEFAULT false,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// 20. shifts
		`CREATE TABLE IF NOT EXISTS shifts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id),
			driver_id UUID NOT NULL REFERENCES drivers(id),
			truck_id UUID NOT NULL,
			started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			ended_at TIMESTAMPTZ,
			status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed')),
			total_km DECIMAL(8,2) DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// 21. alerts
		`CREATE TABLE IF NOT EXISTS alerts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id),
			driver_id UUID NOT NULL REFERENCES drivers(id),
			shift_id UUID,
			type VARCHAR(30) NOT NULL CHECK (type IN ('unauthorized_stop', 'driver_offline', 'speed_exceeded')),
			status VARCHAR(20) NOT NULL DEFAULT 'triggered' CHECK (status IN ('triggered', 'notified', 'acknowledged', 'resolved', 'false_alarm')),
			stop_latitude DECIMAL(10,7),
			stop_longitude DECIMAL(10,7),
			stop_duration_seconds INT DEFAULT 0,
			nearest_zone_id UUID,
			nearest_zone_distance_meters DECIMAL(10,2),
			manager_notes TEXT,
			triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			notified_at TIMESTAMPTZ,
			acknowledged_at TIMESTAMPTZ,
			resolved_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// 22. approved_zones
		`CREATE TABLE IF NOT EXISTS approved_zones (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id),
			label VARCHAR(255) NOT NULL,
			latitude DECIMAL(10,7) NOT NULL,
			longitude DECIMAL(10,7) NOT NULL,
			radius_meters INT NOT NULL DEFAULT 500,
			type VARCHAR(30) DEFAULT 'other' CHECK (type IN ('warehouse', 'client_site', 'gas_station', 'rest_area', 'other')),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// 23. alert_config (one per company)
		`CREATE TABLE IF NOT EXISTS alert_config (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(id),
			max_stop_duration_minutes INT NOT NULL DEFAULT 5,
			alert_on_driver_offline BOOLEAN NOT NULL DEFAULT true,
			offline_threshold_minutes INT NOT NULL DEFAULT 3,
			notify_via_push BOOLEAN NOT NULL DEFAULT true,
			notify_via_email BOOLEAN NOT NULL DEFAULT true,
			notify_via_sms BOOLEAN NOT NULL DEFAULT false,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// Indexes for new tables
		`CREATE INDEX IF NOT EXISTS idx_gps_pings_driver_time ON gps_pings(tenant_id, driver_id, recorded_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_gps_pings_company_time ON gps_pings(tenant_id, recorded_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_gps_pings_shift ON gps_pings(shift_id, recorded_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_shifts_tenant ON shifts(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_shifts_driver ON shifts(tenant_id, driver_id, started_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_shifts_active ON shifts(tenant_id, status) WHERE status = 'active'`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_tenant ON alerts(tenant_id, triggered_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(tenant_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_zones_tenant ON approved_zones(tenant_id)`,
	}

	for i, migration := range migrations {
		if _, err := pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}
