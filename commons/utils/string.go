package utils

import (
	"fmt"
	"strings"
)

func Join[T ~string](slice []T, delimeter, surround string) string {
	stringSlice := make([]string, 0)
	for _, x := range slice {
		x := fmt.Sprintf("%[2]s%[1]s%[2]s", string(x), surround)
		stringSlice = append(stringSlice, x)
	}

	return strings.Join(stringSlice, delimeter)
}
