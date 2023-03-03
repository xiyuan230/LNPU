package utils

import (
	"strings"
)

func TrimSpace(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}

func ScoreStrHandle(str string) string {
	split := strings.Split(str, ":")
	if len(split) == 1 {
		return ""
	}
	return split[1]
}
