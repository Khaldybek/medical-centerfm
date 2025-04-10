package gorm

import (
	"errors"
	"gorm.io/gorm"
	"medical-center/internal/models/appointment"
)

type AppoinmentRepository struct {
	db *gorm.DB
}

func NewAppoinmentRepository(db *gorm.DB) *AppoinmentRepository {
	return &AppoinmentRepository{db: db}
}

func (r *AppoinmentRepository) Create(appoint *appointment.Appointment) error {
	return r.db.Create(appoint).Error
}

func (r *AppoinmentRepository) GetByID(id uint) (*appointment.Appointment, error) {
	var appoint appointment.Appointment
	err := r.db.First(&appoint, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("appoint  not found")
	}
	return &appoint, err
}

func (r *AppoinmentRepository) GetAll() ([]appointment.Appointment, error) {
	var appoint []appointment.Appointment
	err := r.db.Find(&appoint).Error
	return appoint, err
}

func (r *AppoinmentRepository) GetByDepartment(departmentID uint) ([]appointment.Appointment, error) {
	var appoint []appointment.Appointment
	err := r.db.Where("department_id = ?", departmentID).Find(&appoint).Error
	return appoint, err
}

func (r *AppoinmentRepository) GetByPatient(name string) ([]appointment.Appointment, error) {
	var appoint []appointment.Appointment
	err := r.db.Where("patient_name = ?", name).Find(&appoint).Error
	return appoint, err
}

func (r *AppoinmentRepository) Update(appoint *appointment.Appointment) error {
	return r.db.Save(appoint).Error
}

func (r *AppoinmentRepository) Delete(id uint) error {
	return r.db.Delete(&appointment.Appointment{}, id).Error
}

func (r *AppoinmentRepository) GetByDoctor(doctorID uint) ([]appointment.Appointment, error) {
	var appoint []appointment.Appointment
	err := r.db.Where("doctor_id = ?", doctorID).Find(&appoint).Error
	return appoint, err
}
