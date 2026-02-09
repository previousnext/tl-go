package model

import "gorm.io/gorm"

type Issue struct {
	gorm.Model
	Key       string `gorm:"uniqueIndex"`
	Summary   string
	ProjectID uint
	Project   Project `gorm:"foreignkey:ProjectID"`
}
