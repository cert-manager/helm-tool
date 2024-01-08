package markdown

import (
	"fmt"
	"strings"

	"github.com/thatsmrtalbot/helm-docgen/heuristics"
	"github.com/thatsmrtalbot/helm-docgen/parser"
)

func RenderDocument(document *parser.Document) string {
	var sb strings.Builder

	for _, section := range document.Sections {
		if len(section.Properties) == 0 {
			continue
		}

		header := renderSectionHeader(section)
		fmt.Fprint(&sb, header)
		fmt.Fprint(&sb, "\n|property|description|type|default|\n")
		fmt.Fprint(&sb, "|--|--|--|--|\n")

		for _, prop := range section.Properties {
			description := renderCommentAsMarkdown(prop.Description)
			defaultValue := renderStringAsYamlCodeBlock(prop.Default)
			fmt.Fprintf(&sb, "|`%s`|%s|`%s`|%s|\n", prop.Name, description, prop.Type, defaultValue)
		}
	}

	return sb.String()
}

func renderSectionHeader(section parser.Section) string {
	if len(section.Properties) == 0 || section.Name == "" {
		return ""
	}

	description := renderCommentAsMarkdown(section.Description)
	if description == "" {
		return fmt.Sprintf("\n### %s\n", section.Name)
	}

	return fmt.Sprintf("\n### %s\n%s\n", section.Name, description)
}

func renderCommentAsMarkdown(comment parser.Comment) string {
	var sb strings.Builder

	for _, section := range comment.Segments {
		str := section.String()

		switch section.Type {
		case heuristics.ContentTypeYaml:
			str = renderStringAsYamlCodeBlock(str)
		case heuristics.ContentTypeText:
			if strings.TrimSpace(str) != "" {
				str = fmt.Sprintf("<p>%s</p>", strings.ReplaceAll(str, "\n\n", "</p><p>"))
			}
		default:
			continue
		}

		str = strings.ReplaceAll(str, "\n", "<br>")
		sb.WriteString(str)
	}

	return sb.String()
}

func renderStringAsYamlCodeBlock(str string) string {
	str = fmt.Sprintf(`<pre lang="yaml">%s</pre>`, str)
	str = strings.ReplaceAll(str, "\n", "<br>")
	return str
}
