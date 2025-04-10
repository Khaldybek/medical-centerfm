package department

import (
	"medical-center/internal/models/appointment"
	"medical-center/internal/models/doctor"
)

type DepartmentDTO struct {
	Name         string `gorm:"unique"`
	Doctors      []doctor.Doctor
	Appointments []appointment.Appointment
}
