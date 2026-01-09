package tui

import (
	"bytes"
	"context"
	"testing"
)

func TestNewTable(t *testing.T) {
	table := NewTable()
	if table == nil {
		t.Fatal("NewTable() returned nil")
	}
	if len(table.headers) != 0 {
		t.Errorf("Expected 0 headers, got %d", len(table.headers))
	}
	if len(table.rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(table.rows))
	}
}

func TestTableSetHeaders(t *testing.T) {
	table := NewTable().SetHeaders("A", "B", "C")
	if len(table.headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(table.headers))
	}
	if table.headers[0] != "A" || table.headers[1] != "B" || table.headers[2] != "C" {
		t.Errorf("Headers not set correctly: %v", table.headers)
	}
}

func TestTableAddRow(t *testing.T) {
	table := NewTable().
		SetHeaders("Name", "Age").
		AddRow("Alice", "30").
		AddRow("Bob", "25")

	if len(table.rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(table.rows))
	}
}

func TestTableAddRowWithColor(t *testing.T) {
	table := NewTable().
		SetHeaders("Team", "Points").
		AddRowWithColor(ColorGreen, "Italy", "9").
		AddRowWithColor(ColorYellow, "Spain", "6")

	if len(table.rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(table.rows))
	}
	if len(table.colors) != 2 {
		t.Errorf("Expected 2 colors, got %d", len(table.colors))
	}
}

func TestTableRender(t *testing.T) {
	var buf bytes.Buffer
	renderer := NewRenderer(&buf)

	table := NewTable().
		SetHeaders("Name", "Score").
		AddRow("Alice", "100").
		AddRow("Bob", "95")

	ctx := context.Background()
	err := renderer.Render(ctx, table)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Check for table borders
	if !bytes.Contains(buf.Bytes(), []byte("┌")) {
		t.Error("Output missing top-left border")
	}
	if !bytes.Contains(buf.Bytes(), []byte("│")) {
		t.Error("Output missing vertical border")
	}

	// Check for headers
	if !bytes.Contains(buf.Bytes(), []byte("Name")) {
		t.Error("Output missing 'Name' header")
	}
	if !bytes.Contains(buf.Bytes(), []byte("Score")) {
		t.Error("Output missing 'Score' header")
	}

	// Check for data
	if !bytes.Contains(buf.Bytes(), []byte("Alice")) {
		t.Error("Output missing 'Alice' data")
	}
	if !bytes.Contains(buf.Bytes(), []byte("100")) {
		t.Error("Output missing '100' data")
	}
}

func TestTableMinSize(t *testing.T) {
	table := NewTable().
		SetHeaders("A", "B", "C").
		AddRow("1", "2", "3").
		AddRow("4", "5", "6")

	size := table.MinSize()
	if size.Width <= 0 {
		t.Errorf("Expected positive width, got %d", size.Width)
	}
	if size.Height <= 0 {
		t.Errorf("Expected positive height, got %d", size.Height)
	}
}

func TestMultiTable(t *testing.T) {
	table1 := NewTable().
		SetHeaders("Name").
		AddRow("Alice")

	table2 := NewTable().
		SetHeaders("Name").
		AddRow("Bob")

	multiTable := NewMultiTable().
		AddSection("Group A", table1).
		AddSection("Group B", table2)

	size := multiTable.MinSize()
	if size.Height <= 0 {
		t.Error("MultiTable MinSize() returned non-positive height")
	}

	var buf bytes.Buffer
	renderer := NewRenderer(&buf)
	ctx := context.Background()

	err := renderer.Render(ctx, multiTable)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Group A")) {
		t.Error("Output missing 'Group A' header")
	}
	if !bytes.Contains(buf.Bytes(), []byte("Group B")) {
		t.Error("Output missing 'Group B' header")
	}
}
