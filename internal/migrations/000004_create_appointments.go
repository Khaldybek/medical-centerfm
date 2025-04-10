package migrations

import (
	"gorm.io/gorm"
)

type CreateAppointmentsTable struct{}

func (m *CreateAppointmentsTable) ID() string {
	return "000004_create_appointments"
}

func (m *CreateAppointmentsTable) Migrate(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS appointments (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			patient_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			department_id INTEGER NOT NULL,
			doctor_id INTEGER NOT NULL,
			appointment_time TIMESTAMP WITH TIME ZONE NOT NULL,
			CONSTRAINT fk_appointments_department FOREIGN KEY (department_id) REFERENCES departments(id),
			CONSTRAINT fk_appointments_doctor FOREIGN KEY (doctor_id) REFERENCES doctors(id)
		);
		CREATE INDEX IF NOT EXISTS idx_appointments_department_id ON appointments(department_id);
		CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);
	`).Error
}

func (m *CreateAppointmentsTable) Rollback(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS appointments").Error
} 