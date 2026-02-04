package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Key  string `gorm:"index"`
	Name string
}
