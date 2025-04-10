package appointment

import (
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	gorm.Model
	PatientName     string    `gorm:"not null"`
	Email           string    `gorm:"size:255;not null"`
	Phone           string    `gorm:"size:20;not null"`
	DepartmentID    uint      `gorm:"index;not null"`
	DoctorID        uint      `gorm:"index;not null"`
	AppointmentTime time.Time `gorm:"not null"`
}
