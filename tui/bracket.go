package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/yudgnahk/euro21/utils/stringutil"
)

// BracketMatch represents a match in the tournament bracket
type BracketMatch struct {
	ID         int
	HomeTeam   string
	AwayTeam   string
	HomeFlag   string
	AwayFlag   string
	HomeScore  int
	AwayScore  int
	HomeWinner bool
	AwayWinner bool
	Round      string
}

// VisualBracket renders a 2-sided tournament bracket with ASCII art
type VisualBracket struct {
	leftMatches  []BracketMatch
	rightMatches []BracketMatch
	finalMatch   *BracketMatch
	theme        *Theme
	boxWidth     int
}

// NewVisualBracket creates a new visual bracket widget
func NewVisualBracket() *VisualBracket {
	return &VisualBracket{
		leftMatches:  make([]BracketMatch, 0),
		rightMatches: make([]BracketMatch, 0),
		theme:        DefaultTheme,
		boxWidth:     24,
	}
}

// SetTheme sets the visual theme
func (vb *VisualBracket) SetTheme(theme *Theme) *VisualBracket {
	vb.theme = theme
	return vb
}

// BuildEuro2021 creates the Euro 2021 knockout bracket structure
func (vb *VisualBracket) BuildEuro2021() *VisualBracket {
	// LEFT SIDE
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 1, HomeTeam: "Wales", HomeFlag: "🏴󠁧󠁢󠁷󠁬󠁳󠁿", AwayTeam: "Denmark", AwayFlag: "🇩🇰",
		HomeScore: 0, AwayScore: 4, HomeWinner: false, AwayWinner: true, Round: "R16",
	})
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 2, HomeTeam: "Italy", HomeFlag: "🇮🇹", AwayTeam: "Austria", AwayFlag: "🇦🇹",
		HomeScore: 2, AwayScore: 1, HomeWinner: true, AwayWinner: false, Round: "R16",
	})
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 3, HomeTeam: "Netherlands", HomeFlag: "🇳🇱", AwayTeam: "Czech Republic", AwayFlag: "🇨🇿",
		HomeScore: 0, AwayScore: 2, HomeWinner: false, AwayWinner: true, Round: "R16",
	})
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 4, HomeTeam: "Belgium", HomeFlag: "🇧🇪", AwayTeam: "Portugal", AwayFlag: "🇵🇹",
		HomeScore: 1, AwayScore: 0, HomeWinner: true, AwayWinner: false, Round: "R16",
	})
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 5, HomeTeam: "Denmark", HomeFlag: "🇩🇰", AwayTeam: "Czech Republic", AwayFlag: "🇨🇿",
		HomeScore: 2, AwayScore: 1, HomeWinner: true, AwayWinner: false, Round: "QF",
	})
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 6, HomeTeam: "Belgium", HomeFlag: "🇧🇪", AwayTeam: "Italy", AwayFlag: "🇮🇹",
		HomeScore: 1, AwayScore: 2, HomeWinner: false, AwayWinner: true, Round: "QF",
	})
	vb.leftMatches = append(vb.leftMatches, BracketMatch{
		ID: 7, HomeTeam: "Italy", HomeFlag: "🇮🇹", AwayTeam: "Spain", AwayFlag: "🇪🇸",
		HomeScore: 1, AwayScore: 1, HomeWinner: true, AwayWinner: false, Round: "SF",
	})

	// RIGHT SIDE
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 8, HomeTeam: "Croatia", HomeFlag: "🇭🇷", AwayTeam: "Spain", AwayFlag: "🇪🇸",
		HomeScore: 3, AwayScore: 5, HomeWinner: false, AwayWinner: true, Round: "R16",
	})
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 9, HomeTeam: "France", HomeFlag: "🇫🇷", AwayTeam: "Switzerland", AwayFlag: "🇨🇭",
		HomeScore: 3, AwayScore: 3, HomeWinner: false, AwayWinner: true, Round: "R16",
	})
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 10, HomeTeam: "England", HomeFlag: "🏴󠁧󠁢󠁥󠁮󠁧󠁿", AwayTeam: "Germany", AwayFlag: "🇩🇪",
		HomeScore: 2, AwayScore: 0, HomeWinner: true, AwayWinner: false, Round: "R16",
	})
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 11, HomeTeam: "Sweden", HomeFlag: "🇸🇪", AwayTeam: "Ukraine", AwayFlag: "🇺🇦",
		HomeScore: 1, AwayScore: 2, HomeWinner: false, AwayWinner: true, Round: "R16",
	})
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 12, HomeTeam: "Spain", HomeFlag: "🇪🇸", AwayTeam: "Switzerland", AwayFlag: "🇨🇭",
		HomeScore: 1, AwayScore: 1, HomeWinner: true, AwayWinner: false, Round: "QF",
	})
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 13, HomeTeam: "Ukraine", HomeFlag: "🇺🇦", AwayTeam: "England", AwayFlag: "🏴󠁧󠁢󠁥󠁮󠁧󠁿",
		HomeScore: 0, AwayScore: 4, HomeWinner: false, AwayWinner: true, Round: "QF",
	})
	vb.rightMatches = append(vb.rightMatches, BracketMatch{
		ID: 14, HomeTeam: "England", HomeFlag: "🏴󠁧󠁢󠁥󠁮󠁧󠁿", AwayTeam: "Denmark", AwayFlag: "🇩🇰",
		HomeScore: 2, AwayScore: 1, HomeWinner: true, AwayWinner: false, Round: "SF",
	})

	// FINAL
	vb.finalMatch = &BracketMatch{
		ID: 15, HomeTeam: "Italy", HomeFlag: "🇮🇹", AwayTeam: "England", AwayFlag: "🏴󠁧󠁢󠁥󠁮󠁧󠁿",
		HomeScore: 1, AwayScore: 1, HomeWinner: true, AwayWinner: false, Round: "Final",
	}

	return vb
}

