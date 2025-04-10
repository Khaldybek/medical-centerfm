package migrations

import (
	"gorm.io/gorm"
)

type CreateSchedulesTable struct{}

func (m *CreateSchedulesTable) ID() string {
	return "000003_create_schedules"
}

func (m *CreateSchedulesTable) Migrate(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS schedules (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			doctor_id INTEGER NOT NULL,
			start_time TIMESTAMP WITH TIME ZONE NOT NULL,
			end_time TIMESTAMP WITH TIME ZONE NOT NULL,
			booked BOOLEAN NOT NULL DEFAULT FALSE,
			CONSTRAINT fk_schedules_doctor FOREIGN KEY (doctor_id) REFERENCES doctors(id)
		);
		CREATE INDEX IF NOT EXISTS idx_schedules_doctor_id ON schedules(doctor_id);
	`).Error
}

func (m *CreateSchedulesTable) Rollback(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS schedules").Error
} 