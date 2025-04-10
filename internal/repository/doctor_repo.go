package repository

import (
	"medical-center/internal/models/doctor"
)

type DoctorRepository interface {
	Create(doctor *doctor.Doctor) error
	GetByID(id uint) (*doctor.Doctor, error)
	GetAll() ([]doctor.Doctor, error)
	GetByDepartment(departmentID uint) ([]doctor.Doctor, error)
	Update(doctor *doctor.Doctor) error
	Delete(id uint) error
	SetAvailability(id uint, available bool) error
	GetAvailable() ([]doctor.Doctor, error)
}
