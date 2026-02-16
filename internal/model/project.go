package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Key        string `gorm:"index"`
	Name       string
	CategoryID *uint
	Category   *Category `gorm:"foreignkey:CategoryID"`
}

type Category struct {
	gorm.Model
	Name string
}