// Render implements the Widget interface
func (vb *VisualBracket) Render(ctx context.Context, buf *Buffer) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	border := vb.theme.Border
	w := vb.boxWidth

	// Fixed column positions - with proper spacing for connectors
	colR16L := 0
	colQFL := 32   // R16(24) + 8 space for connector
	colSFL := 60   // QF(24) + 4 space
	colFinal := 86 // SF(24) + 4 space + center offset
	colSFR := 116  // mirror of SFL
	colQFR := 142  // mirror of QFL
	colR16R := 174 // mirror of R16L

	// Row positions
	r16Rows := []int{4, 8, 12, 16}
	qfRows := []int{6, 14}
	sfRow := 10
	finalRow := 10

	// Create grid
	grid := newStringGrid(22)

	// Title with styling
	title := "🏆  EURO 2021 - KNOCKOUT STAGE  🏆"
	grid.writeAt(0, 35, ColorYellow.ANSI()+title+ColorReset.ANSI())
	grid.writeAt(1, 0, ColorCyan.ANSI()+strings.Repeat("═", 110)+ColorReset.ANSI())

	// Headers
	grid.writeAt(2, colR16L, "Round of 16")
	grid.writeAt(2, colQFL+2, "Quarter-finals")
	grid.writeAt(2, colSFL+4, "Semi-finals")
	grid.writeAt(2, colFinal+6, "⚽ FINAL")
	grid.writeAt(2, colQFR+2, "Quarter-finals")
	grid.writeAt(2, colR16R, "Round of 16")

	// LEFT SIDE MATCHES (isRightSide = false)
	vb.renderMatchToGrid(grid, vb.leftMatches[0], colR16L, r16Rows[0], false)
	vb.renderMatchToGrid(grid, vb.leftMatches[1], colR16L, r16Rows[1], false)
	vb.renderMatchToGrid(grid, vb.leftMatches[2], colR16L, r16Rows[2], false)
	vb.renderMatchToGrid(grid, vb.leftMatches[3], colR16L, r16Rows[3], false)
	vb.renderMatchToGrid(grid, vb.leftMatches[4], colQFL, qfRows[0], false)
	vb.renderMatchToGrid(grid, vb.leftMatches[5], colQFL, qfRows[1], false)
	vb.renderMatchToGrid(grid, vb.leftMatches[6], colSFL, sfRow, false)

	// RIGHT SIDE MATCHES (isRightSide = true)
	vb.renderMatchToGrid(grid, vb.rightMatches[0], colR16R, r16Rows[0], true)
	vb.renderMatchToGrid(grid, vb.rightMatches[1], colR16R, r16Rows[1], true)
	vb.renderMatchToGrid(grid, vb.rightMatches[2], colR16R, r16Rows[2], true)
	vb.renderMatchToGrid(grid, vb.rightMatches[3], colR16R, r16Rows[3], true)
	vb.renderMatchToGrid(grid, vb.rightMatches[4], colQFR, qfRows[0], true)
	vb.renderMatchToGrid(grid, vb.rightMatches[5], colQFR, qfRows[1], true)
	vb.renderMatchToGrid(grid, vb.rightMatches[6], colSFR, sfRow, true)

	// FINAL (center - use left format, highlighted)
	vb.renderFinalToGrid(grid, *vb.finalMatch, colFinal, finalRow)

	// LEFT SIDE CONNECTORS
	// R16[0,1] -> QF[0] (rows 5,9 -> 7)
	vb.drawConnectorLR(grid, colR16L+w, colQFL, r16Rows[0]+1, r16Rows[1]+1, qfRows[0]+1, border)
	// R16[2,3] -> QF[1] (rows 13,17 -> 15)
	vb.drawConnectorLR(grid, colR16L+w, colQFL, r16Rows[2]+1, r16Rows[3]+1, qfRows[1]+1, border)
	// QF[0,1] -> SF (rows 7,15 -> 11)
	vb.drawConnectorLR(grid, colQFL+w, colSFL, qfRows[0]+1, qfRows[1]+1, sfRow+1, border)
	// SF -> Final
	for x := colSFL + w; x < colFinal; x++ {
		grid.writeAt(sfRow+1, x, border.Horizontal)
	}
	grid.writeAt(sfRow+1, colFinal-1, border.MiddleLeft)

	// RIGHT SIDE CONNECTORS
	// R16[0,1] -> QF[0]
	vb.drawConnectorRL(grid, colR16R-1, colQFR+w-1, r16Rows[0]+1, r16Rows[1]+1, qfRows[0]+1, border)
	// R16[2,3] -> QF[1]
	vb.drawConnectorRL(grid, colR16R-1, colQFR+w-1, r16Rows[2]+1, r16Rows[3]+1, qfRows[1]+1, border)
	// QF[0,1] -> SF
	vb.drawConnectorRL(grid, colQFR-1, colSFR+w-1, qfRows[0]+1, qfRows[1]+1, sfRow+1, border)
	// SF -> Final
	for x := colSFR - 1; x >= colFinal+w; x-- {
		grid.writeAt(sfRow+1, x, border.Horizontal)
	}
	grid.writeAt(sfRow+1, colFinal+w, border.MiddleRight)

	// Final cross
	grid.writeAt(sfRow+1, colFinal-1, border.Cross)
	grid.writeAt(sfRow+1, colFinal+w, border.Cross)

	// Champion banner - centered and styled
	bannerY := 19
	bannerWidth := 76
	bannerX := 17
	championText := "🏆  CHAMPION: 🇮🇹  ITALY  (1-1 vs England, won on penalties)"
	padding := (bannerWidth - 2 - stringutil.GetPrintableLength(championText)) / 2

	grid.writeAt(bannerY, bannerX, ColorGreen.ANSI()+"╔"+strings.Repeat("═", bannerWidth-2)+"╗"+ColorReset.ANSI())
	grid.writeAt(bannerY+1, bannerX, ColorGreen.ANSI()+"║"+strings.Repeat(" ", padding)+championText+strings.Repeat(" ", bannerWidth-2-padding-stringutil.GetPrintableLength(championText))+"║"+ColorReset.ANSI())
	grid.writeAt(bannerY+2, bannerX, ColorGreen.ANSI()+"╚"+strings.Repeat("═", bannerWidth-2)+"╝"+ColorReset.ANSI())

	// Output to buffer
	for _, row := range grid.getRows() {
		trimmed := strings.TrimRight(row, " ")
		buf.WriteString(trimmed)
		buf.NewLine()
	}

	return nil
}

