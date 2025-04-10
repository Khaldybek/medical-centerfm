package service

import (
	"errors"
	"medical-center/internal/models/doctor"
	"medical-center/internal/repository"
)

type DoctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(repo repository.DoctorRepository) *DoctorService {
	return &DoctorService{repo: repo}
}

func (s *DoctorService) CreateDoctor(name string, departmentID uint) (*doctor.Doctor, error) {
	if name == "" {
		return nil, errors.New("doctor name cannot be empty")
	}

	newDoctor := &doctor.Doctor{
		Name:         name,
		DepartmentID: departmentID,
		Available:    true,
	}

	if err := s.repo.Create(newDoctor); err != nil {
		return nil, err
	}
	return newDoctor, nil
}

func (s *DoctorService) GetDoctorByID(id uint) (*doctor.Doctor, error) {
	return s.repo.GetByID(id)
}

func (s *DoctorService) GetAllDoctors() ([]doctor.Doctor, error) {
	return s.repo.GetAll()
}

func (s *DoctorService) UpdateDoctor(id uint, name string, departmentID uint) (*doctor.Doctor, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	doc.Name = name
	doc.DepartmentID = departmentID

	if err := s.repo.Update(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *DoctorService) DeleteDoctor(id uint) error {
	return s.repo.Delete(id)
}

func (s *DoctorService) SetAvailability(id uint, available bool) error {
	return s.repo.SetAvailability(id, available)
}

func (s *DoctorService) GetAvailableDoctors() ([]doctor.Doctor, error) {
	return s.repo.GetAvailable()
}
