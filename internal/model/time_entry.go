package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TimeEntry struct {
	gorm.Model
	IssueKey        string        `gorm:"index"`
	IssueID         uint          `gorm:"index"`
	Issue           *Issue        `gorm:"foreignkey:IssueID"`
	Duration        time.Duration // Duration in minutes
	AISavedDuration time.Duration // Duration of time saved by AI
	Description     string
	Sent            bool
}

func FormatDuration(dur time.Duration) string {
	h := int(dur.Hours())
	m := int(dur.Minutes()) % 60
	s := int(dur.Seconds()) % 60
	var parts []string
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	if s > 0 {
		parts = append(parts, fmt.Sprintf("%ds", s))
	}
	if len(parts) == 0 {
		return "0m"
	}
	return strings.Join(parts, " ")
}

func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC822)
}
