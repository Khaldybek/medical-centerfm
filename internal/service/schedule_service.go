package service

import (
	"errors"
	"medical-center/internal/models/schedule"
	"medical-center/internal/repository"
	"time"
)

type ScheduleService struct {
	repo repository.ScheduleRepository
}

func NewScheduleService(repo repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) CreateSlot(doctorID uint, start, end time.Time) (*schedule.Schedule, error) {
	if start.After(end) {
		return nil, errors.New("start time cannot be after end time")
	}

	newSlot := &schedule.Schedule{
		DoctorID:  doctorID,
		StartTime: start,
		EndTime:   end,
		Booked:    false,
	}

	if err := s.repo.Create(newSlot); err != nil {
		return nil, err
	}
	return newSlot, nil
}

func (s *ScheduleService) GetSlotByID(id uint) (*schedule.Schedule, error) {
	return s.repo.GetByID(id)
}

func (s *ScheduleService) GetDoctorSlots(doctorID uint) ([]schedule.Schedule, error) {
	return s.repo.GetByDoctor(doctorID)
}

func (s *ScheduleService) UpdateSlot(id uint, start, end time.Time) (*schedule.Schedule, error) {
	slot, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	slot.StartTime = start
	slot.EndTime = end

	if err := s.repo.Update(slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *ScheduleService) DeleteSlot(id uint) error {
	return s.repo.Delete(id)
}

func (s *ScheduleService) BookSlot(id uint) error {
	return s.repo.BookSlot(id)
}

func (s *ScheduleService) CancelBooking(id uint) error {
	return s.repo.CancelBooking(id)
}

func (s *ScheduleService) GetAvailableSlots(doctorID uint, date time.Time) ([]schedule.Schedule, error) {
	return s.repo.GetAvailable(doctorID, date)
}
