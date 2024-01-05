package parser

import (
	"strings"
)

type Tags map[string][]string

func (t *Tags) Push(value string) {
	if *t == nil {
		*t = make(Tags)
	}

	trimmed := strings.TrimSpace(value)
	key, value, _ := strings.Cut(trimmed[1:], "=")
	(*t)[key] = append((*t)[key], value)
}

func (t Tags) GetBool(key string) bool {
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

func (t Tags) GetString(key string) string {
	result := ""
	for _, value := range t[key] {
		result = value
	}

	return result
}
