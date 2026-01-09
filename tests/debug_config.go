package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/yudgnahk/euro21/config"
)

func main() {
	configPaths := []string{
		"./euro21.config.json",
		filepath.Join(os.Getenv("HOME"), ".config", "euro21", "config.json"),
	}
	
	fmt.Println("Checking config paths:")
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("  ✓ %s (exists)\n", path)
			data, _ := os.ReadFile(path)
			var cfg map[string]interface{}
			json.Unmarshal(data, &cfg)
			if tournament, ok := cfg["tournament"].(map[string]interface{}); ok {
				fmt.Printf("    Season ID: %.0f\n", tournament["season_id"])
			}
		} else {
			fmt.Printf("  ✗ %s (not found)\n", path)
		}
	}
	
	cfg, _ := config.LoadConfig()
	fmt.Printf("\nLoaded config: Tournament ID=%d, Season ID=%d\n", cfg.Tournament.ID, cfg.Tournament.SeasonID)
}
