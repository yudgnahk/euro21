package tablewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/yudgnahk/euro21/utils/stringutil"
)

type MultiTables struct {
	out        io.Writer
	headers    []string
	subHeaders []string
	data       []TableData
	colLens    []int
}

func NewMultiTables(writer io.Writer) *MultiTables {
	return &MultiTables{
		out:        writer,
		headers:    []string{},
		subHeaders: []string{},
		data:       []TableData{},
		colLens:    []int{},
	}
}

func (t *MultiTables) SetHeaders(headers []string) {
	for i := range headers {
		headers[i] = fmt.Sprintf(" %v ", strings.TrimSpace(headers[i]))
	}

	t.headers = headers

	for i := range headers {
		t.colLens = append(t.colLens, len(t.headers[i])+1)
	}
}

func (t *MultiTables) AppendSubHeaders(subHeader string) {
	t.subHeaders = append(t.subHeaders, subHeader)
}

func (t *MultiTables) AppendTable(tableData TableData) error {
	var data TableData

	for i := range tableData.Data {
		row := tableData.Data[i]

		if len(row) != len(t.headers) {
			return InvalidRowDataLength
		}

		for i := range row {
			row[i] = fmt.Sprintf(" %v ", strings.TrimSpace(row[i]))
		}

		data.Data = append(data.Data, row)
		data.Color = tableData.Color

		for i := range row {
			printableLength := stringutil.GetPrintableLength(row[i])
			if t.colLens[i] < printableLength+1 {
				t.colLens[i] = printableLength + 1
			}
		}
	}

	t.data = append(t.data, data)

	return nil
}

func (t *MultiTables) Render() {
	for i := range t.data {
		fmt.Println(t.subHeaders[i])

		renderBorder(Top, t.colLens)

		renderHeaders(t.headers, t.colLens)
		renderRows(t.data[i], t.colLens)

		renderBorder(Bottom, t.colLens)
	}
}
