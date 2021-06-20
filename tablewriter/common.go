package tablewriter

import (
	"errors"

	"github.com/yudgnahk/euro21/constants"
)

var (
	InvalidRowDataLength = errors.New("invalid row data length")
)

type TableData struct {
	Data  [][]string
	Color []constants.Color
}

const (
	VerticalBlock   = "│"
	HorizontalBlock = "─"
)

type VerticalType string
type HorizontalType string

const (
	Top    HorizontalType = "top"
	Bottom HorizontalType = "bottom"
	Normal HorizontalType = "normal"

	Left   VerticalType = "left"
	Right  VerticalType = "right"
	Middle VerticalType = "middle"
)

var blocks = map[HorizontalType]map[VerticalType]string{
	Top: {
		Left:   "┌",
		Right:  "┐",
		Middle: "┬",
	},
	Bottom: {
		Left:   "└",
		Right:  "┘",
		Middle: "┴",
	},
	Normal: {
		Left:   "├",
		Right:  "┤",
		Middle: "┼",
	},
}
