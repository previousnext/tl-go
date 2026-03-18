package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

var mondayConfig = &now.Config{WeekStartDay: time.Monday}

// ParseHumanDate parses a human-friendly date string relative to ref.
// Supported keywords: today, yesterday, this week, last week, this month, last month.
// Falls back to strict YYYY-MM-DD parsing.
// Returns the resolved start/end times, a display label, and any error.
func ParseHumanDate(s string, ref time.Time) (start, end time.Time, label string, err error) {
	n := mondayConfig.With(ref)

	switch strings.ToLower(strings.TrimSpace(s)) {
	case "today":
		start = n.BeginningOfDay()
		end = n.EndOfDay()
		label = "today"
	case "yesterday":
		yesterday := mondayConfig.With(ref.AddDate(0, 0, -1))
		start = yesterday.BeginningOfDay()
		end = yesterday.EndOfDay()
		label = "yesterday"
	case "this week":
		start = n.BeginningOfWeek()
		end = n.EndOfWeek()
		label = "this week"
	case "last week":
		lastWeek := mondayConfig.With(ref.AddDate(0, 0, -7))
		start = lastWeek.BeginningOfWeek()
		end = lastWeek.EndOfWeek()
		label = "last week"
	case "this month":
		start = n.BeginningOfMonth()
		end = n.EndOfMonth()
		label = "this month"
	case "last month":
		// Derive last month from the beginning of the current month to avoid
		// AddDate day overflow issues on end-of-month reference dates.
		currentMonthStart := n.BeginningOfMonth()
		prevMonthRef := currentMonthStart.AddDate(0, 0, -1)
		lastMonth := mondayConfig.With(prevMonthRef)
		start = lastMonth.BeginningOfMonth()
		end = lastMonth.EndOfMonth()
		label = "last month"
	default:
		t, parseErr := time.ParseInLocation(time.DateOnly, strings.TrimSpace(s), ref.Location())
		if parseErr != nil {
			return time.Time{}, time.Time{}, "", fmt.Errorf("unrecognised date %q: expected YYYY-MM-DD or a keyword like 'today', 'yesterday', 'last week'", s)
		}
		d := mondayConfig.With(t)
		start = d.BeginningOfDay()
		end = d.EndOfDay()
		label = t.Format(time.DateOnly)
	}
	return start, end, label, nil
}
