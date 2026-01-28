package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type TimeEntry struct {
	gorm.Model
	IssueKey    string `gorm:"index"`
	Duration    uint   // Duration in minutes
	Description string
}

func ParseDuration(duration string) (uint, error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}
	return uint(d.Minutes()), nil
}

func FormatDuration(minutes uint) string {
	duration := time.Duration(minutes) * time.Minute
	return strings.TrimSuffix(duration.String(), "0s")
}

func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC822)
}
