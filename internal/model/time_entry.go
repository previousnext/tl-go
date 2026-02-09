package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type TimeEntry struct {
	gorm.Model
	IssueKey    string        `gorm:"index"`
	IssueID     uint          `gorm:"index"`
	Issue       *Issue        `gorm:"foreignkey:IssueID"`
	Duration    time.Duration // Duration in minutes
	Description string
	Sent        bool
}

func FormatDuration(dur time.Duration) string {
	return strings.TrimSuffix(dur.String(), "0s")
}

func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC822)
}
