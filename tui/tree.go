package tui

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// Match represents a tournament match
type Match struct {
	ID        int
	HomeTeam  string
	AwayTeam  string
	HomeScore int
	AwayScore int
	Status    string // "Scheduled", "Live", "Finished"
	Round     string // "Round of 16", "Quarter-finals", etc.
	StartTime time.Time
	ExtraInfo string // "aet", "pen", etc.
}

// BracketNode represents a node in the knockout bracket tree
type BracketNode struct {
	Match    Match
	Left     *BracketNode // Child match 1
	Right    *BracketNode // Child match 2
	Parent   *BracketNode
	Level    int // 0=Final, 1=Semi, 2=Quarter, 3=R16
	Position int // Position within the level
}

// BracketTree renders tournament knockout brackets
type BracketTree struct {
	rounds  map[string][]*BracketNode
	theme   *Theme
	compact bool
}

// NewBracketTree creates a new bracket tree widget
func NewBracketTree() *BracketTree {
	return &BracketTree{
		rounds:  make(map[string][]*BracketNode),
		theme:   DefaultTheme,
		compact: false,
	}
}

// SetCompact enables compact rendering mode
func (bt *BracketTree) SetCompact(compact bool) *BracketTree {
	bt.compact = compact
	return bt
}

// SetTheme sets the theme
func (bt *BracketTree) SetTheme(theme *Theme) *BracketTree {
	bt.theme = theme
	return bt
}

// BuildFromMatches builds the bracket tree from a list of matches
func (bt *BracketTree) BuildFromMatches(matches []Match) error {
	// Group matches by round
	roundMatches := make(map[string][]Match)
	for _, match := range matches {
		roundName := match.Round
		roundMatches[roundName] = append(roundMatches[roundName], match)
	}

	// Sort each round by start time
	for roundName := range roundMatches {
		sort.Slice(roundMatches[roundName], func(i, j int) bool {
			return roundMatches[roundName][i].StartTime.Before(roundMatches[roundName][j].StartTime)
		})
	}

	// Create nodes for each round
	for roundName, matches := range roundMatches {
		nodes := make([]*BracketNode, 0, len(matches))
		for i, match := range matches {
			node := &BracketNode{
				Match:    match,
				Position: i,
			}
			nodes = append(nodes, node)
		}
		bt.rounds[roundName] = nodes
	}

	return nil
}

// Render implements the Widget interface
func (bt *BracketTree) Render(ctx context.Context, buf *Buffer) error {
	if bt.compact {
		return bt.renderCompact(ctx, buf)
	}
	return bt.renderCompact(ctx, buf) // For now, only compact rendering
}

// MinSize implements the Widget interface
func (bt *BracketTree) MinSize() Size {
	height := 0
	maxWidth := 0

	// Calculate size based on rounds
	for roundName, matches := range bt.rounds {
		height += 1 + len(matches) + 1 // Header + matches + spacing
		roundWidth := len(roundName) + 10
		for _, node := range matches {
			matchWidth := bt.getMatchWidth(node)
			if matchWidth > roundWidth {
				roundWidth = matchWidth
			}
		}
		if roundWidth > maxWidth {
			maxWidth = roundWidth
		}
	}

	return Size{Width: maxWidth, Height: height}
}

// renderCompact renders the bracket in compact vertical list format
func (bt *BracketTree) renderCompact(ctx context.Context, buf *Buffer) error {
	// Get rounds in order: Final, Semi-finals, Quarter-finals, Round of 16
	roundOrder := []string{"Final", "Semi-finals", "Quarter-finals", "Round of 16"}

	for _, roundName := range roundOrder {
		matches, ok := bt.rounds[roundName]
		if !ok || len(matches) == 0 {
			continue
		}

		// Check context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Render round header
		buf.WriteStringWithStyle(roundName, Style{Foreground: ColorCyan, Bold: true})
		buf.NewLine()

		// Render matches
		for _, node := range matches {
			buf.WriteString("  ") // Indent
			buf.WriteString(bt.formatMatch(node))
			buf.NewLine()
		}

		// Add spacing between rounds
		buf.NewLine()
	}

	return nil
}

// formatMatch formats a match for display
func (bt *BracketTree) formatMatch(node *BracketNode) string {
	match := node.Match

	if match.Status == "Scheduled" || match.Status == "" {
		// Scheduled match
		return fmt.Sprintf("%s vs %s", match.HomeTeam, match.AwayTeam)
	}

	// Finished match
	result := fmt.Sprintf("%s %d - %d %s",
		match.HomeTeam, match.HomeScore, match.AwayScore, match.AwayTeam)

	// Add extra info if present
	if match.ExtraInfo != "" {
		result += fmt.Sprintf(" (%s)", match.ExtraInfo)
	}

	return result
}

// getMatchWidth calculates the display width of a match
func (bt *BracketTree) getMatchWidth(node *BracketNode) int {
	formatted := bt.formatMatch(node)
	// Simple length calculation (could use stringutil for emoji support)
	return len(formatted) + 2 // +2 for indent
}

// String returns a string representation
func (bt *BracketTree) String() string {
	totalMatches := 0
	for _, matches := range bt.rounds {
		totalMatches += len(matches)
	}
	return fmt.Sprintf("BracketTree{rounds=%d, matches=%d}", len(bt.rounds), totalMatches)
}