// drawConnectorLR draws left-to-right connectors
func (vb *VisualBracket) drawConnectorLR(grid *stringGrid, fromX, toX, y1, y2, midY int, border BorderStyle) {
	// fromX = column just after right border of source box
	// toX = left border column of target box
	// y1, y2 = row positions of the two source match centers
	// midY = row position of the target match center

	// Calculate junction point (midway between fromX and toX)
	junctionX := fromX + 2

	// Draw horizontal lines from source to junction
	for x := fromX; x < junctionX; x++ {
		grid.writeAt(y1, x, border.Horizontal)
		grid.writeAt(y2, x, border.Horizontal)
	}

	// Place join characters at junction
	if y1 < midY {
		grid.writeAt(y1, junctionX, border.TopJoin)
	} else {
		grid.writeAt(y1, junctionX, border.BottomJoin)
	}
	if y2 < midY {
		grid.writeAt(y2, junctionX, border.TopJoin)
	} else {
		grid.writeAt(y2, junctionX, border.BottomJoin)
	}

	// Draw vertical line at junction (only between the two branches)
	startY, endY := y1, y2
	if startY > endY {
		startY, endY = endY, startY
	}
	for y := startY + 1; y < endY; y++ {
		if y == midY {
			grid.writeAt(y, junctionX, border.Cross)
		} else {
			grid.writeAt(y, junctionX, border.Vertical)
		}
	}

	// Draw horizontal line from junction to target (stop before target border)
	for x := junctionX + 1; x < toX; x++ {
		grid.writeAt(midY, x, border.Horizontal)
	}
}

