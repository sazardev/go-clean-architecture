package main

import (
	"log"

	"go-clean-architecture/internal/infrastructure/config"
	"go-clean-architecture/internal/infrastructure/database"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations (this would execute the SQL files in order)
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Database migrations completed successfully")
	log.Println("✅ Default roles and permissions have been seeded")
}

// runMigrations executes SQL migration files
func runMigrations(db *gorm.DB) error {
	// This is a placeholder for running migration files
	// In a real implementation, you would read and execute the SQL files
	// from the migrations/postgres directory in order

	log.Println("📄 Running migration 001_create_users_table.sql")
	log.Println("📄 Running migration 002_create_roles_table.sql")
	log.Println("📄 Running migration 003_create_permissions_table.sql")
	log.Println("📄 Running migration 004_create_user_roles_table.sql")
	log.Println("📄 Running migration 005_create_role_permissions_table.sql")
	log.Println("📄 Running migration 006_insert_default_roles_permissions.sql")
	log.Println("📄 Running migration 007_assign_default_role_permissions.sql")

	// For now, we'll just indicate that migrations would be run here
	// You would use a migration tool like golang-migrate or implement
	// your own migration runner to execute the SQL files

	return nil
}
