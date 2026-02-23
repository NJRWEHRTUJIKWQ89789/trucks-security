# CargoMax API - Enterprise Backend Specification

## Architecture
- **Language**: Go 1.22
- **Module**: `cargomax-api`
- **API**: GraphQL (github.com/graphql-go/graphql + github.com/graphql-go/handler)
- **Database**: PostgreSQL 16 via github.com/jackc/pgx/v5/pgxpool
- **Auth**: EdDSA (Ed25519) JWT in HttpOnly cookies
- **Multi-tenant**: Every table has tenant_id, every query filters by it
- **IDs**: UUID everywhere (github.com/google/uuid)
- **Router**: github.com/go-chi/chi/v5 (just for HTTP + middleware, GraphQL is single endpoint)
- **Password**: golang.org/x/crypto/bcrypt

## Project Structure
```
cargomax-api/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go          (EXISTS - loads .env, Ed25519 keys)
│   ├── database/
│   │   ├── postgres.go           (connection pool)
│   │   └── migrations.go         (CREATE TABLE statements)
│   ├── auth/
│   │   ├── jwt.go                (create/validate EdDSA JWT)
│   │   └── cookies.go            (set/clear/get auth cookies)
│   ├── middleware/
│   │   ├── auth.go               (validate JWT from cookie, set context)
│   │   ├── tenant.go             (extract tenant_id, set context)
│   │   └── logging.go            (request logging)
│   ├── models/                   (Go structs with JSON tags)
│   │   ├── base.go, user.go, tenant.go, shipment.go, vehicle.go,
│   │   ├── driver.go, maintenance.go, warehouse.go, inventory.go,
│   │   ├── order.go, vendor.go, client.go, feedback.go,
│   │   ├── notification.go, role.go, setting.go, activity.go
│   ├── repository/               (database access, ALL queries filter by tenant_id)
│   │   ├── user_repo.go, tenant_repo.go, shipment_repo.go,
│   │   ├── vehicle_repo.go, driver_repo.go, maintenance_repo.go,
│   │   ├── warehouse_repo.go, inventory_repo.go, order_repo.go,
│   │   ├── vendor_repo.go, client_repo.go, feedback_repo.go,
│   │   ├── dashboard_repo.go, report_repo.go, notification_repo.go,
│   │   ├── setting_repo.go, role_repo.go, activity_repo.go
│   ├── graph/
│   │   ├── types/                (GraphQL type definitions)
│   │   │   ├── common.go, auth.go, dashboard.go, shipment.go,
│   │   │   ├── fleet.go, warehouse.go, order.go, vendor.go,
│   │   │   ├── client.go, report.go, settings.go
│   │   ├── resolvers/            (GraphQL resolver implementations)
│   │   │   ├── resolver.go (base struct), auth.go, dashboard.go,
│   │   │   ├── shipments.go, fleet.go, warehouses.go, orders.go,
│   │   │   ├── vendors.go, clients.go, reports.go, settings.go
│   │   └── schema.go            (assembles root Query + Mutation)
│   └── seed/seed.go             (realistic data generator)
```

## Database Schema

### tenants
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    plan VARCHAR(50) DEFAULT 'starter',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### users
```sql
CREATE TABLE users (
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
);
```

### shipments
```sql
CREATE TABLE shipments (
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
);
```

### vehicles
```sql
CREATE TABLE vehicles (
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
);
```

### drivers
```sql
CREATE TABLE drivers (
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
);
```

### maintenance_records
```sql
CREATE TABLE maintenance_records (
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
);
```

### warehouses
```sql
CREATE TABLE warehouses (
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
);
```

### inventory_items
```sql
CREATE TABLE inventory_items (
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
);
```

### orders
```sql
CREATE TABLE orders (
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
);
```

### vendors
```sql
CREATE TABLE vendors (
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
);
```

### clients
```sql
CREATE TABLE clients (
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
);
```

