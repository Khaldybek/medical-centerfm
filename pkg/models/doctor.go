package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	Name           string         `json:"name" binding:"required"`
	Email          string         `json:"email" binding:"required,email" gorm:"unique"`
	DepartmentID   uint           `json:"department_id" binding:"required"`
	Specialization string         `json:"specialization"`
	Availability   []Availability `json:"availability"`
}

type Availability struct {
	gorm.Model
	DoctorID  uint   `json:"doctor_id"`
	Day       string `json:"day" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}
