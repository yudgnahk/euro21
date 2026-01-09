package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/yudgnahk/euro21/dtos"
	"github.com/yudgnahk/euro21/tablewriter"
	"github.com/yudgnahk/euro21/utils/sliceutil"
	emojiflags "github.com/yudgnahk/go-emoji-flags"
)

type stageFilter struct {
	Display string
	Filter  string
}

var stageFilters = []stageFilter{
	{Display: "All Matches", Filter: ""},
	{Display: "Round of 16", Filter: "Round of 16"},
	{Display: "Quarter-finals", Filter: "Quarter"},
	{Display: "Semi-finals", Filter: "Semi"},
	{Display: "Final", Filter: "Final"},
}

func GetMatch() {
	ctx := context.Background()

	// Show stage selection prompt
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "⚽️ {{ .Display | cyan }}",
		Inactive: "  {{ .Display | cyan }}",
		Selected: "⚽️ {{ .Display | yellow | cyan }}",
	}

	searcher := func(input string, index int) bool {
		stage := stageFilters[index]
		return strings.Contains(strings.ToLower(stage.Display), strings.ToLower(input))
	}

	prompt := promptui.Select{
		Label:     "Select round",
		Items:     stageFilters,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	selectedFilter := stageFilters[i]

	// Fetch tournament events from SofaScore API
	response, err := sofaClient.GetTournamentEvents(ctx, cfg.Tournament.ID, cfg.Tournament.SeasonID)
	if err != nil {
		logrus.Errorf("Failed to fetch tournament events: %v", err)
		fmt.Printf("Error: Failed to fetch tournament events. Please check your internet connection.\n")
		return
	}

	if len(response.Events) == 0 {
		fmt.Println("No events found for this tournament.")
		return
	}

	// Filter events based on selection
	var filteredEvents []dtos.SofaEvent
	if selectedFilter.Filter == "" {
		filteredEvents = response.Events
	} else {
		for _, event := range response.Events {
			if event.RoundInfo.Name != "" && strings.Contains(event.RoundInfo.Name, selectedFilter.Filter) {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}

	if len(filteredEvents) == 0 {
		fmt.Printf("No matches found for %s.\n", selectedFilter.Display)
		return
	}

	// Sort events by start time
	sort.Slice(filteredEvents, func(i, j int) bool {
		return filteredEvents[i].StartTimestamp < filteredEvents[j].StartTimestamp
	})

	// Render table
	table := tablewriter.NewTable(os.Stdout)
	table.SetHeader([]string{"Time", "Match", "Status"})

	for _, event := range filteredEvents {
		// Format time
		eventTime := time.Unix(event.StartTimestamp, 0)
		timeStr := eventTime.Format("Mon Jan 02, 15:04")

		// Format match
		matchStr := getDisplayMatch(event)

		// Get status
		statusStr := getEventStatus(event)

		if err := table.Append(sliceutil.ToStringSlice(timeStr, matchStr, statusStr)); err != nil {
			logrus.Warnf("Failed to append table row: %v", err)
		}
	}

	table.Render()
}

func getDisplayMatch(event dtos.SofaEvent) string {
	homeTeam := formatTeamName(event.HomeTeam)
	awayTeam := formatTeamName(event.AwayTeam)

	// Check if match has been played (status is finished or in progress)
	if event.Status.Type == "finished" || event.Status.Type == "inprogress" {
		return fmt.Sprintf("%s %d - %d %s", homeTeam, event.HomeScore.Current, event.AwayScore.Current, awayTeam)
	} else {
		return fmt.Sprintf("%s vs %s", homeTeam, awayTeam)
	}
}

func formatTeamName(team dtos.SofaTeam) string {
	if countryCode, ok := countriesMap[team.Name]; ok {
		return fmt.Sprintf("%v %v", emojiflags.GetFlag(countryCode), team.Name)
	}
	return team.Name
}

func getEventStatus(event dtos.SofaEvent) string {
	switch event.Status.Type {
	case "finished":
		return "FT"
	case "inprogress":
		return "LIVE"
	case "notstarted":
		return "Scheduled"
	case "postponed":
		return "Postponed"
	case "canceled":
		return "Cancelled"
	default:
		return event.Status.Description
	}
}
