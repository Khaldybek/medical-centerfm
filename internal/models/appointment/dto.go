package appointment

import "time"

type AppointmentDTO struct {
	DepartmentID    uint
	PatientName     string
	Email           string
	Phone           string
	AppointmentTime time.Time
}
