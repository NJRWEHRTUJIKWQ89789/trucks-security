package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cargomax-api/internal/auth"
	"cargomax-api/internal/config"
	"cargomax-api/internal/database"
	"cargomax-api/internal/graph"
	"cargomax-api/internal/graph/resolvers"
	"cargomax-api/internal/middleware"
	"cargomax-api/internal/models"
	"cargomax-api/internal/repository"
	"cargomax-api/internal/rest"
	"cargomax-api/internal/seed"
	"cargomax-api/internal/workers"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/graphql-go/handler"
)

// The ResponseWriter context key is defined in models.CtxResponseWriter.
// All resolvers use that key to retrieve the writer for setting cookies.

func main() {
	// Load configuration from environment / .env file.
	cfg := config.Load()

	// Connect to PostgreSQL.
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(pool)
	log.Println("Connected to PostgreSQL")

	// Run database migrations.
	if err := database.RunMigrations(pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations complete")

	// Check for --seed flag: populate database with sample data and exit.
	for _, arg := range os.Args[1:] {
		if arg == "--seed" {
			log.Println("Seeding database...")
			if err := seed.SeedData(pool); err != nil {
				log.Fatalf("Seed failed: %v", err)
			}
			log.Println("Seed complete")
			return
		}
	}

	// Create all 18 repositories.
	userRepo := repository.NewUserRepo(pool)
	tenantRepo := repository.NewTenantRepo(pool)
	shipmentRepo := repository.NewShipmentRepo(pool)
	vehicleRepo := repository.NewVehicleRepo(pool)
	driverRepo := repository.NewDriverRepo(pool)
	maintenanceRepo := repository.NewMaintenanceRepo(pool)
	warehouseRepo := repository.NewWarehouseRepo(pool)
	inventoryRepo := repository.NewInventoryRepo(pool)
	orderRepo := repository.NewOrderRepo(pool)
	vendorRepo := repository.NewVendorRepo(pool)
	clientRepo := repository.NewClientRepo(pool)
	feedbackRepo := repository.NewFeedbackRepo(pool)
	dashboardRepo := repository.NewDashboardRepo(pool)
	reportRepo := repository.NewReportRepo(pool)
	notificationRepo := repository.NewNotificationRepo(pool)
	settingRepo := repository.NewSettingRepo(pool)
	roleRepo := repository.NewRoleRepo(pool)
	activityRepo := repository.NewActivityRepo(pool)

	// Create tracking repositories.
	shiftRepo := repository.NewShiftRepo(pool)
	pingRepo := repository.NewGPSPingRepo(pool)
	alertRepo := repository.NewAlertRepo(pool)
	zoneRepo := repository.NewZoneRepo(pool)

	// Create WebSocket hub and tracking handler.
	wsHub := rest.NewHub(cfg)
	go wsHub.Run()

	trackingHandler := rest.NewTrackingHandler(cfg, driverRepo, vehicleRepo, shiftRepo, pingRepo, alertRepo, zoneRepo, wsHub)
	managerHandler := rest.NewManagerHandler(cfg, driverRepo, vehicleRepo, shiftRepo, pingRepo, alertRepo, zoneRepo)

	// Build the unified resolver that every GraphQL field delegates to.
	resolver := &resolvers.Resolver{
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

	// Assemble the GraphQL schema from all domain resolvers.
	schema, err := graph.NewSchema(resolver)
	if err != nil {
		log.Fatalf("Failed to build GraphQL schema: %v", err)
	}
	log.Println("GraphQL schema built")

	// Create the graphql-go HTTP handler with GraphiQL enabled for development.
	gqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Build the Chi router.
	r := chi.NewRouter()

	// Global middleware.
	r.Use(chiMiddleware.Recoverer) // Prevent panics from crashing the server (e.g. unauthenticated GraphQL requests).
	r.Use(middleware.LoggingMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			allowed := []string{
				cfg.FrontendURL,
				"http://localhost:3000",
				"http://localhost:" + cfg.FrontendPort,
				"http://" + cfg.AppHost + ":" + cfg.FrontendPort,
			}
			for _, a := range allowed {
				if origin == a {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint.
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Optional auth middleware: tries to validate the JWT cookie and populate
	// context with user claims. If the cookie is absent or invalid the request
	// continues without authentication -- individual resolvers that require
	// auth check for the presence of models.CtxUserID in the context.
	optionalAuth := optionalAuthMiddleware(cfg)

	// GraphQL endpoint wrapped with optional auth and ResponseWriter injection.
	r.Route("/graphql", func(sub chi.Router) {
		sub.Use(optionalAuth)
		sub.Handle("/*", injectResponseWriter(gqlHandler))
		sub.Handle("/", injectResponseWriter(gqlHandler))
	})

	// REST API routes for driver mobile app.
	r.Mount("/api/v1", trackingHandler.Routes())

	// REST API routes for manager dashboard.
	r.Mount("/api/v1/manager", managerHandler.Routes())

	// WebSocket endpoints for live dashboard.
	r.Get("/ws/tracking/live", wsHub.HandleTrackingWS)
	r.Get("/ws/alerts", wsHub.HandleAlertsWS)

	// Build the HTTP server.
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine so we can listen for shutdown signals.
	go func() {
		log.Printf("CargoMax API listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Start alert detection worker.
	alertWorker := workers.NewAlertWorker(shiftRepo, pingRepo, alertRepo, zoneRepo, wsHub)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go alertWorker.Start(workerCtx)

	// Graceful shutdown on SIGINT or SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}

// optionalAuthMiddleware returns middleware that attempts JWT validation from
// the access-token cookie. On success the user's claims are injected into the
// request context. On failure (no cookie, expired token, etc.) the request
// continues without authentication -- resolvers decide whether to reject.
func optionalAuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := auth.GetAccessToken(r)
			if tokenString != "" {
				claims, err := auth.ValidateToken(cfg.JWTPublicKey, tokenString)
				if err == nil {
					ctx := r.Context()
					ctx = context.WithValue(ctx, models.CtxTenantID, claims.TenantID)
					ctx = context.WithValue(ctx, models.CtxUserID, claims.UserID)
					ctx = context.WithValue(ctx, models.CtxUserRole, claims.Role)
					ctx = context.WithValue(ctx, models.CtxUserEmail, claims.Email)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// injectResponseWriter wraps an HTTP handler so that the http.ResponseWriter
// is available inside the GraphQL resolve context. Auth resolvers use this to
// set HttpOnly cookies after login/register.
func injectResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), models.CtxResponseWriter, w)
		ctx = context.WithValue(ctx, models.CtxHTTPRequest, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
