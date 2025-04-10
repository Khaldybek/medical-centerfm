package service

import (
	"errors"
	"medical-center/internal/models/appointment"
	"medical-center/internal/repository"
	"time"
)

type AppointmentService struct {
	repo repository.AppRepository
}

func NewAppointmentService(repo repository.AppRepository) *AppointmentService {
	return &AppointmentService{repo: repo}
}

func (s *AppointmentService) CreateAppointment(
	patientName, email, phone string,
	departmentID, doctorID uint,
	appointmentTime time.Time,
) (*appointment.Appointment, error) {

	// Валидация данных
	if patientName == "" {
		return nil, errors.New("patient name is required")
	}
	if appointmentTime.Before(time.Now()) {
		return nil, errors.New("appointment time cannot be in the past")
	}

	newAppointment := &appointment.Appointment{
		PatientName:     patientName,
		Email:           email,
		Phone:           phone,
		DepartmentID:    departmentID,
		DoctorID:        doctorID,
		AppointmentTime: appointmentTime,
	}

	if err := s.repo.Create(newAppointment); err != nil {
		return nil, err
	}
	return newAppointment, nil
}

func (s *AppointmentService) GetAppointmentByID(id uint) (*appointment.Appointment, error) {
	return s.repo.GetByID(id)
}

func (s *AppointmentService) GetAllAppointments() ([]appointment.Appointment, error) {
	return s.repo.GetAll()
}

func (s *AppointmentService) UpdateAppointment(
	id uint,
	patientName, email, phone string,
	departmentID, doctorID uint,
	appointmentTime time.Time,
) (*appointment.Appointment, error) {

	appt, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Обновление полей
	appt.PatientName = patientName
	appt.Email = email
	appt.Phone = phone
	appt.DepartmentID = departmentID
	appt.DoctorID = doctorID
	appt.AppointmentTime = appointmentTime

	if err := s.repo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}

func (s *AppointmentService) DeleteAppointment(id uint) error {
	return s.repo.Delete(id)
}

func (s *AppointmentService) GetByDepartment(departmentID uint) ([]appointment.Appointment, error) {
	return s.repo.GetByDepartment(departmentID)
}

func (s *AppointmentService) GetByPatient(patientName string) ([]appointment.Appointment, error) {
	return s.repo.GetByPatient(patientName)
}