// drawConnectorRL draws right-to-left connectors
func (vb *VisualBracket) drawConnectorRL(grid *stringGrid, fromX, toX, y1, y2, midY int, border BorderStyle) {
	// fromX = column just before left border of source box
	// toX = right border column of target box

	// Calculate junction point
	junctionX := fromX - 2

	// Draw horizontal lines from source to junction
	for x := fromX; x > junctionX; x-- {
		grid.writeAt(y1, x, border.Horizontal)
		grid.writeAt(y2, x, border.Horizontal)
	}

	// Place join characters
	if y1 < midY {
		grid.writeAt(y1, junctionX, border.TopJoin)
	} else {
		grid.writeAt(y1, junctionX, border.BottomJoin)
	}
	if y2 < midY {
		grid.writeAt(y2, junctionX, border.TopJoin)
	} else {
		grid.writeAt(y2, junctionX, border.BottomJoin)
	}

	// Draw vertical line
	startY, endY := y1, y2
	if startY > endY {
		startY, endY = endY, startY
	}
	for y := startY + 1; y < endY; y++ {
		if y == midY {
			grid.writeAt(y, junctionX, border.Cross)
		} else {
			grid.writeAt(y, junctionX, border.Vertical)
		}
	}

	// Draw horizontal line from junction to target
	for x := junctionX - 1; x > toX; x-- {
		grid.writeAt(midY, x, border.Horizontal)
	}
}

// stringGrid handles variable-width characters
type stringGrid struct {
	cells  map[int]map[int]string // row -> col -> string
	maxCol []int                  // max column written for each row
}

func newStringGrid(numRows int) *stringGrid {
	g := &stringGrid{
		cells:  make(map[int]map[int]string),
		maxCol: make([]int, numRows),
	}
	return g
}

func (g *stringGrid) writeAt(row, col int, s string) {
	if row < 0 || row >= len(g.maxCol) {
		return
	}

	if g.cells[row] == nil {
		g.cells[row] = make(map[int]string)
	}

	g.cells[row][col] = s
	if col > g.maxCol[row] {
		g.maxCol[row] = col
	}
}

