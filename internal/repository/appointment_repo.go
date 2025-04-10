package repository

import (
	"medical-center/internal/models/appointment"
)

type AppRepository interface {
	Create(app *appointment.Appointment) error
	GetByID(id uint) (*appointment.Appointment, error)
	GetAll() ([]appointment.Appointment, error)
	GetByDepartment(departmentID uint) ([]appointment.Appointment, error)
	GetByPatient(name string) ([]appointment.Appointment, error)
	Update(appointment *appointment.Appointment) error
	Delete(id uint) error
	GetByDoctor(doctorID uint) ([]appointment.Appointment, error)
}
