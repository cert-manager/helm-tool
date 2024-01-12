package linter

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/cert-manager/helm-docgen/linter/parsetemplates"
	"github.com/cert-manager/helm-docgen/parser"
)

func Lint(
	templatesFolder string,
	exceptionsPath string,
	document *parser.Document,
) error {
	templatePaths, err := parsetemplates.ListTemplatePaths(templatesFolder)
	if err != nil {
		return err
	}

	valuePathsDict := map[string]struct{}{}
	for _, section := range document.Sections {
		for _, property := range section.Properties {
			valuePathsDict[property.Path.String()] = struct{}{}
		}
	}
	valuePathsDict = parsetemplates.MakeUniform(valuePathsDict)
	valuePaths := []string{}
	for path := range valuePathsDict {
		valuePaths = append(valuePaths, path)
	}

	exceptionStrings := []string{}
	if exceptionsPath != "" {
		exceptionsPathsRaw, err := os.ReadFile(exceptionsPath)
		if err != nil {
			return err
		}

		exceptionStrings = strings.Split(string(exceptionsPathsRaw), "\n")
	}

	missingValues, missingTemplates := DiffPaths(valuePaths, templatePaths)

	succeeded := true
	for _, missingValue := range missingValues {
		exceptionString := fmt.Sprintf("value missing from values.yaml: %s", missingValue)

		if !slices.Contains(exceptionStrings, exceptionString) {
			fmt.Println(exceptionString)
			succeeded = false
		}
	}

	for _, missingTemplate := range missingTemplates {
		exceptionString := fmt.Sprintf("value missing from templates: %s", missingTemplate)

		if !slices.Contains(exceptionStrings, exceptionString) {
			fmt.Println(exceptionString)
			succeeded = false
		}
	}

	if !succeeded {
		return fmt.Errorf("values.yaml and templates are not in sync")
	}

	return nil
}