// stripANSI removes ANSI escape codes from string
func stripANSI(s string) string {
	result := make([]rune, 0, len([]rune(s)))
	inEscape := false
	for _, r := range s {
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if r == '\033' || r == '\x1b' {
			inEscape = true
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func (g *stringGrid) getRows() []string {
	rows := make([]string, len(g.maxCol))
	for row := range rows {
		if g.cells[row] == nil {
			continue
		}

		// Build row by iterating through columns in order
		var sb strings.Builder
		currentCol := 0

		// Get sorted columns
		cols := make([]int, 0, len(g.cells[row]))
		for c := range g.cells[row] {
			cols = append(cols, c)
		}
		// Simple bubble sort for small number of columns
		for i := 0; i < len(cols); i++ {
			for j := i + 1; j < len(cols); j++ {
				if cols[i] > cols[j] {
					cols[i], cols[j] = cols[j], cols[i]
				}
			}
		}

		for _, col := range cols {
			s := g.cells[row][col]
			// Pad to reach this column
			for currentCol < col {
				sb.WriteByte(' ')
				currentCol++
			}
			sb.WriteString(s)
			// Use stripped length for width calculation
			currentCol += stringutil.GetPrintableLength(stripANSI(s))
		}
		rows[row] = sb.String()
	}
	return rows
}

// renderMatchToGrid renders a match box
func (vb *VisualBracket) renderMatchToGrid(grid *stringGrid, match BracketMatch, col, row int, isRightSide bool) {
	border := vb.theme.Border
	w := vb.boxWidth

	// Top border
	grid.writeAt(row, col, border.TopLeft+strings.Repeat(border.Horizontal, w-2)+border.TopRight)

	// Determine format based on side
	var homeLine, awayLine string
	if isRightSide {
		homeLine = vb.formatTeamLineRight(match.HomeFlag, match.HomeTeam, match.HomeScore)
		awayLine = vb.formatTeamLineRight(match.AwayFlag, match.AwayTeam, match.AwayScore)
	} else {
		homeLine = vb.formatTeamLineLeft(match.HomeFlag, match.HomeTeam, match.HomeScore)
		awayLine = vb.formatTeamLineLeft(match.AwayFlag, match.AwayTeam, match.AwayScore)
	}

	// Home team with winner highlighting
	if match.HomeWinner {
		homeLine = ColorGreen.ANSI() + homeLine + ColorReset.ANSI()
	}
	grid.writeAt(row+1, col, border.Vertical)
	grid.writeAt(row+1, col+1, homeLine)
	grid.writeAt(row+1, col+w-1, border.Vertical)

	// Away team with winner highlighting
	if match.AwayWinner {
		awayLine = ColorGreen.ANSI() + awayLine + ColorReset.ANSI()
	}
	grid.writeAt(row+2, col, border.BottomLeft)
	grid.writeAt(row+2, col+1, awayLine)
	grid.writeAt(row+2, col+w-1, border.BottomRight)
}

// renderFinalToGrid renders the final match with special styling
func (vb *VisualBracket) renderFinalToGrid(grid *stringGrid, match BracketMatch, col, row int) {
	border := vb.theme.Border
	w := vb.boxWidth

	// Top border with yellow styling
	grid.writeAt(row, col, ColorYellow.ANSI()+border.TopLeft+strings.Repeat(border.Horizontal, w-2)+border.TopRight+ColorReset.ANSI())

	// Format lines
	homeLine := vb.formatTeamLineLeft(match.HomeFlag, match.HomeTeam, match.HomeScore)
	awayLine := vb.formatTeamLineLeft(match.AwayFlag, match.AwayTeam, match.AwayScore)

	// Home team with winner highlighting
	if match.HomeWinner {
		homeLine = ColorGreen.ANSI() + homeLine + ColorReset.ANSI()
	}
	grid.writeAt(row+1, col, ColorYellow.ANSI()+border.Vertical+ColorReset.ANSI())
	grid.writeAt(row+1, col+1, homeLine)
	grid.writeAt(row+1, col+w-1, ColorYellow.ANSI()+border.Vertical+ColorReset.ANSI())

	// Away team with winner highlighting
	if match.AwayWinner {
		awayLine = ColorGreen.ANSI() + awayLine + ColorReset.ANSI()
	}
	grid.writeAt(row+2, col, ColorYellow.ANSI()+border.BottomLeft+ColorReset.ANSI())
	grid.writeAt(row+2, col+1, awayLine)
	grid.writeAt(row+2, col+w-1, ColorYellow.ANSI()+border.BottomRight+ColorReset.ANSI())
}

// formatTeamLineLeft formats a team line for left side (score on right)
func (vb *VisualBracket) formatTeamLineLeft(flag, team string, score int) string {
	scoreStr := fmt.Sprintf("%d", score)
	maxTeamLen := vb.boxWidth - 8 // flag(2) + space + team + padding + score(2)
	displayTeam := team
	if stringutil.GetPrintableLength(team) > maxTeamLen {
		displayTeam = team[:maxTeamLen-1] + "…"
	}
	// Format: "Flag Team        Score"
	return fmt.Sprintf("%s %s%s", flag, padRight(displayTeam, maxTeamLen-3), padLeft(scoreStr, 2))
}

// formatTeamLineRight formats a team line for right side (score on left)
func (vb *VisualBracket) formatTeamLineRight(flag, team string, score int) string {
	scoreStr := fmt.Sprintf("%d", score)
	maxTeamLen := vb.boxWidth - 8
	displayTeam := team
	if stringutil.GetPrintableLength(team) > maxTeamLen {
		displayTeam = team[:maxTeamLen-1] + "…"
	}
	// Format: "Score        Team Flag"
	return fmt.Sprintf("%s%s %s", padRight(scoreStr, 3), padRight(displayTeam, maxTeamLen-3), flag)
}

// MinSize implements the Widget interface
func (vb *VisualBracket) MinSize() Size {
	return Size{Width: 190, Height: 22}
}

// String returns a string representation
func (vb *VisualBracket) String() string {
	total := len(vb.leftMatches) + len(vb.rightMatches)
	if vb.finalMatch != nil {
		total++
	}
	return fmt.Sprintf("VisualBracket{matches=%d}", total)
}

func padRight(s string, width int) string {
	actualLen := stringutil.GetPrintableLength(s)
	if actualLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-actualLen)
}

func padLeft(s string, width int) string {
	actualLen := stringutil.GetPrintableLength(s)
	if actualLen >= width {
		return s
	}
	return strings.Repeat(" ", width-actualLen) + s
}
