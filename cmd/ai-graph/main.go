package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/api/types"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

// daySummary holds aggregated worklog data for a single day.
type daySummary struct {
	Date       time.Time
	Logged     time.Duration
	AISaved    time.Duration
	WorklogIDs int
}

func main() {
	dateFlag := flag.String("date", "this week", "Date range: 'today', 'yesterday', 'this week', 'last week', 'this month', 'last month', or YYYY-MM-DD")
	configFlag := flag.String("config", "", "Path to tl config file (default ~/.config/tl/config.yml)")
	flag.Parse()

	// Load config
	if err := loadConfig(*configFlag); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	baseURL := viper.GetString("jira_base_url")
	username := viper.GetString("jira_username")
	apiToken := viper.GetString("jira_api_token")

	if baseURL == "" || username == "" || apiToken == "" {
		log.Fatal("Missing Jira configuration. Ensure jira_base_url, jira_username, and jira_api_token are set in your tl config file.")
	}

	// Parse date range
	start, end, label, err := util.ParseHumanDate(*dateFlag, time.Now())
	if err != nil {
		log.Fatalf("Invalid date: %v", err)
	}

	fmt.Printf("Fetching worklogs for %s ...\n", label)

	// Build Jira client
	jiraClient := api.NewJiraClient(&http.Client{}, types.JiraClientParams{
		BaseURL:  baseURL,
		Username: username,
		APIToken: apiToken,
	})

	// Fetch all updated worklog IDs since the start of the period
	changes, err := jiraClient.GetUpdatedWorklogIDs(start.UnixMilli())
	if err != nil {
		log.Fatalf("Failed to fetch updated worklog IDs: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("No worklogs found for this period.")
		return
	}

	// Filter to only changes within our date range
	endMillis := end.UnixMilli()
	var filteredIDs []int64
	for _, c := range changes {
		if c.UpdatedTime <= endMillis {
			filteredIDs = append(filteredIDs, c.WorklogID)
		}
	}

	if len(filteredIDs) == 0 {
		fmt.Println("No worklogs found for this period.")
		return
	}

	fmt.Printf("Found %d worklog updates, fetching details...\n", len(filteredIDs))

	// Bulk fetch worklog details
	worklogs, err := jiraClient.BulkGetWorklogs(filteredIDs)
	if err != nil {
		log.Fatalf("Failed to fetch worklog details: %v", err)
	}

	// Filter worklogs to those whose started date falls within our range.
	// The "started" field format from Jira is "2021-01-17T12:34:00.000+0000".
	var filtered []types.Worklog
	for _, w := range worklogs {
		started, err := time.Parse(api.DateFormat, w.Started)
		if err != nil {
			// Try alternate format
			started, err = time.Parse("2006-01-02T15:04:05.000+0000", w.Started)
			if err != nil {
				continue
			}
		}
		if !started.Before(start) && !started.After(end) {
			filtered = append(filtered, w)
		}
	}

	if len(filtered) == 0 {
		fmt.Println("No worklogs with matching start dates found for this period.")
		return
	}

	fmt.Printf("Fetching AI time saved properties for %d worklogs...\n", len(filtered))

	// For each worklog, fetch the AI time saved property and aggregate by day
	days := make(map[string]*daySummary)

	for _, w := range filtered {
		started, _ := time.Parse(api.DateFormat, w.Started)
		if started.IsZero() {
			started, _ = time.Parse("2006-01-02T15:04:05.000+0000", w.Started)
		}
		dayKey := started.Format(time.DateOnly)

		if _, ok := days[dayKey]; !ok {
			days[dayKey] = &daySummary{
				Date: time.Date(started.Year(), started.Month(), started.Day(), 0, 0, 0, 0, started.Location()),
			}
		}

		ds := days[dayKey]
		ds.Logged += time.Duration(w.TimeSpentSeconds) * time.Second
		ds.WorklogIDs++

		// Fetch AI time saved property
		prop, err := jiraClient.GetWorklogProperty(w.IssueID, w.ID, "ai-time-saved")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to fetch AI property for worklog %s: %v\n", w.ID, err)
			continue
		}
		if prop != nil {
			var val types.AITimeSavedPropertyValue
			if err := json.Unmarshal(prop.Value, &val); err == nil {
				ds.AISaved += time.Duration(val.DurationSeconds) * time.Second
			}
		}
	}

	// Build sorted list of day summaries, filling in empty days
	summaries := buildDaySummaries(days, start, end)

	// Render the chart
	renderChart(summaries, label)
}

