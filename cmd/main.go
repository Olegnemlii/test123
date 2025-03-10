package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Olegnemlii/test123/internal/config"
	"github.com/Olegnemlii/test123/internal/service"
	"github.com/Olegnemlii/test123/internal/transport/grpc/handler"
	"github.com/Olegnemlii/test123/internal/transport/grpc/server"
	"github.com/Olegnemlii/test123/pkg/db"

	"github.com/Olegnemlii/test123/internal/repository/postgres"

	_ "github.com/lib/pq"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database Connection
	dbConnection, err := db.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConnection.Close()
	database := dbConnection.GetDB() // Use GetDB to get *sql.DB

	// Run migrations
	if err := runMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Repository
	userRepo := postgres.NewPostgresUserRepository(database)

	// Service
	authService := service.NewUserService(userRepo)

	// gRPC Handler
	authHandler := handler.NewAuthHandler(authService)

	// Start gRPC server
	if err := server.StartGRPCServer(cfg, authHandler); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	// Read the migration file
	migrationFile := "migrations/12312312_users.sql"
	migrationSQL, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
	}

	// Execute the migration
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Migrations ran successfully")
	return nil
}
