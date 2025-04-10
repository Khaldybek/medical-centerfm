package schedule

import "time"

type SchduleDTO struct {
	DoctorID  uint
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
	Booked    bool
}
