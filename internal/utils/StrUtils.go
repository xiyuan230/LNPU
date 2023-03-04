package utils

import (
	"regexp"
	"strconv"
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

func CourseWeekListHandle(str string) []int {
	var result []int
	split := strings.Split(str, "(周)")[0]
	if ok, _ := regexp.MatchString("-", split); ok {
		weeks := strings.Split(split, ",")
		for _, week := range weeks {
			spl := strings.Split(week, "-")
			start, _ := strconv.Atoi(spl[0])
			end, _ := strconv.Atoi(spl[1])
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		}
		return result
	}
	week, _ := strconv.Atoi(split)
	return append(result, week)
}
