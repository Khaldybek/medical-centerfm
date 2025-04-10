package department

import (
	"gorm.io/gorm"
	"medical-center/internal/models/appointment"
	"medical-center/internal/models/doctor"
)

type Department struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Doctors      []doctor.Doctor
	Appointments []appointment.Appointment
}
