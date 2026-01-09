package tui

import (
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// ColorLevel represents the color support level of the terminal
type ColorLevel int

const (
	// ColorNone means no color support
	ColorNone ColorLevel = iota
	// Color16 means basic ANSI color support (16 colors)
	Color16
	// Color256 means extended color support (256 colors)
	Color256
	// ColorTrueColor means RGB/true color support
	ColorTrueColor
)

// Terminal provides terminal capabilities
type Terminal struct {
	Width          int
	Height         int
	ColorSupport   ColorLevel
	UnicodeSupport bool
}

// Detect returns current terminal capabilities
func Detect() *Terminal {
	return &Terminal{
		Width:          detectWidth(),
		Height:         detectHeight(),
		ColorSupport:   detectColorLevel(),
		UnicodeSupport: detectUnicode(),
	}
}

// detectWidth detects terminal width
func detectWidth() int {
	// Try environment variable first
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	// Try syscall on Unix
	ws, err := getWinsize()
	if err == nil && ws.Col > 0 {
		return int(ws.Col)
	}

	// Default fallback
	return 80
}

// detectHeight detects terminal height
func detectHeight() int {
	// Try environment variable first
	if lines := os.Getenv("LINES"); lines != "" {
		if height, err := strconv.Atoi(lines); err == nil && height > 0 {
			return height
		}
	}

	// Try syscall on Unix
	ws, err := getWinsize()
	if err == nil && ws.Row > 0 {
		return int(ws.Row)
	}

	// Default fallback
	return 24
}

// detectColorLevel detects the color support level
func detectColorLevel() ColorLevel {
	// Check COLORTERM for true color support
	if colorTerm := os.Getenv("COLORTERM"); colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorTrueColor
	}

	// Check TERM for 256 color support
	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return Color256
	}

	// Check for basic color support
	if term != "" && term != "dumb" {
		return Color16
	}

	return ColorNone
}

// detectUnicode detects Unicode support
func detectUnicode() bool {
	// Check locale settings
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if val := os.Getenv(env); val != "" {
			val = strings.ToUpper(val)
			if strings.Contains(val, "UTF-8") || strings.Contains(val, "UTF8") {
				return true
			}
		}
	}

	// Default to true on modern systems
	return true
}

// winsize is used for terminal size detection on Unix
type winsize struct {
	Row uint16
	Col uint16
	X   uint16
	Y   uint16
}

// getWinsize gets the terminal window size on Unix systems
func getWinsize() (*winsize, error) {
	ws := &winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return nil, errno
	}
	return ws, nil
}
