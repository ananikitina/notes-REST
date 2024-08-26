package database

import (
	"database/sql"
	"log"

	"github.com/ananikitina/notes-rest/internal/config"
	"github.com/golang-migrate/migrate/v4"
	mg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	// Hash existing passwords in the database
	if err := hashExistingPasswords(DB); err != nil {
		log.Fatalf("Failed to hash existing passwords: %v", err)
	}

	log.Println("Existing passwords hashed successfully")

	return DB, nil
}

// hashExistingPasswords hashes plain-text passwords stored in the database.
func hashExistingPasswords(DB *sql.DB) error {
	rows, err := DB.Query("SELECT id, password FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var password string

		if err := rows.Scan(&id, &password); err != nil {
			return err
		}

		// Check if the password is already hashed (bcrypt hashes start with "$2y$")
		if len(password) > 0 && password[:4] != "$2y$" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			// Update the user's password in the database
			_, err = DB.Exec("UPDATE users SET password = $1 WHERE id = $2", string(hashedPassword), id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
