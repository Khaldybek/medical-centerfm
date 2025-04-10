package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	PatientName     string    `json:"patient_name" binding:"required"`
	PatientPhone    string    `json:"patient_phone" binding:"required"`
	DoctorID        uint      `json:"doctor_id" binding:"required"`
	DepartmentID    uint      `json:"department_id" binding:"required"`
	AppointmentTime time.Time `json:"appointment_time" binding:"required"`
	Status          string    `json:"status" gorm:"default:'scheduled'"`
}
