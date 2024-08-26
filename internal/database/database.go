package database

import (
	"database/sql"
	"log"

	"github.com/ananikitina/notes-rest/internal/config"
	"github.com/golang-migrate/migrate/v4"
	mg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Parse PostgresURL into DSN
	dsn, err := pq.ParseURL(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to parse PostgresURL: %v", err)
	}

	// Initialize the database connection
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run migrations
	driver, err := mg.WithInstance(DB, &mg.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///root/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to initialize migration instance: %v", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Migrations are up to date, no changes applied.")
		} else {
			log.Fatalf("Failed to run migration: %v", err)
		}
	}

	log.Println("Migrations ran successfully")

	return DB, nil
}
