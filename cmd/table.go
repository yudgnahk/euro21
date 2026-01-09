package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/yudgnahk/euro21/constants"
	"github.com/yudgnahk/euro21/dtos"
	"github.com/yudgnahk/euro21/tablewriter"
	"github.com/yudgnahk/euro21/utils/sliceutil"
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

	// Render multi-table layout
	multiTables := tablewriter.NewMultiTables(os.Stdout)
	multiTables.SetHeaders([]string{"Name", "P", "W", "D", "L", "F", "A", "GD"})

	// Sort standings by name (Group A, Group B, etc.)
	sort.Slice(response.Standings, func(i, j int) bool {
		return response.Standings[i].Name < response.Standings[j].Name
	})

	for _, standing := range response.Standings {
		multiTables.AppendSubHeaders(standing.Name)

		tableDetail := tablewriter.TableData{}
		for i, row := range standing.Rows {
			// Get team name with flag emoji
			teamName := row.Team.Name
			if countryCode, ok := countriesMap[teamName]; ok {
				flagAndName := fmt.Sprintf("%v %v", emojiflags.GetFlag(countryCode), teamName)
				teamName = flagAndName
			}

			// Calculate goal difference
			gd := row.ScoresFor - row.ScoresAgainst

			// Append row data
			tableDetail.Data = append(tableDetail.Data,
				sliceutil.ToStringSlice(teamName, row.Points, row.Wins, row.Draws, row.Losses, row.ScoresFor, row.ScoresAgainst, gd))

			// Determine row color based on position
			switch i {
			case 0, 1:
				// First and second place qualify (green)
				tableDetail.Color = append(tableDetail.Color, priorityPlaces)
			case 2:
				// Third place - check if they're in best 4
				if bestThirdPlaces[row.Team.ID] {
					tableDetail.Color = append(tableDetail.Color, bestThirdPlace)
				} else {
					tableDetail.Color = append(tableDetail.Color, failedPlace)
				}
			default:
				// Fourth place and below (white)
				tableDetail.Color = append(tableDetail.Color, failedPlace)
			}
		}

		if err := multiTables.AppendTable(tableDetail); err != nil {
			logrus.Warnf("Failed to append table: %v", err)
		}
	}

	multiTables.Render()
}
