package migrations

import (
	"gorm.io/gorm"
)

type CreateDoctorsTable struct{}

func (m *CreateDoctorsTable) ID() string {
	return "000002_create_doctors"
}

func (m *CreateDoctorsTable) Migrate(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS doctors (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			name VARCHAR(100) NOT NULL,
			department_id INTEGER NOT NULL,
			available BOOLEAN NOT NULL DEFAULT TRUE,
			CONSTRAINT fk_doctors_department FOREIGN KEY (department_id) REFERENCES departments(id)
		);
		CREATE INDEX IF NOT EXISTS idx_doctors_department_id ON doctors(department_id);
	`).Error
}

func (m *CreateDoctorsTable) Rollback(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS doctors").Error
} 