### client_feedback
```sql
CREATE TABLE client_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES clients(id),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    category VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### notifications
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255),
    message TEXT,
    type VARCHAR(50),
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### notification_preferences
```sql
CREATE TABLE notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    event_type VARCHAR(100),
    email_enabled BOOLEAN DEFAULT TRUE,
    sms_enabled BOOLEAN DEFAULT FALSE,
    push_enabled BOOLEAN DEFAULT FALSE,
    UNIQUE(tenant_id, user_id, event_type)
);
```

### roles
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    permissions JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);
```

### settings
```sql
CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    category VARCHAR(100),
    updated_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);
```

### activity_log
```sql
CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    action VARCHAR(255),
    entity_type VARCHAR(100),
    entity_id UUID,
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Config Package (EXISTS at internal/config/config.go)
```go
type Config struct {
    Port          string
    DatabaseURL   string
    JWTPrivateKey ed25519.PrivateKey
    JWTPublicKey  ed25519.PublicKey
    CookieDomain  string
    CookieSecure  bool
    SMTPHost      string
    SMTPPort      string
    SMTPUser      string
    SMTPPass      string
    FrontendURL   string
}
func Load() *Config
```

## Context Keys (define in internal/models/base.go)
```go
type ContextKey string
const (
    CtxTenantID ContextKey = "tenant_id"
    CtxUserID   ContextKey = "user_id"
    CtxUserRole ContextKey = "user_role"
    CtxUserEmail ContextKey = "user_email"
)
```
Usage: `tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)`
In resolvers: `tenantID := p.Context.Value(models.CtxTenantID).(uuid.UUID)`

## Model Pattern
```go
package models
type Shipment struct {
    ID                uuid.UUID  `json:"id"`
    TenantID          uuid.UUID  `json:"tenant_id"`
    TrackingNumber    string     `json:"tracking_number"`
    // ... fields matching DB columns
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}
```

## Repository Pattern
```go
package repository
type ShipmentRepo struct { db *pgxpool.Pool }
func NewShipmentRepo(db *pgxpool.Pool) *ShipmentRepo { return &ShipmentRepo{db: db} }

// CRITICAL: EVERY query MUST include tenant_id in WHERE
func (r *ShipmentRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Shipment, error) {
    row := r.db.QueryRow(ctx,
        "SELECT id, tenant_id, tracking_number, ... FROM shipments WHERE id = $1 AND tenant_id = $2",
        id, tenantID)
    // scan and return
}
func (r *ShipmentRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Shipment, int, error)
func (r *ShipmentRepo) Create(ctx context.Context, s *models.Shipment) error
func (r *ShipmentRepo) Update(ctx context.Context, tenantID, id uuid.UUID, s *models.Shipment) error
func (r *ShipmentRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error
```

## GraphQL Type Pattern (internal/graph/types/)
```go
package types

import "github.com/graphql-go/graphql"

var ShipmentType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Shipment",
    Fields: graphql.Fields{
        "id":              &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
        "trackingNumber":  &graphql.Field{Type: graphql.String},
        "origin":          &graphql.Field{Type: graphql.String},
        "destination":     &graphql.Field{Type: graphql.String},
        "status":          &graphql.Field{Type: graphql.String},
        // ... all fields
    },
})

var ShipmentInputType = graphql.NewInputObject(graphql.InputObjectConfig{
    Name: "ShipmentInput",
    Fields: graphql.InputObjectConfigFieldMap{
        "trackingNumber": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
        "origin":         &graphql.InputObjectFieldConfig{Type: graphql.String},
        // ... input fields (no id, tenant_id, timestamps)
    },
})

// Connection type for pagination
var ShipmentConnectionType = graphql.NewObject(graphql.ObjectConfig{
    Name: "ShipmentConnection",
    Fields: graphql.Fields{
        "items":      &graphql.Field{Type: graphql.NewList(ShipmentType)},
        "totalCount": &graphql.Field{Type: graphql.Int},
        "page":       &graphql.Field{Type: graphql.Int},
        "perPage":    &graphql.Field{Type: graphql.Int},
        "totalPages": &graphql.Field{Type: graphql.Int},
    },
})
```

## GraphQL Resolver Pattern (internal/graph/resolvers/)
```go
package resolvers

