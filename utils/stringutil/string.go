package stringutil

import (
	"strconv"

	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

func GetPrintableLength(s string) int {
	emojis := emoji.FindAll(s)
	l := len(s)

	if len(emojis) > 0 {
		for i := range emojis {
			l = l - getLenFromSliceIndex(emojis[i].Locations) + 2
		}
	}

	return l
}

func getLenFromSliceIndex(indexSlice [][]int) int {
	res := 0

	for i := range indexSlice {
		res += indexSlice[i][1] - indexSlice[i][0]
	}

	if res > 2 {
		res = res * 4
	} else {
		res = res*4 + 1
	}

	return res
}

func ToInt(s string) int {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return int(n)
}
