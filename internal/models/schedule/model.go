package schedule

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	gorm.Model
	DoctorID  uint      `gorm:"index;not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Booked    bool      `gorm:"default:false"`
}

func (s *Schedule) IsValid() error {
	if s.StartTime.After(s.EndTime) {
		return errors.New("invalid time slot: start time after end time")
	}
	return nil
}
