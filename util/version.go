package util

import (
	"regexp"
	"strconv"
)

var pattern = regexp.MustCompile(`[^0-9]+`)

func CompareVersion(target string, global string) bool {
	targetVerStr := pattern.ReplaceAllString(target, "")
	globalVerStr := pattern.ReplaceAllString(global, "")
	targetInt, _ := strconv.Atoi(targetVerStr)
	globalInt, _ := strconv.Atoi(globalVerStr)
	return targetInt < globalInt
}
