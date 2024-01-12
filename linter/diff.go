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
