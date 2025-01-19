package stringUtil

import "strings"

func ReplaceSpaceWithUnderscore(str string) string {
	return strings.ReplaceAll(str, " ", "_")
}