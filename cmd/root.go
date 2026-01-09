package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yudgnahk/euro21/adapters/sofascore"
	"github.com/yudgnahk/euro21/config"
	"github.com/yudgnahk/euro21/dtos"
)

var (
	countriesMap = make(map[string]string)
	sofaClient   *sofascore.Client
	cfg          *config.Config
)

//go:embed countries.json
var countriesData string

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
	// Load configuration
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		logrus.Warnf("Failed to load config, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	// Initialize SofaScore client
	sofaClient = sofascore.NewClient()

	// Setup cache if enabled
	cacheConfig := cfg.ToCacheConfig()
	if cacheConfig.Enabled {
		cache, err := sofascore.NewCache(cacheConfig)
		if err != nil {
			logrus.Warnf("Failed to initialize cache: %v", err)
		} else {
			sofaClient.SetCache(cache)
			logrus.Debugf("Cache enabled: %s (TTL: %v)", cacheConfig.CacheDir, cacheConfig.TTL)
		}
	}

	// Setup proxy if enabled
	proxyConfig := cfg.ToProxyConfig()
	if proxyConfig.Enabled {
		sofaClient.SetProxyConfig(&proxyConfig)
		logrus.Debugf("Proxy enabled: %s (%s)", proxyConfig.Address, proxyConfig.Type)
	}

	// Load countries mapping
	countries := make([]dtos.Country, 0)
	err = json.Unmarshal([]byte(countriesData), &countries)
	if err != nil {
		logrus.Warnf("Failed to load countries data: %v", err)
	}

	for i := range countries {
		countriesMap[countries[i].Name] = countries[i].Code
	}
}
