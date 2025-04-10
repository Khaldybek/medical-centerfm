package models

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
