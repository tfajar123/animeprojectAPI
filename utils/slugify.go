package utils

import (
	"regexp"
	"strings"
)

func Slugify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	reg, _ := regexp.Compile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	reg2, _ := regexp.Compile("-+")
	s = reg2.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}