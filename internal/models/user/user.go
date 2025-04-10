package user

import (
	"time"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleDoctor Role = "doctor"
	RolePatient Role = "patient"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // Password is not exposed in JSON
	Name      string    `json:"name"`
	Role      Role      `json:"role" gorm:"type:varchar(20);default:'patient'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 