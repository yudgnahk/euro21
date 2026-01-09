package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yudgnahk/euro21/constants"
	"github.com/yudgnahk/euro21/dtos"
	"github.com/yudgnahk/euro21/tui"
	emojiflags "github.com/yudgnahk/go-emoji-flags"
)

const (
	priorityPlaces = constants.ColorGreen
	bestThirdPlace = constants.ColorYellow
	failedPlace    = constants.ColorWhite
)

func GetTable() {
	ctx := context.Background()

	// Fetch standings from SofaScore API
	response, err := sofaClient.GetStandings(ctx, cfg.Tournament.ID, cfg.Tournament.SeasonID)
	if err != nil {
		logrus.Errorf("Failed to fetch standings: %v", err)
		fmt.Printf("Error: Failed to fetch standings. Please check your internet connection.\n")
		return
	}

	if len(response.Standings) == 0 {
		fmt.Println("No standings data available.")
		return
	}

	// Extract third place teams from all groups for ranking
	thirdPlaceTeams := make([]dtos.SofaStandingRow, 0)
	for _, standing := range response.Standings {
		if len(standing.Rows) >= 3 {
			thirdPlaceTeams = append(thirdPlaceTeams, standing.Rows[2])
		}
	}

	// Sort third place teams by points, goal difference, and goals for
	sort.Slice(thirdPlaceTeams, func(i, j int) bool {
		if thirdPlaceTeams[i].Points != thirdPlaceTeams[j].Points {
			return thirdPlaceTeams[i].Points > thirdPlaceTeams[j].Points
		}

		gdI := thirdPlaceTeams[i].ScoresFor - thirdPlaceTeams[i].ScoresAgainst
		gdJ := thirdPlaceTeams[j].ScoresFor - thirdPlaceTeams[j].ScoresAgainst

		if gdI != gdJ {
			return gdI > gdJ
		}

		return thirdPlaceTeams[i].ScoresFor > thirdPlaceTeams[j].ScoresFor
	})

	// Best 4 third place teams qualify
	bestThirdPlaces := make(map[int]bool)
	for i := 0; i < 4 && i < len(thirdPlaceTeams); i++ {
		bestThirdPlaces[thirdPlaceTeams[i].Team.ID] = true
	}

	// Create renderer and multi-table using new TUI package
	renderer := tui.NewRenderer(os.Stdout)
	multiTable := tui.NewMultiTable()

	// Sort standings by name (Group A, Group B, etc.)
	sort.Slice(response.Standings, func(i, j int) bool {
		return response.Standings[i].Name < response.Standings[j].Name
	})

	for _, standing := range response.Standings {
		// Create table for this group
		table := tui.NewTable().
			SetHeaders("Name", "P", "W", "D", "L", "F", "A", "GD").
			// Set right alignment for all number columns (columns 1-7)
			SetAlignment(1, tui.AlignRight). // P
			SetAlignment(2, tui.AlignRight). // W
			SetAlignment(3, tui.AlignRight). // D
			SetAlignment(4, tui.AlignRight). // L
			SetAlignment(5, tui.AlignRight). // F
			SetAlignment(6, tui.AlignRight). // A
			SetAlignment(7, tui.AlignRight)  // GD

		for i, row := range standing.Rows {
			// Get team name with flag emoji
			teamName := row.Team.Name
			if countryCode, ok := countriesMap[teamName]; ok {
				flag := emojiflags.GetFlag(countryCode)
				// emojiflags library adds trailing space for some flags but not others
				// Trim it and add consistent spacing
				flag = strings.TrimSpace(flag)
				flagAndName := fmt.Sprintf("%v  %v", flag, teamName)
				teamName = flagAndName
			}

			// Calculate goal difference
			gd := row.ScoresFor - row.ScoresAgainst

			// Determine row color based on position
			var color tui.Color
			switch i {
			case 0, 1:
				// First and second place qualify (green)
				color = tui.NewColor(priorityPlaces)
			case 2:
				// Third place - check if they're in best 4
				if bestThirdPlaces[row.Team.ID] {
					color = tui.NewColor(bestThirdPlace)
				} else {
					color = tui.NewColor(failedPlace)
				}
			default:
				// Fourth place and below (white)
				color = tui.NewColor(failedPlace)
			}

			// Add row with color
			table.AddRowWithColor(color,
				teamName,
				fmt.Sprint(row.Points),
				fmt.Sprint(row.Wins),
				fmt.Sprint(row.Draws),
				fmt.Sprint(row.Losses),
				fmt.Sprint(row.ScoresFor),
				fmt.Sprint(row.ScoresAgainst),
				fmt.Sprint(gd),
			)
		}

		// Add section to multi-table
		multiTable.AddSection(standing.Name, table)
	}

	// Normalize widths across all tables to make them the same size
	multiTable.NormalizeWidths()

	// Render the multi-table
	if err := renderer.Render(ctx, multiTable); err != nil {
		logrus.Errorf("Failed to render table: %v", err)
	}
}
