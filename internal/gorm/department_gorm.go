package gorm

import (
	"errors"
	"gorm.io/gorm"
	"medical-center/internal/models/department"
	"medical-center/internal/models/schedule"
	"time"
)

type DepartmentRepositoryImpl struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepositoryImpl {
	return &DepartmentRepositoryImpl{db: db}
}

func (r *DepartmentRepositoryImpl) Create(depart *department.Department) error {
	return r.db.Create(depart).Error
}

func (r *DepartmentRepositoryImpl) GetByID(id uint) (*department.Department, error) {
	var depart department.Department
	err := r.db.First(&depart, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("department not found")
	}
	return &depart, err
}

func (r *DepartmentRepositoryImpl) GetAll() ([]department.Department, error) {
	var depart []department.Department
	err := r.db.Find(&depart).Error
	return depart, err
}

func (r *DepartmentRepositoryImpl) Update(depart *department.Department) error {
	return r.db.Save(depart).Error
}

func (r *DepartmentRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&department.Department{}, id).Error
}

func (r *DepartmentRepositoryImpl) GetWithDoctors(id uint) (*department.Department, error) {
	var depart department.Department
	err := r.db.Preload("Doctors").First(&depart, id).Error
	if err != nil {
		return nil, err
	}
	return &depart, nil
}

func (r *DepartmentRepositoryImpl) GetAvailableSlots(id uint, date time.Time) ([]time.Time, error) {
	var slots []time.Time
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	// Пример логики: собираем все свободные слоты врачей отделения
	err := r.db.Model(&schedule.Schedule{}).
		Joins("JOIN doctors ON doctors.id = schedules.doctor_id").
		Where("doctors.department_id = ? AND schedules.booked = ? AND schedules.start_time BETWEEN ? AND ?",
			id,
			false,
			start,
			end).
		Pluck("DISTINCT schedules.start_time", &slots).
		Error

	return slots, err
}
