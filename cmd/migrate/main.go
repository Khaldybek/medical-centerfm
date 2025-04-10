package main

import (
	"flag"
	"fmt"
	"log"
	"medical-center/internal/migrations"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Define command line flags
	dsn := flag.String("dsn", "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable", "Database connection string")
	rollback := flag.Bool("rollback", false, "Rollback the last migration")
	flag.Parse()

	// Connect to the database
	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create migrator
	migrator := migrations.NewMigrator(db)

	// Add migrations
	migrator.AddMigration(&migrations.CreateDepartmentsTable{})
	migrator.AddMigration(&migrations.CreateDoctorsTable{})
	migrator.AddMigration(&migrations.CreateSchedulesTable{})
	migrator.AddMigration(&migrations.CreateAppointmentsTable{})
	migrator.AddMigration(&migrations.CreateUsersTable{})

	// Run migrations or rollback
	if *rollback {
		fmt.Println("Rolling back the last migration...")
		if err := migrator.RollbackLast(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")
	} else {
		fmt.Println("Running migrations...")
		if err := migrator.Migrate(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	}

	os.Exit(0)
} 