// loadConfig reads the tl configuration file using viper.
func loadConfig(configPath string) error {
	if configPath != "" {
		viper.SetConfigFile(configPath)
		return viper.ReadInConfig()
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user config directory: %w", err)
	}

	configDir := filepath.Join(userConfigDir, "tl")
	viper.AddConfigPath(configDir)
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}
	return nil
}

// buildDaySummaries creates a sorted slice of daySummary with entries for every day in the range.
func buildDaySummaries(days map[string]*daySummary, start, end time.Time) []daySummary {
	var summaries []daySummary

	current := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	for !current.After(endDay) {
		key := current.Format(time.DateOnly)
		if ds, ok := days[key]; ok {
			summaries = append(summaries, *ds)
		} else {
			summaries = append(summaries, daySummary{Date: current})
		}
		current = current.AddDate(0, 0, 1)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date.Before(summaries[j].Date)
	})

	return summaries
}

// renderChart displays a bar chart of daily time logged vs AI time saved.
func renderChart(summaries []daySummary, label string) {
	orange := gchalk.Hex("#ee5622")
	cyan := gchalk.Hex("#00bcd4")
	dim := gchalk.Dim
	bold := gchalk.Bold

	const barMaxWidth = 40

	// Find max duration for scaling
	var maxLogged time.Duration
	var totalLogged, totalAISaved time.Duration

	for _, ds := range summaries {
		if ds.Logged > maxLogged {
			maxLogged = ds.Logged
		}
		totalLogged += ds.Logged
		totalAISaved += ds.AISaved
	}

	fmt.Println()
	fmt.Println(bold(fmt.Sprintf("AI Time Savings - %s", label)))
	fmt.Println()

	for _, ds := range summaries {
		dayLabel := ds.Date.Format("Mon 02 Jan")

		if ds.Logged == 0 {
			fmt.Printf("  %s  %s\n", dayLabel, dim("-"))
			continue
		}

		// Calculate bar widths
		loggedWidth := 0
		aiWidth := 0
		if maxLogged > 0 {
			loggedWidth = int(math.Round(float64(ds.Logged) / float64(maxLogged) * barMaxWidth))
			aiWidth = int(math.Round(float64(ds.AISaved) / float64(maxLogged) * barMaxWidth))
		}
		if loggedWidth < 1 && ds.Logged > 0 {
			loggedWidth = 1
		}

		// The AI bar is drawn within the logged bar to show the proportion
		// [=====AI=====|---rest---]
		if aiWidth > loggedWidth {
			aiWidth = loggedWidth
		}

		restWidth := loggedWidth - aiWidth

		bar := cyan(strings.Repeat("█", aiWidth)) + orange(strings.Repeat("█", restWidth))
		padding := strings.Repeat(" ", barMaxWidth-loggedWidth)

		pct := float64(0)
		if ds.Logged > 0 {
			pct = float64(ds.AISaved) / float64(ds.Logged) * 100
		}

		stats := fmt.Sprintf("%s logged, %s AI saved",
			model.FormatDuration(ds.Logged),
			model.FormatDuration(ds.AISaved),
		)
		if ds.AISaved > 0 {
			stats += fmt.Sprintf(" (%s)", bold(fmt.Sprintf("%.0f%%", pct)))
		}

		fmt.Printf("  %s  %s%s  %s\n", dayLabel, bar, padding, stats)
	}

	// Footer
	fmt.Println()
	totalPct := float64(0)
	if totalLogged > 0 {
		totalPct = float64(totalAISaved) / float64(totalLogged) * 100
	}

	fmt.Printf("  %s  %s logged, %s AI saved",
		bold("Total      "),
		bold(model.FormatDuration(totalLogged)),
		bold(model.FormatDuration(totalAISaved)),
	)
	if totalAISaved > 0 {
		fmt.Printf(" (%s)", bold(fmt.Sprintf("%.0f%%", totalPct)))
	}
	fmt.Println()
	fmt.Println()
	fmt.Printf("  %s Time logged  %s AI time saved\n", orange("█"), cyan("█"))
	fmt.Println()
}
