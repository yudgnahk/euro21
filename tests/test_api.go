package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/yudgnahk/euro21/adapters/sofascore"
	"github.com/yudgnahk/euro21/config"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		cfg = config.DefaultConfig()
	}
	
	fmt.Printf("Testing Euro 2021 (Tournament ID: %d, Season ID: %d)\n", cfg.Tournament.ID, cfg.Tournament.SeasonID)
	
	// Create client
	client := sofascore.NewClient()
	
	// Setup cache
	if cfg.Cache.Enabled {
		cacheConfig := cfg.ToCacheConfig()
		cache, err := sofascore.NewCache(cacheConfig)
		if err == nil {
			client.SetCache(cache)
			fmt.Println("Cache enabled")
		}
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Test standings
	fmt.Println("\n=== Testing GetStandings ===")
	standings, err := client.GetStandings(ctx, cfg.Tournament.ID, cfg.Tournament.SeasonID)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("SUCCESS: Got %d standings/groups\n", len(standings.Standings))
		for i, standing := range standings.Standings {
			if i < 3 {
				fmt.Printf("  - %s (%d teams)\n", standing.Name, len(standing.Rows))
				for j, row := range standing.Rows {
					if j < 2 {
						fmt.Printf("    %d. %s - %d pts\n", row.Position, row.Team.Name, row.Points)
					}
				}
			}
		}
	}
	
	// Test events
	fmt.Println("\n=== Testing GetTournamentEvents ===")
	events, err := client.GetTournamentEvents(ctx, cfg.Tournament.ID, cfg.Tournament.SeasonID)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("SUCCESS: Got %d events\n", len(events.Events))
		for i, event := range events.Events {
			if i < 5 {
				eventTime := time.Unix(event.StartTimestamp, 0)
				fmt.Printf("  - %s: %s vs %s [%s]\n", 
					eventTime.Format("2006-01-02"), 
					event.HomeTeam.Name, 
					event.AwayTeam.Name,
					event.Status.Type)
			}
		}
	}
}
