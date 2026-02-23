package config

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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

func Load() *Config {
	_ = godotenv.Load()

	privKey, pubKey := loadOrGenerateKeys()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://cargomax:cargomax_secret@localhost:5432/cargomax?sslmode=disable"),
		JWTPrivateKey: privKey,
		JWTPublicKey:  pubKey,
		CookieDomain:  getEnv("COOKIE_DOMAIN", "localhost"),
		CookieSecure:  getEnv("COOKIE_SECURE", "false") == "true",
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnv("SMTP_PORT", "1025"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:3333"),
	}
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
