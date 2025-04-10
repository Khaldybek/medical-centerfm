package doctor

import "medical-center/internal/models/schedule"

type DoctorDTO struct {
	Name         string
	DepartmentID uint
	Schedule     []schedule.Schedule
}
