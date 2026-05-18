package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yudgnahk/euro21/tui"
)

// bracketCmd represents the bracket command
var bracketCmd = &cobra.Command{
	Use:   "bracket",
	Short: "Display the Euro 2021 knockout bracket",
	Long: `Display the Euro 2021 knockout stage bracket with all matches
from Round of 16 through to the Final, including visual connectors
showing the path to victory.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Create and render the visual Euro 2021 bracket
		bracket := tui.NewVisualBracket().BuildEuro2021()

		renderer := tui.NewRenderer(os.Stdout)
		if err := renderer.Render(ctx, bracket); err != nil {
			cmd.PrintErrf("Failed to render bracket: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(bracketCmd)
}
