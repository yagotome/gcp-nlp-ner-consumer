package stringutil

import (
	"log"
	"regexp"
	"strings"
)

// RemovePunctuation ...
func RemovePunctuation(s string) string {
	reg, err := regexp.Compile(`[\!"#\$%&\\'\(\)\*\+,\-\./\:;\<\=\>\?@\[\\\]\^_` + "`" + `\{\|\}~]`)
	// reg, err := regexp.Compile(`["#\$%&\\'\(\)\*\+,\-/\:;\<\=\>@\[\\\]\^_` + "`" + `\{\|\}~]`)
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(strings.Replace(s, "\n", " ", -1), "")
}

// SplitByPunctuation ...
func SplitByPunctuation(s string) []string {
	r, err := regexp.Compile(`[^\.\!\?]+`)
	if err != nil {
		log.Fatal(err)
	}
	arr := r.FindAllString(s, -1)
	newArr := make([]string, len(arr))
	for i, e := range arr {
		newArr[i] = strings.Trim(e, " \t\n")
	}
	return newArr
}
