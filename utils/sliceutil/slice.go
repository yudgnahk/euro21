package sliceutil

import "fmt"

func ToStringSlice(args ...interface{}) []string {
	res := make([]string, len(args))
	for i := range args {
		res[i] = fmt.Sprintf("%v", args[i])
	}

	return res
}
