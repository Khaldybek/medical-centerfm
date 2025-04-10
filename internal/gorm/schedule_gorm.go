package gorm

import (
	"errors"
	"gorm.io/gorm"
	"medical-center/internal/models/schedule"
	"time"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(slot *schedule.Schedule) error {
	if err := slot.IsValid(); err != nil {
		return err
	}
	return r.db.Create(slot).Error
}

func (r *ScheduleRepository) GetByID(id uint) (*schedule.Schedule, error) {
	var slot schedule.Schedule
	err := r.db.First(&slot, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("schedule slot not found")
	}
	return &slot, err
}

func (r *ScheduleRepository) GetByDoctor(doctorID uint) ([]schedule.Schedule, error) {
	var slots []schedule.Schedule
	err := r.db.Where("doctor_id = ?", doctorID).Find(&slots).Error
	return slots, err
}

func (r *ScheduleRepository) GetAvailable(doctorID uint, date time.Time) ([]schedule.Schedule, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var slots []schedule.Schedule
	err := r.db.Where(
		"doctor_id = ? AND booked = ? AND start_time BETWEEN ? AND ?",
		doctorID,
		false,
		start,
		end,
	).Find(&slots).Error

	return slots, err
}

func (r *ScheduleRepository) Update(slot *schedule.Schedule) error {
	if err := slot.IsValid(); err != nil {
		return err
	}
	return r.db.Save(slot).Error
}

func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&schedule.Schedule{}, id).Error
}

func (r *ScheduleRepository) BookSlot(id uint) error {
	return r.db.Model(&schedule.Schedule{}).
		Where("id = ? AND booked = ?", id, false).
		Update("booked", true).Error
}

func (r *ScheduleRepository) CancelBooking(id uint) error {
	return r.db.Model(&schedule.Schedule{}).
		Where("id = ? AND booked = ?", id, true).
		Update("booked", false).Error
}
