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
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/cert-manager/helm-tool/linter/parsetemplates"
	"github.com/cert-manager/helm-tool/linter/sets"
	"github.com/cert-manager/helm-tool/parser"
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

	valuePaths := sets.Set[string]{}
	for _, section := range document.Sections {
		for _, property := range section.Properties {
			valuePaths.Insert(property.Path.PatternString())
		}
	}
	valuePaths = sets.RemovePrefixes(valuePaths)

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
	for missingValue := range missingValues {
		exceptionString := fmt.Sprintf("value missing from values.yaml: %s", missingValue)

		if !slices.Contains(exceptionStrings, exceptionString) {
			fmt.Println(exceptionString)
			succeeded = false
		}
	}

	for missingTemplate := range missingTemplates {
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
