package events

import (
	"math"
	"strings"
)

func IsCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "/")
}

func IsCallBack(cb string) bool {
	return strings.Contains(cb, "cb")
}

func ParseCallBack(cb string) []string {
	return strings.Split(cb, ":")
}

func GetPagesCount(recipesCount, recipesPerPage int64) int64 {
	return int64(math.Ceil(float64(recipesCount) / float64(recipesPerPage)))
}
