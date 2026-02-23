package config

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost       string
	Port          string
	FrontendPort  string
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

func Load() *Config {
	_ = godotenv.Load()

	privKey, pubKey := loadOrGenerateKeys()

	// APP_HOST is the single source of truth for the deployment host.
	// It can be an IP (e.g. "157.230.168.249") or a domain (e.g. "cargomax.io").
	// All other host-dependent values derive from it.
	appHost := getEnv("APP_HOST", "localhost")
	port := getEnv("PORT", "8080")
	frontendPort := getEnv("FRONTEND_PORT", "3000")

	// Cookie domain: for IP addresses leave it empty (browsers handle it);
	// for named domains use the value directly.
	cookieDomain := getEnv("COOKIE_DOMAIN", "")
	if cookieDomain == "" {
		if net.ParseIP(appHost) == nil && appHost != "localhost" {
			// Named domain – set cookie domain so it works across subdomains.
			cookieDomain = appHost
		}
		// For IP addresses or localhost, leave blank – the browser will
		// scope the cookie to the exact origin automatically.
	}

	// Derive frontend URL from APP_HOST if FRONTEND_URL is not explicitly set.
	frontendURL := getEnv("FRONTEND_URL", "")
	if frontendURL == "" {
		scheme := "http"
		if getEnv("COOKIE_SECURE", "false") == "true" {
			scheme = "https"
		}
		frontendURL = fmt.Sprintf("%s://%s:%s", scheme, appHost, frontendPort)
	}

	cfg := &Config{
		AppHost:       appHost,
		Port:          port,
		FrontendPort:  frontendPort,
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://cargomax:cargomax_secret@localhost:5432/cargomax?sslmode=disable"),
		JWTPrivateKey: privKey,
		JWTPublicKey:  pubKey,
		CookieDomain:  cookieDomain,
		CookieSecure:  getEnv("COOKIE_SECURE", "false") == "true",
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnv("SMTP_PORT", "1025"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		FrontendURL:   frontendURL,
	}

	log.Printf("Config: APP_HOST=%s, FrontendURL=%s, CookieDomain=%q, CookieSecure=%v",
		cfg.AppHost, cfg.FrontendURL, cfg.CookieDomain, cfg.CookieSecure)

	return cfg
}

func loadOrGenerateKeys() (ed25519.PrivateKey, ed25519.PublicKey) {
	privKeyB64 := os.Getenv("JWT_PRIVATE_KEY")
	if privKeyB64 != "" {
		privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyB64)
		if err == nil && len(privKeyBytes) == ed25519.PrivateKeySize {
			privKey := ed25519.PrivateKey(privKeyBytes)
			pubKey := privKey.Public().(ed25519.PublicKey)
			return privKey, pubKey
		}
	}

	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate Ed25519 keys: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(privKey)
	log.Printf("Generated new Ed25519 key pair. Set JWT_PRIVATE_KEY=%s in .env to persist", encoded)

	return privKey, pubKey
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
