package gorm

import (
	"errors"
	"gorm.io/gorm"
	"medical-center/internal/models/doctor"
)

type DoctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) Create(doct *doctor.Doctor) error {
	return r.db.Create(doct).Error
}

func (r *DoctorRepository) GetByID(id uint) (*doctor.Doctor, error) {
	var doct doctor.Doctor
	err := r.db.Preload("Schedule").First(&doct, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("doctor not found")
	}
	return &doct, err
}

func (r *DoctorRepository) GetAll() ([]doctor.Doctor, error) {
	var doctors []doctor.Doctor
	err := r.db.Find(&doctors).Error
	return doctors, err
}

func (r *DoctorRepository) GetByDepartment(departmentID uint) ([]doctor.Doctor, error) {
	var doctors []doctor.Doctor
	err := r.db.Where("department_id = ?", departmentID).Find(&doctors).Error
	return doctors, err
}

func (r *DoctorRepository) Update(doctor *doctor.Doctor) error {
	return r.db.Save(doctor).Error
}

func (r *DoctorRepository) Delete(id uint) error {
	return r.db.Delete(&doctor.Doctor{}, id).Error
}

func (r *DoctorRepository) SetAvailability(id uint, available bool) error {
	return r.db.Model(&doctor.Doctor{}).
		Where("id = ?", id).
		Update("available", available).Error
}

func (r *DoctorRepository) GetAvailable() ([]doctor.Doctor, error) {
	var doctors []doctor.Doctor
	err := r.db.Where("available = ?", true).Find(&doctors).Error
	return doctors, err
}
