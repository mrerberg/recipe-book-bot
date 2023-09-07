package events

import (
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
