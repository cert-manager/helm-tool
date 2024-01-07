package parser

import (
	"strings"
)

type tags map[string][]string

func (t *tags) Push(value string) {
	if *t == nil {
		*t = make(tags)
	}

	trimmed := strings.TrimSpace(value)
	key, value, _ := strings.Cut(trimmed[1:], "=")
	(*t)[key] = append((*t)[key], value)
}

func (t tags) GetBool(key string) bool {
	result := false

	for _, value := range t[key] {
		switch value {
		case "true", "enabled", "1", "yes", "":
			result = true
		case "false", "disabled", "0", "no":
			result = false
		default:
			result = true
		}
	}

	return result
}

func (t tags) GetString(key string) string {
	result := ""
	for _, value := range t[key] {
		result = value
	}

	return result
}
