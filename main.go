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

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/cert-manager/helm-tool/linter"
	"github.com/cert-manager/helm-tool/parser"
	"github.com/cert-manager/helm-tool/render"
	"github.com/cert-manager/helm-tool/schema"
	"github.com/spf13/cobra"
)

var (
	valuesFile      string
	templatesFolder string
	exceptionsFile  string
	targetFile      string
	templateName    string
	headerSearch    = regexValue{regexp.MustCompile(`(?m)^##\s+Parameters *$`)}
	footerSearch    = regexValue{regexp.MustCompile(`(?m)^##?\s+.*$`)}
)

var Cmd = cobra.Command{
	Use: "helm-tool",
}

var Render = cobra.Command{
	Use:   "render",
	Short: "render documentation to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		document, err := parser.Load(valuesFile, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %q: %s\n", valuesFile, err)
			os.Exit(1)
		}

		result, err := render.Render(templateName, document)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering template: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

var Inject = cobra.Command{
	Use:   "inject",
	Short: "generate documentation and inject into existing markdown file",
	Run: func(cmd *cobra.Command, args []string) {
		document, err := parser.Load(valuesFile, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %q: %s\n", valuesFile, err)
			os.Exit(1)
		}

		if err := render.Inject(targetFile, templateName, document, headerSearch.regexp, footerSearch.regexp); err != nil {
			fmt.Fprintf(os.Stderr, "Could inject markdown into %q: %s\n", targetFile, err)
			os.Exit(1)
		}
	},
}

var Schema = cobra.Command{
	Use: "schema",
	Run: func(cmd *cobra.Command, args []string) {
		document, err := parser.Load(valuesFile, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %q: %s\n", valuesFile, err)
			os.Exit(1)
		}

		renderedSchema, err := schema.Render(document)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not render schema: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(renderedSchema)
	},
}

var Lint = cobra.Command{
	Use: "lint",
	Run: func(cmd *cobra.Command, args []string) {
		document, err := parser.Load(valuesFile, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %q: %s\n", valuesFile, err)
			os.Exit(1)
		}

		err = linter.Lint(templatesFolder, exceptionsFile, document)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not lint: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("No errors found")
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&valuesFile, "values", "i", "values.yaml", "values file used to generate the documentation")

	Cmd.AddCommand(&Inject)
	Inject.PersistentFlags().StringVarP(&templateName, "template", "t", "markdown-plain", "template to render documentation with")
	Inject.PersistentFlags().StringVarP(&targetFile, "output", "o", "README.md", "file to inject the generated markdown into")
	Inject.PersistentFlags().Var(&headerSearch, "header-search", "set the regex used to match the start of the injected markdown")
	Inject.PersistentFlags().Var(&footerSearch, "footer-search", "set the regex used to match the end of the injected markdown")

	Cmd.AddCommand(&Render)
	Render.PersistentFlags().StringVarP(&templateName, "template", "t", "markdown-plain", "template to render documentation with")

	Cmd.AddCommand(&Schema)

	Cmd.AddCommand(&Lint)
	Lint.PersistentFlags().StringVarP(&templatesFolder, "templates", "d", "templates", "templates folder used to lint the values file")
	Lint.PersistentFlags().StringVarP(&exceptionsFile, "exceptions", "e", "", "file containing exceptions to the linting rules")
}

func main() {
	Cmd.Execute()
}

type regexValue struct {
	regexp *regexp.Regexp
}

func (r *regexValue) String() string {
	return r.regexp.String()
}

func (r *regexValue) Set(value string) error {
	compiled, err := regexp.Compile("(?m)" + value)
	if err != nil {
		return err
	}

	r.regexp = compiled
	return nil
}

func (r *regexValue) Type() string {
	return "regex"
}
