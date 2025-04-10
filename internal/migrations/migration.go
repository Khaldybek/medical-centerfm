package migrations

import (
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration interface {
	ID() string
	Migrate(*gorm.DB) error
	Rollback(*gorm.DB) error
}

// MigrationRecord is the database model to keep track of applied migrations
type MigrationRecord struct {
	ID        string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"not null"`
}

// Migrator handles the migration process
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// AddMigration adds a migration to the migrator
func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// Migrate runs all pending migrations
func (m *Migrator) Migrate() error {
	// Ensure migration table exists
	err := m.db.AutoMigrate(&MigrationRecord{})
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Sort migrations by ID to ensure consistent order
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID() < m.migrations[j].ID()
	})

	// Get applied migrations
	var appliedMigrations []MigrationRecord
	if err := m.db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map for faster lookup
	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.ID] = true
	}

	// Apply pending migrations
	for _, migration := range m.migrations {
		if !appliedMap[migration.ID()] {
			log.Printf("Running migration: %s", migration.ID())
			
			err := m.db.Transaction(func(tx *gorm.DB) error {
				if err := migration.Migrate(tx); err != nil {
					return err
				}
				
				// Record migration as applied
				record := MigrationRecord{
					ID:        migration.ID(),
					AppliedAt: time.Now(),
				}
				
				return tx.Create(&record).Error
			})
			
			if err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.ID(), err)
			}
			
			log.Printf("Migration %s completed successfully", migration.ID())
		}
	}
	
	return nil
}

// RollbackLast rolls back the last applied migration
func (m *Migrator) RollbackLast() error {
	var lastMigration MigrationRecord
	result := m.db.Order("id DESC").First(&lastMigration)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", result.Error)
	}
	
	// Find the migration to rollback
	var migrationToRollback Migration
	for _, migration := range m.migrations {
		if migration.ID() == lastMigration.ID {
			migrationToRollback = migration
			break
		}
	}
	
	if migrationToRollback == nil {
		return fmt.Errorf("migration %s not found in registered migrations", lastMigration.ID)
	}
	
	// Rollback the migration
	log.Printf("Rolling back migration: %s", lastMigration.ID)
	
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := migrationToRollback.Rollback(tx); err != nil {
			return err
		}
		
		// Remove the migration record
		return tx.Delete(&lastMigration).Error
	})
	
	if err != nil {
		return fmt.Errorf("rollback of migration %s failed: %w", lastMigration.ID, err)
	}
	
	log.Printf("Rollback of migration %s completed successfully", lastMigration.ID)
	return nil
} 