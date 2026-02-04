package model

import "gorm.io/gorm"

type Issue struct {
	gorm.Model
	Key       string `gorm:"index"`
	Summary   string
	ProjectID *uint
	Project   *Project `gorm:"foreignkey:ProjectID"`
}