// resolver.go - base struct (shared by all resolver files)
type Resolver struct {
    UserRepo      *repository.UserRepo
    TenantRepo    *repository.TenantRepo
    ShipmentRepo  *repository.ShipmentRepo
    VehicleRepo   *repository.VehicleRepo
    DriverRepo    *repository.DriverRepo
    MaintenanceRepo *repository.MaintenanceRepo
    WarehouseRepo *repository.WarehouseRepo
    InventoryRepo *repository.InventoryRepo
    OrderRepo     *repository.OrderRepo
    VendorRepo    *repository.VendorRepo
    ClientRepo    *repository.ClientRepo
    FeedbackRepo  *repository.FeedbackRepo
    DashboardRepo *repository.DashboardRepo
    ReportRepo    *repository.ReportRepo
    NotificationRepo *repository.NotificationRepo
    SettingRepo   *repository.SettingRepo
    RoleRepo      *repository.RoleRepo
    ActivityRepo  *repository.ActivityRepo
    Config        *config.Config
}

// In resolver files, create functions that return graphql.Fields:
// shipments.go
func (r *Resolver) ShipmentQueries() graphql.Fields {
    return graphql.Fields{
        "shipments": &graphql.Field{
            Type: types.ShipmentConnectionType,
            Args: graphql.FieldConfigArgument{
                "page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
                "perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
                "status":  &graphql.ArgumentConfig{Type: graphql.String},
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                tenantID := p.Context.Value(models.CtxTenantID).(uuid.UUID)
                page := p.Args["page"].(int)
                perPage := p.Args["perPage"].(int)
                items, total, err := r.ShipmentRepo.List(p.Context, tenantID, page, perPage)
                if err != nil { return nil, err }
                return map[string]interface{}{
                    "items": items, "totalCount": total,
                    "page": page, "perPage": perPage,
                    "totalPages": (total + perPage - 1) / perPage,
                }, nil
            },
        },
        // ... more query fields
    }
}

func (r *Resolver) ShipmentMutations() graphql.Fields {
    return graphql.Fields{
        "createShipment": &graphql.Field{ ... },
        "updateShipment": &graphql.Field{ ... },
        "deleteShipment": &graphql.Field{ ... },
    }
}
```

## Schema Assembly (internal/graph/schema.go)
```go
package graph

func NewSchema(r *resolvers.Resolver) (graphql.Schema, error) {
    queryFields := graphql.Fields{}
    mutationFields := graphql.Fields{}

    // Merge all resolver query/mutation fields
    for k, v := range r.AuthQueries() { queryFields[k] = v }
    for k, v := range r.DashboardQueries() { queryFields[k] = v }
    for k, v := range r.ShipmentQueries() { queryFields[k] = v }
    // ... all domains
    for k, v := range r.AuthMutations() { mutationFields[k] = v }
    // ... all domains

    return graphql.NewSchema(graphql.SchemaConfig{
        Query: graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: queryFields}),
        Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields}),
    })
}
```

## Auth Details
- JWT signed with Ed25519 (EdDSA) using golang-jwt/jwt/v5
- Token stored in HttpOnly, Secure, SameSite=Strict cookie named "cargomax_token"
- Token claims: user_id (UUID), tenant_id (UUID), email (string), role (string), exp, iat
- Access token: 15 min expiry
- Refresh token: 7 day expiry, stored in "cargomax_refresh" cookie
- Email verification: random token stored in user record, verified via mutation
- Password: bcrypt hashed

## CRITICAL Security Rules
1. EVERY database query MUST include `AND tenant_id = $N` in WHERE clause
2. NEVER trust tenant_id from request body - ALWAYS get from JWT/context
3. All IDs are UUIDs, validated before use
4. No sequential/guessable IDs anywhere
5. Cookie: HttpOnly=true, Secure=cfg.CookieSecure, SameSite=Strict, Path=/
6. Rate limiting on auth endpoints
7. Input validation on all mutations
8. SQL injection prevention via parameterized queries only ($1, $2, etc.)
