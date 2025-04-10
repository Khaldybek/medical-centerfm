package repository

import (
	"medical-center/internal/models/department"
	"time"
)

type DepartmentRepository interface {
	Create(depart *department.Department) error
	GetByID(id uint) (*department.Department, error)
	GetAll() ([]department.Department, error)
	Update(depart *department.Department) error
	Delete(id uint) error
	GetWithDoctors(id uint) (*department.Department, error)
	GetAvailableSlots(id uint, date time.Time) ([]time.Time, error)
}
