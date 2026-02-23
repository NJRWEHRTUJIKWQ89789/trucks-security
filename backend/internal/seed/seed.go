package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// SeedData populates the database with realistic demo data for two tenants.
// All inserts use ON CONFLICT DO NOTHING so the function is idempotent.
func SeedData(pool *pgxpool.Pool) error {
	ctx := context.Background()

	log.Println("Seeding database with demo data...")

	// Hash the shared demo password once at startup.
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seed: bcrypt hash failed: %w", err)
	}
	pwHash := string(passwordHash)

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// ---------------------------------------------------------------
	// Tenant IDs (deterministic UUIDs so we can reference them)
	// ---------------------------------------------------------------
	acmeTenantID := uuid.New()
	betaTenantID := uuid.New()

	// ---------------------------------------------------------------
	// 1. TENANTS
	// ---------------------------------------------------------------
	tenantSQL := `INSERT INTO tenants (id, name, domain, plan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING`

	if _, err := pool.Exec(ctx, tenantSQL, acmeTenantID, "Acme Logistics", "acme-logistics.com", "enterprise", now.Add(-180*24*time.Hour), now); err != nil {
		return fmt.Errorf("seed tenants: %w", err)
	}
	if _, err := pool.Exec(ctx, tenantSQL, betaTenantID, "Beta Transport", "betatransport.io", "starter", now.Add(-60*24*time.Hour), now); err != nil {
		return fmt.Errorf("seed tenants: %w", err)
	}

	// ---------------------------------------------------------------
	// 2. USERS — Acme (5) + Beta (2)
	// ---------------------------------------------------------------
	type seedUser struct {
		id        uuid.UUID
		tenantID  uuid.UUID
		email     string
		firstName string
		lastName  string
		role      string
	}

	acmeAdmin := seedUser{uuid.New(), acmeTenantID, "admin@acme.com", "Alice", "Anderson", "admin"}
	acmeManager := seedUser{uuid.New(), acmeTenantID, "john@acme.com", "John", "Mitchell", "manager"}
	acmeDispatcher := seedUser{uuid.New(), acmeTenantID, "sarah@acme.com", "Sarah", "Chen", "dispatcher"}
	acmeDriver := seedUser{uuid.New(), acmeTenantID, "mike@acme.com", "Mike", "Rodriguez", "driver"}
	acmeViewer := seedUser{uuid.New(), acmeTenantID, "viewer@acme.com", "Emma", "Davis", "viewer"}

	betaAdmin := seedUser{uuid.New(), betaTenantID, "admin@beta.com", "Robert", "Kim", "admin"}
	betaViewer := seedUser{uuid.New(), betaTenantID, "viewer@beta.com", "Lisa", "Park", "viewer"}

	allUsers := []seedUser{acmeAdmin, acmeManager, acmeDispatcher, acmeDriver, acmeViewer, betaAdmin, betaViewer}

	userSQL := `INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, $8, $9)
		ON CONFLICT (tenant_id, email) DO NOTHING`

	for _, u := range allUsers {
		if _, err := pool.Exec(ctx, userSQL, u.id, u.tenantID, u.email, pwHash, u.firstName, u.lastName, u.role, now.Add(-90*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed users (%s): %w", u.email, err)
		}
	}

	// ---------------------------------------------------------------
	// 3. VEHICLES — Acme (45) + Beta (5)
	// ---------------------------------------------------------------
	type seedVehicle struct {
		id       uuid.UUID
		tenantID uuid.UUID
		vid      string
		name     string
		vtype    string
		status   string
		fuel     int
		mileage  int
		plate    string
		year     int
	}

	truckNames := []string{
		"Freightliner Cascadia", "Kenworth T680", "Peterbilt 579", "Volvo VNL 860",
		"Mack Anthem", "International LT", "Western Star 5700XE", "Kenworth W990",
		"Peterbilt 389", "Freightliner M2 106", "Hino L6", "Isuzu NPR-HD",
		"Ford F-650", "Ram 5500", "Chevrolet Silverado 6500HD", "Kenworth T880",
		"Peterbilt 567", "Mack Granite", "Volvo VHD", "International HV",
		"Freightliner 114SD", "Western Star 4700SF", "Hino XL8", "Isuzu FTR",
		"Ford F-750", "Freightliner eCascadia", "Volvo VNR Electric", "Kenworth T680E",
		"Peterbilt 220EV", "Mack LR Electric", "International eMV", "Ford E-Transit",
		"Freightliner M2 112", "Kenworth T370", "Peterbilt 536", "Volvo VNL 300",
		"Mack MD Electric", "International CV", "Western Star 4900", "Hino M5",
		"Isuzu NRR", "Ford F-600", "Ram 4500", "Chevrolet Low Cab Forward 5500HD",
		"Kenworth K270E",
	}

	vehicleStatuses := []string{"available", "in_transit", "maintenance", "out_of_service", "available", "in_transit", "available", "in_transit", "available"}
	vehicleTypes := []string{"semi_truck", "box_truck", "flatbed", "refrigerated", "tanker", "semi_truck", "box_truck", "flatbed", "semi_truck"}

	acmeVehicles := make([]seedVehicle, 45)
	for i := 0; i < 45; i++ {
		status := vehicleStatuses[i%len(vehicleStatuses)]
		fuel := 30 + (i*17)%71 // 30-100
		mileage := 12000 + (i * 3847) % 188000
		year := 2019 + i%6 // 2019-2024
		acmeVehicles[i] = seedVehicle{
			id:       uuid.New(),
			tenantID: acmeTenantID,
			vid:      fmt.Sprintf("TRK-%03d", i+1),
			name:     truckNames[i],
			vtype:    vehicleTypes[i%len(vehicleTypes)],
			status:   status,
			fuel:     fuel,
			mileage:  mileage,
			plate:    fmt.Sprintf("%s-%s-%04d", string(rune('A'+i%26)), string(rune('A'+(i+7)%26)), 1000+i*73%9000),
			year:     year,
		}
	}

	betaVehicleNames := []string{"Freightliner Cascadia", "Kenworth T680", "Volvo VNL 860", "Peterbilt 579", "Mack Anthem"}
	betaVehicles := make([]seedVehicle, 5)
	for i := 0; i < 5; i++ {
		betaVehicles[i] = seedVehicle{
			id:       uuid.New(),
			tenantID: betaTenantID,
			vid:      fmt.Sprintf("BT-%03d", i+1),
			name:     betaVehicleNames[i],
			vtype:    vehicleTypes[i],
			status:   vehicleStatuses[i],
			fuel:     50 + i*10,
			mileage:  20000 + i*15000,
			plate:    fmt.Sprintf("BT-%04d", 5000+i),
			year:     2021 + i%3,
		}
	}

	vehicleSQL := `INSERT INTO vehicles (id, tenant_id, vehicle_id, name, type, status, fuel_level, mileage, last_service, next_service, license_plate, year, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (tenant_id, vehicle_id) DO NOTHING`

	allVehicles := append(acmeVehicles, betaVehicles...)
	for _, v := range allVehicles {
		lastService := now.Add(-time.Duration(30+v.mileage%120) * 24 * time.Hour)
		nextService := now.Add(time.Duration(30+v.fuel*2) * 24 * time.Hour)
		if _, err := pool.Exec(ctx, vehicleSQL, v.id, v.tenantID, v.vid, v.name, v.vtype, v.status, v.fuel, v.mileage, lastService, nextService, v.plate, v.year, now.Add(-120*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed vehicles (%s): %w", v.vid, err)
		}
	}

	// ---------------------------------------------------------------
	// 4. DRIVERS — Acme (30) + Beta (3)
	// ---------------------------------------------------------------
	type seedDriver struct {
		id        uuid.UUID
		tenantID  uuid.UUID
		empID     string
		firstName string
		lastName  string
		email     string
		phone     string
		status    string
		rating    float64
		delivers  int
		vehicleID *uuid.UUID
	}

	driverFirstNames := []string{
		"James", "Robert", "David", "Carlos", "Marcus", "Antonio", "William", "Thomas",
		"Daniel", "Richard", "Joseph", "Kevin", "Brian", "Steven", "Edward",
		"Alejandro", "Terrence", "Gregory", "Raymond", "Lawrence",
		"Christopher", "Patrick", "Timothy", "Jeffrey", "Frank",
		"Samuel", "Dennis", "Henry", "Peter", "Eugene",
	}
	driverLastNames := []string{
		"Thompson", "Garcia", "Williams", "Martinez", "Johnson", "Rossi", "O'Brien", "Washington",
		"Lee", "Nguyen", "Brown", "Taylor", "Anderson", "Jackson", "White",
		"Hernandez", "Robinson", "Clark", "Lewis", "Walker",
		"Young", "King", "Scott", "Green", "Baker",
		"Adams", "Nelson", "Carter", "Mitchell", "Perez",
	}
	driverStatuses := []string{"available", "on_trip", "off_duty", "available", "on_trip", "available"}

	acmeDrivers := make([]seedDriver, 30)
	for i := 0; i < 30; i++ {
		rating := 3.5 + float64(i%16)*0.1 // 3.5 - 5.0
		if rating > 5.0 {
			rating = 5.0
		}
		var vehiclePtr *uuid.UUID
		if i < 20 { // first 20 drivers assigned to vehicles
			vid := acmeVehicles[i].id
			vehiclePtr = &vid
		}
		acmeDrivers[i] = seedDriver{
			id:        uuid.New(),
			tenantID:  acmeTenantID,
			empID:     fmt.Sprintf("DRV-%03d", i+1),
			firstName: driverFirstNames[i],
			lastName:  driverLastNames[i],
			email:     fmt.Sprintf("%s.%s@acme.com", driverFirstNames[i], driverLastNames[i]),
			phone:     fmt.Sprintf("(555) %03d-%04d", 100+i*31%900, 1000+i*137%9000),
			status:    driverStatuses[i%len(driverStatuses)],
			rating:    rating,
			delivers:  80 + (i*47)%420,
			vehicleID: vehiclePtr,
		}
	}

	betaDriverFirstNames := []string{"Marco", "Diana", "Oscar"}
	betaDriverLastNames := []string{"Reyes", "Foster", "Chen"}
	betaDrivers := make([]seedDriver, 3)
	for i := 0; i < 3; i++ {
		vid := betaVehicles[i].id
		betaDrivers[i] = seedDriver{
			id:        uuid.New(),
			tenantID:  betaTenantID,
			empID:     fmt.Sprintf("BD-%03d", i+1),
			firstName: betaDriverFirstNames[i],
			lastName:  betaDriverLastNames[i],
			email:     fmt.Sprintf("%s.%s@betatransport.io", betaDriverFirstNames[i], betaDriverLastNames[i]),
			phone:     fmt.Sprintf("(555) %03d-%04d", 600+i, 2000+i*300),
			status:    driverStatuses[i],
			rating:    4.0 + float64(i)*0.3,
			delivers:  50 + i*40,
			vehicleID: &vid,
		}
	}

	driverSQL := `INSERT INTO drivers (id, tenant_id, employee_id, first_name, last_name, email, phone, license_number, license_expiry, status, rating, total_deliveries, vehicle_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (tenant_id, employee_id) DO NOTHING`

	allDrivers := append(acmeDrivers, betaDrivers...)
	for i, d := range allDrivers {
		licenseNum := fmt.Sprintf("CDL-%06d", 100000+i*7919)
		licenseExpiry := now.Add(time.Duration(90+i*30) * 24 * time.Hour)
		if _, err := pool.Exec(ctx, driverSQL, d.id, d.tenantID, d.empID, d.firstName, d.lastName, d.email, d.phone, licenseNum, licenseExpiry, d.status, d.rating, d.delivers, d.vehicleID, now.Add(-100*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed drivers (%s): %w", d.empID, err)
		}
	}

	// ---------------------------------------------------------------
	// 5. SHIPMENTS — Acme (42) + Beta (10)
	// ---------------------------------------------------------------
	type seedShipment struct {
		id           uuid.UUID
		tenantID     uuid.UUID
		tracking     string
		origin       string
		destination  string
		status       string
		carrier      string
		weight       float64
		dimensions   string
		estDelivery  time.Time
		actDelivery  *time.Time
		custName     string
		custEmail    string
		notes        string
	}

	usCities := []string{
		"Los Angeles, CA", "New York, NY", "Chicago, IL", "Houston, TX",
		"Phoenix, AZ", "Philadelphia, PA", "San Antonio, TX", "San Diego, CA",
		"Dallas, TX", "San Jose, CA", "Austin, TX", "Jacksonville, FL",
		"Fort Worth, TX", "Columbus, OH", "Charlotte, NC", "Indianapolis, IN",
		"San Francisco, CA", "Seattle, WA", "Denver, CO", "Nashville, TN",
		"Oklahoma City, OK", "Portland, OR", "Las Vegas, NV", "Memphis, TN",
		"Louisville, KY", "Baltimore, MD", "Milwaukee, WI", "Albuquerque, NM",
		"Tucson, AZ", "Fresno, CA", "Sacramento, CA", "Mesa, AZ",
		"Kansas City, MO", "Atlanta, GA", "Omaha, NE", "Miami, FL",
		"Minneapolis, MN", "Raleigh, NC", "Cleveland, OH", "Tampa, FL",
		"St. Louis, MO", "Pittsburgh, PA",
	}

	carriers := []string{
		"Acme Fleet", "FastFreight Express", "Continental Hauling", "Pacific Cargo Lines",
		"Midwest Logistics", "Southern Transport Co", "Atlantic Freight", "Mountain Express",
	}

	// We ensure 7 delayed, 28 delivered today. Build statuses explicitly.
	acmeShipmentStatuses := make([]string, 42)
	for i := 0; i < 42; i++ {
		switch {
		case i < 28:
			acmeShipmentStatuses[i] = "delivered"
		case i < 33:
			acmeShipmentStatuses[i] = "in_transit"
		case i < 35:
			acmeShipmentStatuses[i] = "pending"
		default: // 35-41 = 7 delayed
			acmeShipmentStatuses[i] = "delayed"
		}
	}

	customerFirstNames := []string{"Michael", "Jennifer", "David", "Sarah", "Robert", "Emily", "William", "Jessica", "Richard", "Amanda", "Charles", "Stephanie", "Joseph", "Nicole", "Thomas"}
	customerLastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Davis", "Miller", "Wilson", "Moore", "Taylor", "Anderson", "Thomas", "Jackson", "White", "Harris"}

	acmeShipments := make([]seedShipment, 42)
	for i := 0; i < 42; i++ {
		status := acmeShipmentStatuses[i]
		origin := usCities[i%len(usCities)]
		dest := usCities[(i+7)%len(usCities)]
		weight := 150.0 + float64(i*73%4850)
		custFirst := customerFirstNames[i%len(customerFirstNames)]
		custLast := customerLastNames[i%len(customerLastNames)]

		est := today.Add(time.Duration(i%5) * 24 * time.Hour)
		var act *time.Time
		if status == "delivered" {
			// 28 delivered today
			t := today.Add(time.Duration(8+i%10) * time.Hour) // delivered between 8am-5pm today
			act = &t
		}
		if status == "delayed" {
			est = today.Add(-time.Duration(1+i%3) * 24 * time.Hour) // past due
		}

		acmeShipments[i] = seedShipment{
			id:          uuid.New(),
			tenantID:    acmeTenantID,
			tracking:    fmt.Sprintf("SHP-%d", 1001+i),
			origin:      origin,
			destination: dest,
			status:      status,
			carrier:     carriers[i%len(carriers)],
			weight:      weight,
			dimensions:  fmt.Sprintf("%dx%dx%d", 24+i%48, 20+i%36, 18+i%30),
			estDelivery: est,
			actDelivery: act,
			custName:    fmt.Sprintf("%s %s", custFirst, custLast),
			custEmail:   fmt.Sprintf("%s.%s%d@example.com", custFirst, custLast, i),
			notes:       fmt.Sprintf("Priority shipment from %s to %s", origin, dest),
		}
	}

	betaShipmentStatuses := []string{"delivered", "delivered", "delivered", "delivered", "in_transit", "in_transit", "pending", "pending", "delayed", "delayed"}
	betaShipments := make([]seedShipment, 10)
	for i := 0; i < 10; i++ {
		status := betaShipmentStatuses[i]
		est := today.Add(time.Duration(i%4) * 24 * time.Hour)
		var act *time.Time
		if status == "delivered" {
			t := today.Add(time.Duration(-i) * 24 * time.Hour)
			act = &t
		}
		betaShipments[i] = seedShipment{
			id:          uuid.New(),
			tenantID:    betaTenantID,
			tracking:    fmt.Sprintf("BSH-%d", 2001+i),
			origin:      usCities[i*3%len(usCities)],
			destination: usCities[(i*3+10)%len(usCities)],
			status:      status,
			carrier:     carriers[i%4],
			weight:      200.0 + float64(i)*150,
			dimensions:  fmt.Sprintf("%dx%dx%d", 30+i*5, 24+i*3, 20+i*2),
			estDelivery: est,
			actDelivery: act,
			custName:    fmt.Sprintf("%s %s", customerFirstNames[i%5], customerLastNames[i%5+5]),
			custEmail:   fmt.Sprintf("beta.customer%d@example.com", i+1),
			notes:       fmt.Sprintf("Beta Transport shipment #%d", i+1),
		}
	}

	shipmentSQL := `INSERT INTO shipments (id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (tenant_id, tracking_number) DO NOTHING`

	allShipments := append(acmeShipments, betaShipments...)
	for _, s := range allShipments {
		createdAt := now.Add(-time.Duration(5+int(s.weight)%10) * 24 * time.Hour)
		if _, err := pool.Exec(ctx, shipmentSQL, s.id, s.tenantID, s.tracking, s.origin, s.destination, s.status, s.carrier, s.weight, s.dimensions, s.estDelivery, s.actDelivery, s.custName, s.custEmail, s.notes, createdAt, now); err != nil {
			return fmt.Errorf("seed shipments (%s): %w", s.tracking, err)
		}
	}

	// ---------------------------------------------------------------
	// 6. MAINTENANCE RECORDS — Acme (20)
	// ---------------------------------------------------------------
	maintenanceTypes := []string{
		"Oil Change", "Tire Rotation", "Brake Inspection", "Engine Tune-up",
		"Transmission Service", "Air Filter Replacement", "Coolant Flush", "Battery Replacement",
		"Wheel Alignment", "Fuel System Cleaning", "Suspension Check", "Electrical Diagnostic",
		"DOT Inspection", "Exhaust System Repair", "A/C Service", "Windshield Replacement",
		"Clutch Replacement", "Differential Service", "Power Steering Flush", "Belt Replacement",
	}
	maintenanceStatuses := []string{"completed", "completed", "completed", "in_progress", "scheduled", "completed", "scheduled", "completed", "completed", "in_progress"}
	mechanics := []string{"Tony's Truck Shop", "FleetCare Mechanics", "Highway Diesel Service", "AllPro Truck Repair", "Summit Fleet Maintenance"}

	maintenanceSQL := `INSERT INTO maintenance_records (id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT DO NOTHING`

	for i := 0; i < 20; i++ {
		mID := uuid.New()
		vehicle := acmeVehicles[i%len(acmeVehicles)]
		mType := maintenanceTypes[i]
		mStatus := maintenanceStatuses[i%len(maintenanceStatuses)]
		schedDate := now.Add(-time.Duration(60-i*3) * 24 * time.Hour)
		var compDate *time.Time
		cost := 150.0 + float64(i*97%1850)
		if mStatus == "completed" {
			t := schedDate.Add(time.Duration(1+i%3) * 24 * time.Hour)
			compDate = &t
		}
		desc := fmt.Sprintf("Routine %s for %s (VIN: %s)", mType, vehicle.name, vehicle.vid)
		if _, err := pool.Exec(ctx, maintenanceSQL, mID, acmeTenantID, vehicle.id, mType, desc, mStatus, schedDate, compDate, cost, mechanics[i%len(mechanics)], now.Add(-60*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed maintenance (%d): %w", i, err)
		}
	}

	// ---------------------------------------------------------------
	// 7. WAREHOUSES — Acme (8) + Beta (2)
	// ---------------------------------------------------------------
	type seedWarehouse struct {
		id       uuid.UUID
		tenantID uuid.UUID
		name     string
		location string
		address  string
		capacity int
		used     int
		manager  string
		phone    string
		status   string
	}

	acmeWarehouses := []seedWarehouse{
		{uuid.New(), acmeTenantID, "LA Distribution Center", "Los Angeles, CA", "1200 Alameda St, Los Angeles, CA 90012", 50000, 38500, "Robert Chang", "(213) 555-0101", "active"},
		{uuid.New(), acmeTenantID, "Chicago Hub", "Chicago, IL", "4500 W Roosevelt Rd, Chicago, IL 60624", 45000, 31200, "Maria Santos", "(312) 555-0202", "active"},
		{uuid.New(), acmeTenantID, "Dallas Mega Warehouse", "Dallas, TX", "8900 Stemmons Fwy, Dallas, TX 75247", 60000, 42000, "Tom Bradley", "(214) 555-0303", "active"},
		{uuid.New(), acmeTenantID, "Atlanta Fulfillment", "Atlanta, GA", "2100 Donald Lee Hollowell Pkwy, Atlanta, GA 30318", 35000, 28700, "Keisha Williams", "(404) 555-0404", "active"},
		{uuid.New(), acmeTenantID, "Seattle Cold Storage", "Seattle, WA", "3700 E Marginal Way S, Seattle, WA 98134", 25000, 19800, "Derek Olson", "(206) 555-0505", "active"},
		{uuid.New(), acmeTenantID, "Miami Port Warehouse", "Miami, FL", "1000 Port Blvd, Miami, FL 33132", 40000, 35600, "Carlos Mendez", "(305) 555-0606", "active"},
		{uuid.New(), acmeTenantID, "Denver Transit Hub", "Denver, CO", "5600 E 56th Ave, Commerce City, CO 80022", 30000, 18900, "Angela Torres", "(720) 555-0707", "active"},
		{uuid.New(), acmeTenantID, "NYC Metro Facility", "Newark, NJ", "100 Port St, Newark, NJ 07114", 55000, 48200, "Frank DeLuca", "(973) 555-0808", "active"},
	}

	betaWarehouses := []seedWarehouse{
		{uuid.New(), betaTenantID, "Beta Central Depot", "Houston, TX", "7200 Navigation Blvd, Houston, TX 77011", 15000, 8500, "Marco Reyes", "(713) 555-0901", "active"},
		{uuid.New(), betaTenantID, "Beta East Hub", "Charlotte, NC", "4800 Statesville Ave, Charlotte, NC 28269", 12000, 6200, "Diana Foster", "(704) 555-0902", "active"},
	}

	warehouseSQL := `INSERT INTO warehouses (id, tenant_id, name, location, address, capacity, used_capacity, manager, phone, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT DO NOTHING`

	allWarehouses := append(acmeWarehouses, betaWarehouses...)
	for _, w := range allWarehouses {
		if _, err := pool.Exec(ctx, warehouseSQL, w.id, w.tenantID, w.name, w.location, w.address, w.capacity, w.used, w.manager, w.phone, w.status, now.Add(-150*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed warehouses (%s): %w", w.name, err)
		}
	}

	// ---------------------------------------------------------------
	// 8. INVENTORY ITEMS — Acme (50)
	// ---------------------------------------------------------------
	type seedInventory struct {
		sku      string
		name     string
		category string
		quantity int
		minQty   int
		price    float64
		weight   float64
		status   string
	}

	inventoryItems := []seedInventory{
		{"SKU-001", "Standard Cardboard Box 24x18x12", "Packaging", 1200, 200, 2.50, 0.8, "in_stock"},
		{"SKU-002", "Heavy-Duty Pallet Wrap", "Packaging", 800, 100, 18.99, 3.2, "in_stock"},
		{"SKU-003", "Bubble Wrap Roll 100ft", "Packaging", 450, 50, 24.99, 5.0, "in_stock"},
		{"SKU-004", "Shipping Labels (1000pk)", "Supplies", 300, 50, 34.99, 1.5, "in_stock"},
		{"SKU-005", "Packing Tape 6-Pack", "Supplies", 620, 100, 12.99, 2.1, "in_stock"},
		{"SKU-006", "Foam Packing Peanuts 14cu ft", "Packaging", 180, 30, 42.99, 8.0, "in_stock"},
		{"SKU-007", "Wooden Pallet 48x40", "Equipment", 540, 100, 15.00, 22.0, "in_stock"},
		{"SKU-008", "Stretch Film 18in", "Packaging", 720, 80, 28.50, 4.5, "in_stock"},
		{"SKU-009", "Corrugated Mailer 12x9x4", "Packaging", 2400, 500, 1.25, 0.3, "in_stock"},
		{"SKU-010", "Thermal Printer Ribbon", "Supplies", 95, 20, 45.00, 0.5, "low_stock"},
		{"SKU-011", "Cargo Straps 2in x 27ft", "Equipment", 160, 30, 19.99, 1.8, "in_stock"},
		{"SKU-012", "Moving Blankets 72x80", "Equipment", 240, 40, 32.00, 5.5, "in_stock"},
		{"SKU-013", "Hand Truck Dolly", "Equipment", 35, 10, 89.99, 18.0, "in_stock"},
		{"SKU-014", "Pallet Jack Manual", "Equipment", 12, 3, 329.99, 75.0, "in_stock"},
		{"SKU-015", "Safety Vest High-Vis", "Safety", 480, 100, 8.99, 0.3, "in_stock"},
		{"SKU-016", "Work Gloves (12 pair)", "Safety", 360, 60, 24.99, 1.2, "in_stock"},
		{"SKU-017", "Floor Marking Tape Yellow", "Supplies", 140, 20, 16.99, 0.8, "in_stock"},
		{"SKU-018", "Barcode Scanner Handheld", "Technology", 28, 5, 249.99, 0.4, "in_stock"},
		{"SKU-019", "Insulated Shipping Box", "Packaging", 75, 25, 8.50, 1.2, "low_stock"},
		{"SKU-020", "Dry Ice Packs (24pk)", "Cold Chain", 90, 30, 36.00, 12.0, "in_stock"},
		{"SKU-021", "Hazmat Labels Assorted", "Compliance", 550, 100, 22.00, 0.5, "in_stock"},
		{"SKU-022", "Bill of Lading Forms (500)", "Documentation", 180, 30, 29.99, 2.0, "in_stock"},
		{"SKU-023", "Ratchet Tie Down 1.5in", "Equipment", 320, 50, 14.50, 1.5, "in_stock"},
		{"SKU-024", "Dock Bumper 10x4.5x10", "Equipment", 24, 6, 65.00, 9.0, "in_stock"},
		{"SKU-025", "Warehouse Fan 24in", "Equipment", 16, 4, 159.99, 14.0, "in_stock"},
		{"SKU-026", "LED Bay Light 150W", "Equipment", 60, 10, 79.99, 3.8, "in_stock"},
		{"SKU-027", "Dunnage Air Bags", "Packaging", 400, 80, 5.25, 0.6, "in_stock"},
		{"SKU-028", "Temperature Logger Digital", "Cold Chain", 42, 10, 55.00, 0.2, "in_stock"},
		{"SKU-029", "Poly Strapping Kit", "Packaging", 110, 20, 38.00, 6.0, "in_stock"},
		{"SKU-030", "Corner Protectors (100pk)", "Packaging", 280, 40, 18.00, 4.0, "in_stock"},
		{"SKU-031", "Warehouse Shelving Unit", "Equipment", 18, 5, 189.99, 45.0, "in_stock"},
		{"SKU-032", "Anti-Fatigue Floor Mat", "Safety", 65, 15, 42.00, 6.5, "in_stock"},
		{"SKU-033", "Pallet Collar Wooden", "Equipment", 200, 40, 12.00, 4.0, "in_stock"},
		{"SKU-034", "Shrink Wrap Gun", "Equipment", 8, 2, 275.00, 3.5, "in_stock"},
		{"SKU-035", "Spill Containment Kit", "Safety", 22, 5, 145.00, 15.0, "in_stock"},
		{"SKU-036", "First Aid Kit Industrial", "Safety", 30, 8, 68.00, 3.0, "in_stock"},
		{"SKU-037", "Fire Extinguisher 10lb", "Safety", 45, 10, 55.00, 10.0, "in_stock"},
		{"SKU-038", "Dock Plate Aluminum 60x36", "Equipment", 6, 2, 425.00, 38.0, "in_stock"},
		{"SKU-039", "Label Printer Thermal", "Technology", 15, 3, 399.99, 4.5, "in_stock"},
		{"SKU-040", "Cargo Net 6x8ft", "Equipment", 85, 15, 52.00, 3.0, "in_stock"},
		{"SKU-041", "VCI Paper Rust Prevention", "Packaging", 130, 25, 32.00, 7.0, "in_stock"},
		{"SKU-042", "Moisture Barrier Bag 48x42", "Packaging", 350, 60, 4.75, 0.4, "in_stock"},
		{"SKU-043", "Forklift Battery 36V", "Equipment", 4, 1, 1850.00, 550.0, "low_stock"},
		{"SKU-044", "Conveyor Belt Section 10ft", "Equipment", 3, 1, 2200.00, 65.0, "low_stock"},
		{"SKU-045", "RFID Tags (500pk)", "Technology", 40, 10, 195.00, 0.8, "in_stock"},
		{"SKU-046", "Loading Dock Seal", "Equipment", 8, 2, 780.00, 35.0, "in_stock"},
		{"SKU-047", "Platform Scale 1000lb", "Equipment", 6, 2, 550.00, 25.0, "in_stock"},
		{"SKU-048", "Drum Dolly Steel", "Equipment", 14, 4, 78.00, 12.0, "in_stock"},
		{"SKU-049", "Hydraulic Lift Table", "Equipment", 5, 1, 1100.00, 120.0, "in_stock"},
		{"SKU-050", "Warehouse Management Tablet", "Technology", 20, 5, 499.99, 0.6, "in_stock"},
	}

	inventorySQL := `INSERT INTO inventory_items (id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (tenant_id, sku) DO NOTHING`

	for i, item := range inventoryItems {
		wh := acmeWarehouses[i%len(acmeWarehouses)]
		if _, err := pool.Exec(ctx, inventorySQL, uuid.New(), acmeTenantID, wh.id, item.sku, item.name, item.category, item.quantity, item.minQty, item.price, item.weight, item.status, now.Add(-90*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed inventory (%s): %w", item.sku, err)
		}
	}

	// ---------------------------------------------------------------
	// 9. ORDERS — Acme (156) + Beta (20)
	// ---------------------------------------------------------------
	// Revenue target for Acme: ~$284,590
	// 156 orders, average ~$1824 each
	orderStatuses := []string{"delivered", "delivered", "delivered", "shipped", "processing", "pending", "cancelled", "returned"}
	orderTypes := []string{"standard", "express", "freight", "standard", "express", "standard"}

	orderSQL := `INSERT INTO orders (id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (tenant_id, order_number) DO NOTHING`

	// Pre-calculated amounts that sum to ~$284,590 across 156 orders.
	// Base amounts with variation.
	acmeOrderAmounts := make([]float64, 156)
	runningTotal := 0.0
	for i := 0; i < 156; i++ {
		// Vary between $50 and $5000
		base := 50.0 + float64((i*1973+500)%4951)
		acmeOrderAmounts[i] = base
		runningTotal += base
	}
	// Scale to hit ~$284,590
	scaleFactor := 284590.0 / runningTotal
	for i := range acmeOrderAmounts {
		acmeOrderAmounts[i] = float64(int(acmeOrderAmounts[i]*scaleFactor*100)) / 100
	}

	for i := 0; i < 156; i++ {
		status := orderStatuses[i%len(orderStatuses)]
		oType := orderTypes[i%len(orderTypes)]
		custFirst := customerFirstNames[i%len(customerFirstNames)]
		custLast := customerLastNames[(i+3)%len(customerLastNames)]
		amount := acmeOrderAmounts[i]

		// Link some orders to shipments
		var shipmentPtr *uuid.UUID
		if i < 42 {
			sid := acmeShipments[i].id
			shipmentPtr = &sid
		}

		schedDate := today.Add(-time.Duration(i%30) * 24 * time.Hour)

		var returnReason *string
		var cancelReason *string
		if status == "returned" {
			r := "Customer requested return - item not as described"
			returnReason = &r
		}
		if status == "cancelled" {
			c := "Customer cancelled before shipment"
			cancelReason = &c
		}

		createdAt := now.Add(-time.Duration(1+i%60) * 24 * time.Hour)

		if _, err := pool.Exec(ctx, orderSQL, uuid.New(), acmeTenantID, fmt.Sprintf("ORD-%d", 1001+i),
			fmt.Sprintf("%s %s", custFirst, custLast),
			fmt.Sprintf("%s.%s%d@example.com", custFirst, custLast, i),
			status, oType, amount, shipmentPtr, schedDate, returnReason, cancelReason,
			createdAt, now); err != nil {
			return fmt.Errorf("seed orders (ORD-%d): %w", 1001+i, err)
		}
	}

	// Beta orders (20)
	for i := 0; i < 20; i++ {
		status := orderStatuses[i%len(orderStatuses)]
		oType := orderTypes[i%len(orderTypes)]
		amount := 100.0 + float64((i*347)%4900)

		var shipmentPtr *uuid.UUID
		if i < 10 {
			sid := betaShipments[i].id
			shipmentPtr = &sid
		}

		var returnReason *string
		var cancelReason *string
		if status == "returned" {
			r := "Damaged in transit"
			returnReason = &r
		}
		if status == "cancelled" {
			c := "Duplicate order"
			cancelReason = &c
		}

		createdAt := now.Add(-time.Duration(1+i%20) * 24 * time.Hour)

		if _, err := pool.Exec(ctx, orderSQL, uuid.New(), betaTenantID, fmt.Sprintf("BORD-%d", 3001+i),
			fmt.Sprintf("%s %s", customerFirstNames[i%5], customerLastNames[i%5+5]),
			fmt.Sprintf("beta.order%d@example.com", i+1),
			status, oType, amount, shipmentPtr, today.Add(-time.Duration(i%15)*24*time.Hour), returnReason, cancelReason,
			createdAt, now); err != nil {
			return fmt.Errorf("seed beta orders (BORD-%d): %w", 3001+i, err)
		}
	}

	// ---------------------------------------------------------------
	// 10. VENDORS — Acme (12)
	// ---------------------------------------------------------------
	type seedVendor struct {
		name     string
		contact  string
		email    string
		phone    string
		address  string
		category string
		rating   float64
		status   string
	}

	acmeVendors := []seedVendor{
		{"Pacific Fuel Supply", "David Park", "david@pacificfuel.com", "(310) 555-1001", "800 Harbor Blvd, Long Beach, CA 90802", "Fuel", 4.5, "active"},
		{"Continental Tire Depot", "Frank Romano", "frank@contitires.com", "(312) 555-1002", "2300 S Cicero Ave, Chicago, IL 60804", "Tires", 4.8, "active"},
		{"Summit Fleet Parts", "Karen Mitchell", "karen@summitparts.com", "(214) 555-1003", "5600 Stemmons Fwy, Dallas, TX 75207", "Parts", 4.2, "active"},
		{"National Insurance Group", "Robert Chen", "robert@natinsure.com", "(212) 555-1004", "350 5th Ave, New York, NY 10118", "Insurance", 4.0, "active"},
		{"TechFleet GPS Solutions", "Priya Sharma", "priya@techfleet.com", "(415) 555-1005", "100 Spear St, San Francisco, CA 94105", "Technology", 4.6, "active"},
		{"SafeHaul Training Co", "Jim Patterson", "jim@safehaul.com", "(404) 555-1006", "1400 Peachtree Rd, Atlanta, GA 30309", "Training", 3.9, "active"},
		{"Midwest Diesel Repair", "Tony Kowalski", "tony@mwdiesel.com", "(312) 555-1007", "7800 S Pulaski Rd, Chicago, IL 60652", "Maintenance", 4.3, "active"},
		{"GreenFleet EV Charging", "Sarah Lin", "sarah@greenfleet.com", "(206) 555-1008", "400 Occidental Ave S, Seattle, WA 98104", "Infrastructure", 4.1, "active"},
		{"QuickWash Fleet Cleaning", "Miguel Ortiz", "miguel@quickwash.com", "(713) 555-1009", "3200 Telephone Rd, Houston, TX 77023", "Cleaning", 3.8, "active"},
		{"LoadStar Freight Brokers", "Nancy Williams", "nancy@loadstar.com", "(615) 555-1010", "900 Broadway, Nashville, TN 37203", "Brokerage", 4.4, "active"},
		{"Atlas Warehouse Supply", "George Patel", "george@atlassupply.com", "(305) 555-1011", "2000 NW 21st St, Miami, FL 33142", "Warehouse Supplies", 4.0, "active"},
		{"Pioneer Logistics Software", "Lisa Yamamoto", "lisa@pioneerlogi.com", "(503) 555-1012", "1500 SW 1st Ave, Portland, OR 97201", "Software", 4.7, "active"},
	}

	vendorSQL := `INSERT INTO vendors (id, tenant_id, name, contact_person, email, phone, address, category, rating, contract_start, contract_end, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT DO NOTHING`

	for i, v := range acmeVendors {
		contractStart := now.Add(-time.Duration(365+i*30) * 24 * time.Hour)
		contractEnd := now.Add(time.Duration(365-i*20) * 24 * time.Hour)
		if _, err := pool.Exec(ctx, vendorSQL, uuid.New(), acmeTenantID, v.name, v.contact, v.email, v.phone, v.address, v.category, v.rating, contractStart, contractEnd, v.status, now.Add(-180*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed vendors (%s): %w", v.name, err)
		}
	}

	// ---------------------------------------------------------------
	// 11. CLIENTS — Acme (20) + Beta (3)
	// ---------------------------------------------------------------
	type seedClient struct {
		id          uuid.UUID
		tenantID    uuid.UUID
		company     string
		contact     string
		email       string
		phone       string
		address     string
		industry    string
		shipments   int
		spent       float64
		satisfaction float64
		status      string
	}

	acmeClients := []seedClient{
		{uuid.New(), acmeTenantID, "TechVista Electronics", "Amanda Rodriguez", "amanda@techvista.com", "(408) 555-2001", "2500 Great America Pkwy, Santa Clara, CA 95054", "Electronics", 145, 52300.00, 4.7, "active"},
		{uuid.New(), acmeTenantID, "GreenLeaf Organics", "Brian Park", "brian@greenleaf.com", "(503) 555-2002", "800 NE Broadway, Portland, OR 97232", "Food & Beverage", 89, 31200.00, 4.5, "active"},
		{uuid.New(), acmeTenantID, "Pinnacle Construction", "Carlos Vega", "carlos@pinnaclecon.com", "(214) 555-2003", "1600 Pacific Ave, Dallas, TX 75201", "Construction", 67, 98500.00, 4.2, "active"},
		{uuid.New(), acmeTenantID, "MedSupply Direct", "Diana Foster", "diana@medsupply.com", "(617) 555-2004", "200 Seaport Blvd, Boston, MA 02210", "Healthcare", 210, 125000.00, 4.9, "active"},
		{uuid.New(), acmeTenantID, "StyleHouse Fashion", "Emily Chen", "emily@stylehouse.com", "(212) 555-2005", "520 Broadway, New York, NY 10012", "Retail", 178, 43800.00, 4.3, "active"},
		{uuid.New(), acmeTenantID, "AutoParts Unlimited", "Frank Marino", "frank@autoparts.com", "(313) 555-2006", "3000 W Grand Blvd, Detroit, MI 48202", "Automotive", 95, 67200.00, 4.1, "active"},
		{uuid.New(), acmeTenantID, "Harvest Moon Farms", "Grace Thompson", "grace@harvestmoon.com", "(559) 555-2007", "1200 Clovis Ave, Fresno, CA 93612", "Agriculture", 42, 18900.00, 4.6, "active"},
		{uuid.New(), acmeTenantID, "BlueWave Marine", "Henry Nakamura", "henry@bluewave.com", "(206) 555-2008", "1500 Westlake Ave N, Seattle, WA 98109", "Marine", 38, 29400.00, 4.4, "active"},
		{uuid.New(), acmeTenantID, "SolarBright Energy", "Irene Walsh", "irene@solarbright.com", "(602) 555-2009", "4400 N Scottsdale Rd, Scottsdale, AZ 85251", "Energy", 56, 41200.00, 4.0, "active"},
		{uuid.New(), acmeTenantID, "PetCare Holdings", "Jack Sullivan", "jack@petcare.com", "(615) 555-2010", "200 4th Ave N, Nashville, TN 37219", "Pet Care", 132, 22100.00, 4.8, "active"},
		{uuid.New(), acmeTenantID, "QuickPrint Services", "Karen Lee", "karen@quickprint.com", "(404) 555-2011", "1100 Spring St NW, Atlanta, GA 30309", "Printing", 64, 15700.00, 4.2, "active"},
		{uuid.New(), acmeTenantID, "Rocky Mountain Brewing", "Larry Hudson", "larry@rockymtn.com", "(720) 555-2012", "2000 Lawrence St, Denver, CO 80205", "Beverage", 88, 35600.00, 4.5, "active"},
		{uuid.New(), acmeTenantID, "Coastal Imports LLC", "Maria Santos", "maria@coastalimports.com", "(305) 555-2013", "800 Brickell Ave, Miami, FL 33131", "Import/Export", 156, 87300.00, 4.3, "active"},
		{uuid.New(), acmeTenantID, "NorthStar Manufacturing", "Nathan Berg", "nathan@northstar.com", "(612) 555-2014", "500 Hennepin Ave, Minneapolis, MN 55403", "Manufacturing", 73, 54800.00, 4.1, "active"},
		{uuid.New(), acmeTenantID, "Pacific Rim Trading", "Olivia Tanaka", "olivia@pacificrim.com", "(310) 555-2015", "400 Continental Blvd, El Segundo, CA 90245", "Trading", 112, 62400.00, 4.6, "active"},
		{uuid.New(), acmeTenantID, "Phoenix Home Goods", "Patrick Rivera", "patrick@phoenixhome.com", "(480) 555-2016", "7000 E Mayo Blvd, Phoenix, AZ 85054", "Home Goods", 94, 28900.00, 4.4, "active"},
		{uuid.New(), acmeTenantID, "Quantum Labs Inc", "Rachel Kim", "rachel@quantumlabs.com", "(858) 555-2017", "9500 Gilman Dr, La Jolla, CA 92093", "Biotech", 31, 45200.00, 4.7, "active"},
		{uuid.New(), acmeTenantID, "Redwood Furniture Co", "Sam Mitchell", "sam@redwoodfurn.com", "(916) 555-2018", "1800 Capitol Ave, Sacramento, CA 95811", "Furniture", 48, 39100.00, 4.0, "active"},
		{uuid.New(), acmeTenantID, "Summit Sports Gear", "Tara O'Brien", "tara@summitsports.com", "(801) 555-2019", "600 S Main St, Salt Lake City, UT 84101", "Sporting Goods", 81, 21600.00, 4.5, "active"},
		{uuid.New(), acmeTenantID, "United Chemical Supply", "Victor Patel", "victor@unitedchem.com", "(713) 555-2020", "1200 Smith St, Houston, TX 77002", "Chemical", 27, 72400.00, 3.9, "active"},
	}

	betaClients := []seedClient{
		{uuid.New(), betaTenantID, "Delta Electronics", "Amy Wang", "amy@deltaelec.com", "(704) 555-3001", "300 S Tryon St, Charlotte, NC 28202", "Electronics", 25, 18500.00, 4.3, "active"},
		{uuid.New(), betaTenantID, "Lone Star Textiles", "Bob Turner", "bob@lonestar.com", "(512) 555-3002", "100 Congress Ave, Austin, TX 78701", "Textiles", 18, 12200.00, 4.1, "active"},
		{uuid.New(), betaTenantID, "Bayou Foods Inc", "Claire Dubois", "claire@bayoufoods.com", "(504) 555-3003", "800 Poydras St, New Orleans, LA 70112", "Food & Beverage", 32, 24800.00, 4.6, "active"},
	}

	clientSQL := `INSERT INTO clients (id, tenant_id, company_name, contact_person, email, phone, address, industry, total_shipments, total_spent, satisfaction_rating, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT DO NOTHING`

	allClients := append(acmeClients, betaClients...)
	for _, c := range allClients {
		if _, err := pool.Exec(ctx, clientSQL, c.id, c.tenantID, c.company, c.contact, c.email, c.phone, c.address, c.industry, c.shipments, c.spent, c.satisfaction, c.status, now.Add(-120*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed clients (%s): %w", c.company, err)
		}
	}

	// ---------------------------------------------------------------
	// 12. CLIENT FEEDBACK — Acme (15)
	// ---------------------------------------------------------------
	type seedFeedback struct {
		clientIdx int
		rating    int
		comment   string
		category  string
	}

	feedbacks := []seedFeedback{
		{0, 5, "Outstanding delivery speed and package condition. TechVista is very satisfied with Acme's services.", "delivery"},
		{1, 4, "Good overall service. Temperature-controlled shipments arrived in perfect condition.", "service"},
		{2, 4, "Reliable heavy equipment transport. Minor delay on one shipment but well communicated.", "delivery"},
		{3, 5, "Critical medical supplies always arrive on time. Excellent handling of fragile items.", "reliability"},
		{4, 4, "Fashion inventory shipments well-handled. Appreciate the careful packaging.", "packaging"},
		{5, 3, "Some parts arrived with minor damage. Packaging could be improved for auto parts.", "packaging"},
		{6, 5, "Refrigerated transport for produce is exceptional. Zero spoilage this quarter.", "service"},
		{7, 5, "Marine equipment shipping handled professionally. Great communication throughout.", "communication"},
		{8, 4, "Solar panel shipments require extra care - Acme delivers consistently.", "delivery"},
		{9, 5, "Pet supplies always arrive on schedule. Our customers are happy, so we are happy.", "reliability"},
		{10, 4, "Print materials delivered without damage. Good tracking system.", "technology"},
		{11, 5, "Brewery equipment moved safely across country. Highly recommend Acme.", "service"},
		{12, 4, "International import coordination is smooth. Customs documentation well managed.", "documentation"},
		{13, 3, "Manufacturing parts delivery could be faster. Current lead times are acceptable but not great.", "delivery"},
		{14, 5, "Pacific Rim is consistently impressed with Acme's cross-country logistics capability.", "reliability"},
	}

	feedbackSQL := `INSERT INTO client_feedback (id, tenant_id, client_id, rating, comment, category, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT DO NOTHING`

	for i, f := range feedbacks {
		createdAt := now.Add(-time.Duration(i*5+1) * 24 * time.Hour)
		if _, err := pool.Exec(ctx, feedbackSQL, uuid.New(), acmeTenantID, acmeClients[f.clientIdx].id, f.rating, f.comment, f.category, createdAt); err != nil {
			return fmt.Errorf("seed feedback (%d): %w", i, err)
		}
	}

	// ---------------------------------------------------------------
	// 13. ROLES — Acme (5)
	// ---------------------------------------------------------------
	type seedRole struct {
		name        string
		permissions string
	}

	roles := []seedRole{
		{"Admin", `{"shipments":["create","read","update","delete"],"vehicles":["create","read","update","delete"],"drivers":["create","read","update","delete"],"warehouses":["create","read","update","delete"],"orders":["create","read","update","delete"],"vendors":["create","read","update","delete"],"clients":["create","read","update","delete"],"users":["create","read","update","delete"],"settings":["create","read","update","delete"],"reports":["read"],"dashboard":["read"]}`},
		{"Manager", `{"shipments":["create","read","update"],"vehicles":["create","read","update"],"drivers":["create","read","update"],"warehouses":["read","update"],"orders":["create","read","update"],"vendors":["read","update"],"clients":["create","read","update"],"users":["read"],"settings":["read"],"reports":["read"],"dashboard":["read"]}`},
		{"Dispatcher", `{"shipments":["create","read","update"],"vehicles":["read"],"drivers":["read","update"],"warehouses":["read"],"orders":["read","update"],"vendors":["read"],"clients":["read"],"reports":["read"],"dashboard":["read"]}`},
		{"Driver", `{"shipments":["read"],"vehicles":["read"],"drivers":["read"],"orders":["read"],"dashboard":["read"]}`},
		{"Viewer", `{"shipments":["read"],"vehicles":["read"],"drivers":["read"],"warehouses":["read"],"orders":["read"],"vendors":["read"],"clients":["read"],"reports":["read"],"dashboard":["read"]}`},
	}

	roleSQL := `INSERT INTO roles (id, tenant_id, name, permissions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tenant_id, name) DO NOTHING`

	for _, r := range roles {
		if _, err := pool.Exec(ctx, roleSQL, uuid.New(), acmeTenantID, r.name, r.permissions, now.Add(-180*24*time.Hour), now); err != nil {
			return fmt.Errorf("seed roles (%s): %w", r.name, err)
		}
	}

	// ---------------------------------------------------------------
	// 14. SETTINGS — Acme (10)
	// ---------------------------------------------------------------
	type seedSetting struct {
		key      string
		value    string
		category string
	}

	settings := []seedSetting{
		{"company_name", "Acme Logistics", "general"},
		{"timezone", "America/Chicago", "general"},
		{"currency", "USD", "general"},
		{"date_format", "MM/DD/YYYY", "general"},
		{"email_notifications", "true", "notifications"},
		{"sms_notifications", "false", "notifications"},
		{"low_fuel_threshold", "20", "fleet"},
		{"maintenance_reminder_days", "30", "fleet"},
		{"default_page_size", "20", "display"},
		{"dashboard_refresh_interval", "300", "display"},
	}

	settingSQL := `INSERT INTO settings (id, tenant_id, key, value, category, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, key) DO NOTHING`

	for _, s := range settings {
		if _, err := pool.Exec(ctx, settingSQL, uuid.New(), acmeTenantID, s.key, s.value, s.category, acmeAdmin.id, now); err != nil {
			return fmt.Errorf("seed settings (%s): %w", s.key, err)
		}
	}

	// ---------------------------------------------------------------
	// 15. ACTIVITY LOG — Acme (30)
	// ---------------------------------------------------------------
	type seedActivity struct {
		userID     uuid.UUID
		action     string
		entityType string
		details    string
		ip         string
	}

	activities := []seedActivity{
		{acmeAdmin.id, "user.login", "user", `{"email":"admin@acme.com"}`, "192.168.1.10"},
		{acmeAdmin.id, "settings.update", "setting", `{"key":"timezone","old":"America/New_York","new":"America/Chicago"}`, "192.168.1.10"},
		{acmeManager.id, "user.login", "user", `{"email":"john@acme.com"}`, "192.168.1.15"},
		{acmeManager.id, "shipment.create", "shipment", `{"tracking":"SHP-1001","origin":"Los Angeles, CA"}`, "192.168.1.15"},
		{acmeManager.id, "shipment.update", "shipment", `{"tracking":"SHP-1002","status":"in_transit"}`, "192.168.1.15"},
		{acmeDispatcher.id, "user.login", "user", `{"email":"sarah@acme.com"}`, "192.168.1.20"},
		{acmeDispatcher.id, "driver.assign", "driver", `{"driver":"DRV-001","vehicle":"TRK-001"}`, "192.168.1.20"},
		{acmeDispatcher.id, "shipment.update", "shipment", `{"tracking":"SHP-1005","status":"delivered"}`, "192.168.1.20"},
		{acmeDispatcher.id, "shipment.create", "shipment", `{"tracking":"SHP-1010","origin":"Chicago, IL"}`, "192.168.1.20"},
		{acmeDriver.id, "user.login", "user", `{"email":"mike@acme.com"}`, "10.0.0.50"},
		{acmeDriver.id, "shipment.update", "shipment", `{"tracking":"SHP-1003","status":"delivered","location":"Dallas, TX"}`, "10.0.0.50"},
		{acmeAdmin.id, "user.create", "user", `{"email":"viewer@acme.com","role":"viewer"}`, "192.168.1.10"},
		{acmeAdmin.id, "role.update", "role", `{"name":"Dispatcher","added_permission":"shipments.create"}`, "192.168.1.10"},
		{acmeManager.id, "order.create", "order", `{"order_number":"ORD-1001","amount":1250.00}`, "192.168.1.15"},
		{acmeManager.id, "client.create", "client", `{"company":"TechVista Electronics"}`, "192.168.1.15"},
		{acmeManager.id, "vendor.update", "vendor", `{"name":"Pacific Fuel Supply","rating":4.5}`, "192.168.1.15"},
		{acmeDispatcher.id, "vehicle.update", "vehicle", `{"vehicle_id":"TRK-005","status":"maintenance"}`, "192.168.1.20"},
		{acmeDispatcher.id, "maintenance.create", "maintenance", `{"vehicle":"TRK-005","type":"Oil Change"}`, "192.168.1.20"},
		{acmeAdmin.id, "warehouse.update", "warehouse", `{"name":"LA Distribution Center","used_capacity":38500}`, "192.168.1.10"},
		{acmeManager.id, "shipment.create", "shipment", `{"tracking":"SHP-1020","destination":"Seattle, WA"}`, "192.168.1.15"},
		{acmeDispatcher.id, "driver.update", "driver", `{"employee_id":"DRV-005","status":"on_trip"}`, "192.168.1.20"},
		{acmeManager.id, "order.update", "order", `{"order_number":"ORD-1010","status":"shipped"}`, "192.168.1.15"},
		{acmeAdmin.id, "settings.update", "setting", `{"key":"low_fuel_threshold","old":"15","new":"20"}`, "192.168.1.10"},
		{acmeDispatcher.id, "shipment.update", "shipment", `{"tracking":"SHP-1025","status":"in_transit"}`, "192.168.1.20"},
		{acmeManager.id, "report.generate", "report", `{"type":"monthly_revenue","month":"January"}`, "192.168.1.15"},
		{acmeAdmin.id, "user.update", "user", `{"email":"john@acme.com","role":"manager"}`, "192.168.1.10"},
		{acmeDispatcher.id, "shipment.create", "shipment", `{"tracking":"SHP-1035","origin":"Miami, FL"}`, "192.168.1.20"},
		{acmeManager.id, "client.update", "client", `{"company":"MedSupply Direct","satisfaction_rating":4.9}`, "192.168.1.15"},
		{acmeDriver.id, "shipment.update", "shipment", `{"tracking":"SHP-1040","status":"delivered","location":"Phoenix, AZ"}`, "10.0.0.50"},
		{acmeViewer.id, "user.login", "user", `{"email":"viewer@acme.com"}`, "192.168.1.30"},
	}

	activitySQL := `INSERT INTO activity_log (id, tenant_id, user_id, action, entity_type, entity_id, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT DO NOTHING`

	for i, a := range activities {
		createdAt := now.Add(-time.Duration(len(activities)-i) * time.Hour)
		if _, err := pool.Exec(ctx, activitySQL, uuid.New(), acmeTenantID, a.userID, a.action, a.entityType, uuid.New(), a.details, a.ip, createdAt); err != nil {
			return fmt.Errorf("seed activity_log (%d): %w", i, err)
		}
	}

	// ---------------------------------------------------------------
	// 16. NOTIFICATIONS — Acme (varied across users)
	// ---------------------------------------------------------------
	type seedNotification struct {
		userID  uuid.UUID
		title   string
		message string
		ntype   string
		read    bool
	}

	notifications := []seedNotification{
		{acmeAdmin.id, "System Update Complete", "CargoMax platform has been updated to version 3.2.1 with new reporting features.", "system", true},
		{acmeAdmin.id, "New User Registration", "A new user viewer@acme.com has been registered and requires role assignment.", "user", true},
		{acmeAdmin.id, "Monthly Report Available", "The January 2026 monthly performance report is ready for review.", "report", false},
		{acmeManager.id, "Shipment Delayed", "SHP-1035 has been delayed due to weather conditions in Denver, CO.", "shipment", false},
		{acmeManager.id, "New Client Added", "TechVista Electronics has been added as a new client.", "client", true},
		{acmeManager.id, "Revenue Milestone", "Monthly revenue has exceeded $280,000 target.", "financial", false},
		{acmeManager.id, "Vehicle Maintenance Due", "TRK-012 is due for scheduled maintenance in 3 days.", "fleet", false},
		{acmeDispatcher.id, "Driver Assignment", "Driver DRV-001 has been assigned to vehicle TRK-001.", "fleet", true},
		{acmeDispatcher.id, "Route Optimization", "New optimized routes available for 5 pending shipments.", "shipment", false},
		{acmeDispatcher.id, "Low Fuel Alert", "Vehicle TRK-023 fuel level is below 20%.", "fleet", false},
		{acmeDriver.id, "New Delivery Assignment", "You have been assigned shipment SHP-1040 for delivery to Phoenix, AZ.", "shipment", true},
		{acmeDriver.id, "Route Update", "Your route for SHP-1040 has been updated to avoid construction on I-10.", "shipment", false},
		{acmeViewer.id, "Welcome to CargoMax", "Welcome to the Acme Logistics dashboard. Contact admin for role upgrades.", "system", false},
	}

	notificationSQL := `INSERT INTO notifications (id, tenant_id, user_id, title, message, type, read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT DO NOTHING`

	for i, n := range notifications {
		createdAt := now.Add(-time.Duration(len(notifications)-i) * 2 * time.Hour)
		if _, err := pool.Exec(ctx, notificationSQL, uuid.New(), acmeTenantID, n.userID, n.title, n.message, n.ntype, n.read, createdAt); err != nil {
			return fmt.Errorf("seed notifications (%d): %w", i, err)
		}
	}

	// ---------------------------------------------------------------
	// 17. NOTIFICATION PREFERENCES — for all Acme users
	// ---------------------------------------------------------------
	notifPrefSQL := `INSERT INTO notification_preferences (id, tenant_id, user_id, event_type, email_enabled, sms_enabled, push_enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, user_id, event_type) DO NOTHING`

	eventTypes := []string{"shipment_update", "order_update", "fleet_alert", "system_notification"}
	acmeUserList := []seedUser{acmeAdmin, acmeManager, acmeDispatcher, acmeDriver, acmeViewer}
	for _, u := range acmeUserList {
		for _, et := range eventTypes {
			emailEnabled := true
			smsEnabled := u.role == "admin" || u.role == "manager"
			pushEnabled := u.role != "viewer"
			if _, err := pool.Exec(ctx, notifPrefSQL, uuid.New(), acmeTenantID, u.id, et, emailEnabled, smsEnabled, pushEnabled); err != nil {
				return fmt.Errorf("seed notification_preferences (%s/%s): %w", u.email, et, err)
			}
		}
	}

	log.Println("Seed data inserted successfully.")
	log.Printf("  Acme Logistics (tenant: %s) — 5 users, 42 shipments, 45 vehicles, 30 drivers, 20 maintenance, 8 warehouses, 50 inventory, 156 orders, 12 vendors, 20 clients, 15 feedback, 5 roles, 10 settings, 30 activity logs, 13 notifications", acmeTenantID)
	log.Printf("  Beta Transport  (tenant: %s) — 2 users, 10 shipments, 5 vehicles, 3 drivers, 2 warehouses, 20 orders, 3 clients", betaTenantID)

	return nil
}
