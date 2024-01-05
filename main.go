package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thatsmrtalbot/helm-docgen/internal/markdown"
	"github.com/thatsmrtalbot/helm-docgen/internal/parser"
)

var (
	valuesFile string
	targetFile string
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

		fmt.Println(markdown.RenderDocument(document))
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

		if err := markdown.InjectIntoMarkdown(targetFile, document); err != nil {
			fmt.Fprintf(os.Stderr, "Could inject markdown into %q: %s\n", targetFile, err)
			os.Exit(1)
		}
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&valuesFile, "values", "i", "values.yaml", "values file used to generate the documentation")

	Cmd.AddCommand(&Inject)
	Inject.PersistentFlags().StringVarP(&targetFile, "output", "o", "README.md", "file to inject the generated markdown into")

	Cmd.AddCommand(&Render)
}

func main() {
	Cmd.Execute()
}
