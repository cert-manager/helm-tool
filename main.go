package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/cert-manager/helm-docgen/linter"
	"github.com/cert-manager/helm-docgen/parser"
	"github.com/cert-manager/helm-docgen/render"
	"github.com/cert-manager/helm-docgen/schema"
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
	Use: "helm-docgen",
}

var Render = cobra.Command{
	Use:   "render",
	Short: "render documentation to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		document, err := parser.Load(valuesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %q: %s\n", valuesFile, err)
			os.Exit(1)
		}

		fmt.Println(render.Render(templateName, document))
	},
}

var Inject = cobra.Command{
	Use:   "inject",
	Short: "generate documentation and inject into existing markdown file",
	Run: func(cmd *cobra.Command, args []string) {
		document, err := parser.Load(valuesFile)
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
		document, err := parser.Load(valuesFile)
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
		document, err := parser.Load(valuesFile)
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
	Inject.PersistentFlags().StringVarP(&templateName, "template", "t", "markdown-table", "template to render documentation with")
	Inject.PersistentFlags().StringVarP(&targetFile, "output", "o", "README.md", "file to inject the generated markdown into")
	Inject.PersistentFlags().Var(&headerSearch, "header-search", "set the regex used to match the start of the injected markdown")
	Inject.PersistentFlags().Var(&footerSearch, "footer-search", "set the regex used to match the end of the injected markdown")

	Cmd.AddCommand(&Render)
	Render.PersistentFlags().StringVarP(&templateName, "template", "t", "markdown-table", "template to render documentation with")

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
