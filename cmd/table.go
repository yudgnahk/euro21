package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/yudgnahk/euro21/adapters"
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
	data, _ := adapters.GetTables()

	// get rank to set color
	thirdPlaces := make([]dtos.TeamData, 0)

	for _, stage := range data.Stages {
		teams := stage.LeagueTable.L[0].Tables[0].Team

		for i := range teams {
			if i == 2 {
				thirdPlaces = append(thirdPlaces, teams[i])
			}
		}
	}

	sort.Slice(thirdPlaces, func(i, j int) bool {
		if thirdPlaces[i].Points != thirdPlaces[j].Points {
			if thirdPlaces[i].Points < thirdPlaces[j].Points {
				return false
			}

			return true
		}

		if thirdPlaces[i].Gd != thirdPlaces[j].Gd {
			if thirdPlaces[i].Gd < thirdPlaces[j].Gd {
				return false
			}

			return true
		}

		if thirdPlaces[i].Ga != thirdPlaces[j].Ga {
			if thirdPlaces[i].Ga < thirdPlaces[j].Ga {
				return false
			}

			return true
		}

		return false
	})

	bestThirdPlaces := make(map[string]dtos.TeamData)

	for i := 0; i < 4; i++ {
		bestThirdPlaces[thirdPlaces[i].Name] = thirdPlaces[i]
	}

	multiTables := tablewriter.NewMultiTables(os.Stdout)
	multiTables.SetHeaders([]string{"Name", "P", "W", "D", "L", "F", "A", "GD"})

	for _, stage := range data.Stages {

		multiTables.AppendSubHeaders(stage.StageName)

		teams := stage.LeagueTable.L[0].Tables[0].Team

		tableDetail := tablewriter.TableData{}
		for i, team := range teams {
			flagAndName := fmt.Sprintf("%v %v", emojiflags.GetFlag(countriesMap[team.Name]), team.Name)
			tableDetail.Data = append(tableDetail.Data,
				sliceutil.ToStringSlice(flagAndName, team.Points, team.Win, team.Draw, team.Lost, team.Gf, team.Ga, team.Gd))

			switch i {
			case 0, 1:
				tableDetail.Color = append(tableDetail.Color, priorityPlaces)
			case 2:
				if _, ok := bestThirdPlaces[team.Name]; ok {
					tableDetail.Color = append(tableDetail.Color, bestThirdPlace)
				} else {
					tableDetail.Color = append(tableDetail.Color, failedPlace)
				}
			case 3:
				tableDetail.Color = append(tableDetail.Color, failedPlace)
			}
		}

		multiTables.AppendTable(tableDetail)
	}

	multiTables.Render()
}
