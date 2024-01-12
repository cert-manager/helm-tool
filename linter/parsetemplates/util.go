package parsetemplates

import (
	"fmt"
	"strings"
)

func MakeUniform(paths map[string]struct{}) map[string]struct{} {
	results := map[string]struct{}{}

	for path := range paths {
		sections := strings.Split(path, ".")
		buildPath := ""
	SectionLoop:
		for _, section := range sections {
			if section == "" {
				continue SectionLoop
			}

			buildPath = fmt.Sprintf("%s.%s", buildPath, section)
			results[buildPath] = struct{}{}
		}
	}

	return results
}
