package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/yudgnahk/euro21/adapters/sofascore"
	"github.com/yudgnahk/euro21/config"
)

func main() {
	cfg, _ := config.LoadConfig()
	client := sofascore.NewClient()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	events, err := client.GetTournamentEvents(ctx, cfg.Tournament.ID, cfg.Tournament.SeasonID)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Printf("Total events: %d\n\n", len(events.Events))
	
	// Check round info
	roundCounts := make(map[string]int)
	for _, event := range events.Events {
		if event.RoundInfo.Name != "" {
			roundCounts[event.RoundInfo.Name]++
		} else {
			roundCounts["<empty>"]++
		}
	}
	
	fmt.Println("Rounds found:")
	for round, count := range roundCounts {
		fmt.Printf("  %s: %d events\n", round, count)
	}
}
