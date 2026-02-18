package model

import (
	"time"
)

// CurrentTimeEntry tracks the ongoing time entry session.
type CurrentTimeEntry struct {
	ID        uint   `gorm:"primaryKey"`
	IssueKey  string `gorm:"index"`
	StartTime time.Time
	Paused    bool
	PauseTime time.Time
	Duration  time.Duration // accumulated duration
}
