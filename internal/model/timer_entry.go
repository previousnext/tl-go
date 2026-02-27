package model

import (
	"time"
)

// TimerEntry tracks the ongoing time entry session.
type TimerEntry struct {
	ID             uint   `gorm:"primaryKey"`
	IssueKey       string `gorm:"index"`
	StartTime      time.Time
	LastActiveTime time.Time
	Paused         bool
	PauseTime      time.Time
	Duration       time.Duration // accumulated duration
	Description    *string
}
