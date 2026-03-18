package util

import (
	"testing"
	"time"
)

// ref is a fixed Wednesday 2026-03-18 12:00:00 UTC for deterministic tests.
var ref = time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)

func TestParseHumanDate(t *testing.T) {
	tests := []struct {
		input     string
		wantStart time.Time
		wantEnd   time.Time
		wantLabel string
		wantErr   bool
	}{
		{
			input:     "today",
			wantStart: time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 18, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "today",
		},
		{
			input:     "TODAY",
			wantStart: time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 18, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "today",
		},
		{
			input:     "yesterday",
			wantStart: time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 17, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "yesterday",
		},
		{
			// ref is Wednesday 2026-03-18; week starts Monday, so this week = Mon 2026-03-16
			input:     "this week",
			wantStart: time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 22, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "this week",
		},
		{
			// last week = Mon 2026-03-09 to Sun 2026-03-15
			input:     "last week",
			wantStart: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 15, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "last week",
		},
		{
			input:     "Last Week",
			wantStart: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 15, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "last week",
		},
		{
			// ref is in March 2026; this month = 2026-03-01 to 2026-03-31
			input:     "this month",
			wantStart: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 3, 31, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "this month",
		},
		{
			// last month = 2026-02-01 to 2026-02-28
			input:     "last month",
			wantStart: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 2, 28, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "last month",
		},
		{
			// YYYY-MM-DD fallback
			input:     "2026-01-15",
			wantStart: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 1, 15, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC),
			wantLabel: "2026-01-15",
		},
		{
			input:   "not-a-date",
			wantErr: true,
		},
		{
			input:   "2026-99-99",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			start, end, label, err := ParseHumanDate(tt.input, ref)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !start.Equal(tt.wantStart) {
				t.Errorf("start: got %v, want %v", start, tt.wantStart)
			}
			if !end.Equal(tt.wantEnd) {
				t.Errorf("end: got %v, want %v", end, tt.wantEnd)
			}
			if label != tt.wantLabel {
				t.Errorf("label: got %q, want %q", label, tt.wantLabel)
			}
		})
	}
}

func TestParseHumanDate_LastMonth_EndOfMonthRef(t *testing.T) {
	endOfMonthRef := time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC)

	gotStart, gotEnd, gotLabel, err := ParseHumanDate("last month", endOfMonthRef)
	if err != nil {
		t.Fatalf("ParseHumanDate returned error: %v", err)
	}

	wantStart := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2026, 2, 28, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
	wantLabel := "last month"

	if !gotStart.Equal(wantStart) || !gotEnd.Equal(wantEnd) || gotLabel != wantLabel {
		t.Fatalf("ParseHumanDate(%q, %v) = (%v, %v, %q), want (%v, %v, %q)",
			"last month", endOfMonthRef, gotStart, gotEnd, gotLabel, wantStart, wantEnd, wantLabel)
	}
}
