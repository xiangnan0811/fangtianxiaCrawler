package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func ExtractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return strings.TrimSpace(string(match[1]))
	} else {
		return "暂无资料"
	}
}

func ExtractFloat64(contents []byte, re *regexp.Regexp) float64 {
	float64String := ExtractString(contents, re)
	result, err := strconv.ParseFloat(float64String, 64)
	if err == nil {
		return result
	} else {
		return 0
	}
}

func ExtractInt(contents []byte, re *regexp.Regexp) int {
	intString := ExtractString(contents, re)
	result, err := strconv.Atoi(intString)
	if err == nil {
		return result
	} else {
		return 0
	}
}

func ExtractAll(contents []byte, re *regexp.Regexp) []string {
	var result []string
	matches := re.FindAllSubmatch(contents, -1)
	if len(matches) >= 1 {
		for _, match := range matches {
			result = append(result, strings.TrimSpace(string(match[1])))
		}
	}
	return result
}
