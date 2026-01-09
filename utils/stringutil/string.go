package stringutil

import (
	"strconv"

	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

func GetPrintableLength(s string) int {
	emojis := emoji.FindAll(s)

	// Start with the rune count (visual character count)
	// This handles Unicode characters like 'ü' correctly (1 rune = 1 visual char)
	l := len([]rune(s))

	if len(emojis) > 0 {
		for i := range emojis {
			// Calculate the rune count used by this emoji
			emojiRuneCount := getEmojiRuneCount(emojis[i].Locations)
			// Replace emoji rune count with visual width (2 for all emojis)
			l = l - emojiRuneCount + 2
		}
	}

	return l
}

// getEmojiRuneCount calculates the total rune count of an emoji
// from its location information
func getEmojiRuneCount(locations [][]int) int {
	// For subdivision flags like 🏴󠁧󠁢󠁥󠁮󠁧󠁿, locations is [[0 7]] meaning 7 runes
	// For regular flags like 🇭🇷, locations is [[0 2]] meaning 2 runes

	runeCount := 0
	for i := range locations {
		runeCount += locations[i][1] - locations[i][0]
	}

	return runeCount
}

func ToInt(s string) int {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return int(n)
}
