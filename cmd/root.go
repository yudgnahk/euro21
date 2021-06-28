package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yudgnahk/euro21/dtos"
)

var countriesMap = make(map[string]string)

type menu struct {
	Choice string
}

var choices = []menu{
	{
		Choice: "Get group stage table",
	},
	{
		Choice: "Get match list",
	},
}

// RootCmd is the root command of kit
var RootCmd = &cobra.Command{
	Use:   "euro21",
	Short: "Euro 2021 CLI",
	Run: func(cmd *cobra.Command, args []string) {
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "👉🏻 {{ .Choice | cyan }}",
			Inactive: "  {{ .Choice | cyan }}",
			Selected: "👉🏻 {{ .Choice | yellow | cyan }}",
		}

		searcher := func(input string, index int) bool {
			choice := choices[index]

			return strings.Contains(choice.Choice, input)
		}

		prompt := promptui.Select{
			Label:     "Select menu",
			Items:     choices,
			Templates: templates,
			Size:      10,
			Searcher:  searcher,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch i {
		case 0:
			GetTable()
		case 1:
			GetMatch()
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	countries := make([]dtos.Country, 0)
	bytes, err := ioutil.ReadFile("./data/countries.json")
	if err != nil {
		fmt.Println("read countries data got error: ", err)
	}

	_ = json.Unmarshal(bytes, &countries)

	for i := range countries {
		countriesMap[countries[i].Name] = countries[i].Code
	}

	fmt.Println("finish init countries data")
}
