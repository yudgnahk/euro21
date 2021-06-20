package tablewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/yudgnahk/euro21/constants"
	"github.com/yudgnahk/euro21/utils/stringutil"
)

type Table struct {
	out     io.Writer
	headers []string
	data    TableData
	colLens []int
}

func NewTable(writer io.Writer) *Table {
	return &Table{
		out:     writer,
		headers: []string{},
		data:    TableData{},
		colLens: []int{},
	}
}

func (t *Table) SetHeader(headers []string) {
	for i := range headers {
		headers[i] = fmt.Sprintf(" %v ", strings.TrimSpace(headers[i]))
	}

	t.headers = headers

	for i := range headers {
		t.colLens = append(t.colLens, len(t.headers[i])+1)
	}
}

func (t *Table) Append(row []string) error {
	if len(row) != len(t.headers) {
		return InvalidRowDataLength
	}

	for i := range row {
		row[i] = fmt.Sprintf(" %v ", strings.TrimSpace(row[i]))
	}

	t.data.Data = append(t.data.Data, row)

	for i := range row {
		printableLength := stringutil.GetPrintableLength(row[i])
		if t.colLens[i] < printableLength+1 {
			t.colLens[i] = printableLength + 1
		}
	}

	return nil
}

func render(s string, l int, color constants.Color) {
	sLen := stringutil.GetPrintableLength(s) + 1

	for i := sLen; i < l; i++ {
		s += " "
	}

	fmt.Print(VerticalBlock, color, s, constants.ColorWhite)
}

func renderHeaders(headers []string, colLens []int) {
	for i := range headers {
		render(headers[i], colLens[i], constants.ColorWhite)
	}
	fmt.Println(VerticalBlock)
	renderBorder(Normal, colLens)
}

func renderRow(s []string, l []int, color constants.Color) {
	for i := range s {
		render(s[i], l[i], color)
	}

	fmt.Println(VerticalBlock, constants.ColorWhite)
}

func renderRows(data TableData, colLens []int) {
	for i := range data.Data {
		renderRow(data.Data[i], colLens, data.Color[i])
		if i != len(data.Data)-1 {
			renderBorder(Normal, colLens)
		}
	}
}

func renderBorder(horType HorizontalType, colLens []int) {
	border := blocks[horType][Left]
	start := 1
	for i := range colLens {
		for j := start; j < start+colLens[i]-1; j++ {
			border += HorizontalBlock
		}
		if i != len(colLens)-1 {
			border += blocks[horType][Middle]
		}

		start = start + colLens[i]
	}

	border += blocks[horType][Right]
	fmt.Println(border)
}

func (t *Table) Render() {
	renderBorder(Top, t.colLens)

	renderHeaders(t.headers, t.colLens)
	renderRows(t.data, t.colLens)

	renderBorder(Bottom, t.colLens)
}
