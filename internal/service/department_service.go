package service

import (
	"errors"
	"medical-center/internal/models/department"
	"medical-center/internal/repository"
	"time"
)

type DepartmentService struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) CreateDepartment(name string) (*department.Department, error) {
	if name == "" {
		return nil, errors.New("department name cannot be empty")
	}

	newDept := &department.Department{
		Name: name,
	}

	if err := s.repo.Create(newDept); err != nil {
		return nil, err
	}
	return newDept, nil
}

func (s *DepartmentService) GetDepartmentByID(id uint) (*department.Department, error) {
	return s.repo.GetByID(id)
}

func (s *DepartmentService) GetAllDepartments() ([]department.Department, error) {
	return s.repo.GetAll()
}

func (s *DepartmentService) UpdateDepartment(id uint, name string) (*department.Department, error) {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	dept.Name = name

	if err := s.repo.Update(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *DepartmentService) DeleteDepartment(id uint) error {
	return s.repo.Delete(id)
}

func (s *DepartmentService) GetWithDoctors(id uint) (*department.Department, error) {
	return s.repo.GetWithDoctors(id)
}

func (s *DepartmentService) GetAvailableSlots(id uint, date string) ([]string, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	slots, err := s.repo.GetAvailableSlots(id, parsedDate)
	if err != nil {
		return nil, err
	}

	formattedSlots := make([]string, 0, len(slots))
	for _, slot := range slots {
		formattedSlots = append(formattedSlots, slot.Format("15:04"))
	}

	return formattedSlots, nil
}
