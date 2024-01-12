/*
Copyright 2024 The cert-manager Authors.

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

package linter

import (
	"sort"
	"strings"
)

func DiffPaths(pathsA []string, pathsB []string) ([]string, []string) {
	sort.Strings(pathsA)
	sort.Strings(pathsB)

	missingA := []string{}
	missingB := []string{}

	prefix := "<NOT A PREFIX>"
	var i, j int
	for i < len(pathsA) && j < len(pathsB) {
		pathA := pathsA[i]
		pathB := pathsB[j]

		pathAHasPrefix := strings.HasPrefix(pathA, prefix)
		pathBHasPrefix := strings.HasPrefix(pathB, prefix)

		if pathA == pathB {
			prefix = pathA
			i++
			j++
			continue
		}

		if pathBHasPrefix && pathAHasPrefix {
			prefix = "<NOT A PREFIX>"
			continue
		}

		if pathA < pathB {
			if !pathAHasPrefix {
				missingB = append(missingB, pathA)
			}
			i++
		} else {
			if !pathBHasPrefix {
				missingA = append(missingA, pathB)
			}
			j++
		}
	}

	for i < len(pathsA) {
		pathA := pathsA[i]

		pathAHasPrefix := strings.HasPrefix(pathA, prefix)
		if !pathAHasPrefix {
			missingB = append(missingB, pathA)
		}

		i++
	}

	for j < len(pathsB) {
		pathB := pathsB[j]

		pathBHasPrefix := strings.HasPrefix(pathB, prefix)
		if !pathBHasPrefix {
			missingA = append(missingA, pathB)
		}

		j++
	}

	return missingA, missingB
}
