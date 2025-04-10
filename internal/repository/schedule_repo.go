package repository

import (
	"medical-center/internal/models/schedule"
	"time"
)

type ScheduleRepository interface {
	Create(slot *schedule.Schedule) error
	GetByID(id uint) (*schedule.Schedule, error)
	GetByDoctor(doctorID uint) ([]schedule.Schedule, error)
	GetAvailable(doctorID uint, date time.Time) ([]schedule.Schedule, error)
	Update(slot *schedule.Schedule) error
	Delete(id uint) error
	BookSlot(id uint) error
	CancelBooking(id uint) error
}
