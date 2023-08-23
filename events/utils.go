package events

import "strings"

func IsCommand(cmd string) bool {
	return strings.Contains(cmd, "/")
}
