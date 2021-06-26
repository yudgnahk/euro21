package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/yudgnahk/euro21/adapters"
	"github.com/yudgnahk/euro21/constants"
	"github.com/yudgnahk/euro21/dtos"
	"github.com/yudgnahk/euro21/tablewriter"
	"github.com/yudgnahk/euro21/utils/sliceutil"
	"github.com/yudgnahk/euro21/utils/stringutil"
	emojiflags "github.com/yudgnahk/go-emoji-flags"
)

type stageSelection struct {
	Display string
	Name    string
}

var stageSelections = []stageSelection{
	{
		Display: constants.StageGroupA,
		Name:    adapters.GroupA,
	},
	{
		Display: constants.StageGroupB,
		Name:    adapters.GroupB,
	},
	{
		Display: constants.StageGroupC,
		Name:    adapters.GroupC,
	},
	{
		Display: constants.StageGroupD,
		Name:    adapters.GroupD,
	},
	{
		Display: constants.StageGroupE,
		Name:    adapters.GroupE,
	},
	{
		Display: constants.StageGroupF,
		Name:    adapters.GroupF,
	},
	{
		Display: constants.RoundOf16,
		Name:    adapters.RoundOf16,
	},
	{
		Display: constants.QuarterFinals,
		Name:    adapters.QuarterFinals,
	},
	{
		Display: constants.SemiFinals,
		Name:    adapters.SemiFinals,
	},
	{
		Display: constants.Final,
		Name:    adapters.Final,
	},
}

func GetMatch() {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "⚽️ {{ .Display | cyan }}",
		Inactive: "  {{ .Display | cyan }}",
		Selected: "⚽️ {{ .Display | yellow | cyan }}",
	}

	searcher := func(input string, index int) bool {
		stage := stageSelections[index]

		return strings.Contains(stage.Display, input)
	}

	prompt := promptui.Select{
		Label:     "Select round",
		Items:     stageSelections,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	data, _ := adapters.GetStage(stageSelections[i].Name)

	stageData := data.Stages[0]
	matches := stageData.Events

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Esd < matches[j].Esd
	})

	table := tablewriter.NewTable(os.Stdout)
	table.SetHeader([]string{"Time", "Match"})

	for _, m := range matches {
		table.Append(sliceutil.ToStringSlice(getTime(m.Esd).Format(time.RFC1123), getDisplayMatch(m)))
	}

	table.Render()
}

func getTime(z int64) time.Time {
	var (
		s      = fmt.Sprintf("%v", z)
		year   = stringutil.ToInt(s[0:4])
		month  = stringutil.ToInt(s[4:6])
		date   = stringutil.ToInt(s[6:8])
		hour   = stringutil.ToInt(s[8:10])
		minute = stringutil.ToInt(s[10:12])
		second = stringutil.ToInt(s[12:])
	)

	t := time.Date(year, time.Month(month), date, hour, minute, second, 0, time.Local)
	return t
}

func getDisplayMatch(e dtos.Event) string {
	if len(e.Tr1) > 0 {
		return fmt.Sprintf("%v %v - %v %v", getTeamName(e.T1[0]), e.Tr1, e.Tr2, getTeamName(e.T2[0]))
	} else {
		return fmt.Sprintf("%v ? - ? %v", getTeamName(e.T1[0]), getTeamName(e.T2[0]))
	}
}

func getTeamName(t dtos.Team) string {
	if t.Tbd == 0 {
		return fmt.Sprintf("%v %v", emojiflags.GetFlag(countriesMap[t.Nm]), t.Nm)
	} else {
		return t.Nm
	}
}
