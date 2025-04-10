package doctor

import (
	"gorm.io/gorm"
	"medical-center/internal/models/schedule"
)

type Doctor struct {
	gorm.Model
	Name         string              `gorm:"not null;size:100"`
	DepartmentID uint                `gorm:"index;not null"`
	Available    bool                `gorm:"default:true"`
	Schedule     []schedule.Schedule `gorm:"foreignKey:DoctorID"` // Связь с расписанием
}
