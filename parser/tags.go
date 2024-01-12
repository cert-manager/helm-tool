/*
Copyright 2021 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
