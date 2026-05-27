package main

import (
	"agro-data-management-system/internal/config"
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

func main() {
	var email, password, role, configPath string
	flag.StringVar(&email, "email", "", "User email (required)")
	flag.StringVar(&password, "password", "", "User password (required)")
	flag.StringVar(&role, "role", "user", "User role (default: user)")
	flag.StringVar(&configPath, "config", "config/local.yaml", "Path to config file")
	flag.Parse()

	if email == "" || password == "" {
		fmt.Println("Usage: createuser --email user@example.com --password pass123 [--role user|admin] [--config path]")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	email = strings.ToLower(strings.TrimSpace(email))
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		role = "user"
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	query := `INSERT INTO users (email, password_hash, role, created_at) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, user.Email, user.PasswordHash, user.Role, user.CreatedAt)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	fmt.Printf("User %s created successfully with role %s\n", user.Email, user.Role)
}
