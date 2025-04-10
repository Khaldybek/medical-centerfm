package migrations

import (
	"gorm.io/gorm"
)

type CreateDepartmentsTable struct{}

func (m *CreateDepartmentsTable) ID() string {
	return "000001_create_departments"
}

func (m *CreateDepartmentsTable) Migrate(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS departments (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			name VARCHAR(100) NOT NULL UNIQUE
		)
	`).Error
}

func (m *CreateDepartmentsTable) Rollback(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS departments").Error
} 