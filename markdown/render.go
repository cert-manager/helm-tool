package markdown

import (
	"fmt"
	"strings"

	"github.com/cert-manager/helm-docgen/heuristics"
	"github.com/cert-manager/helm-docgen/parser"
)

func RenderDocument(document *parser.Document) string {
	var sb strings.Builder

	for _, section := range document.Sections {
		if len(section.Properties) == 0 {
			continue
		}

		header := renderSectionHeader(section)
		fmt.Fprint(&sb, header)

		fmt.Fprint(&sb, "\n<table>")
		fmt.Fprint(&sb, "\n<tr>")
		fmt.Fprint(&sb, "\n<th>Property</th>")
		fmt.Fprint(&sb, "\n<th>Description</th>")
		fmt.Fprint(&sb, "\n<th>Type</th>")
		fmt.Fprint(&sb, "\n<th>Default</th>")
		fmt.Fprint(&sb, "\n</tr>")

		for _, prop := range section.Properties {
			description := renderCommentAsMarkdown(prop.Description)
			defaultValue := renderStringAsYamlCodeBlock(prop.Default)

			fmt.Fprint(&sb, "\n<tr>")
			fmt.Fprintf(&sb, "\n<td>%s</td>", prop.Name)
			fmt.Fprintf(&sb, "\n<td>%s\n</td>", description)
			fmt.Fprintf(&sb, "\n<td>%s</td>", prop.Type)
			fmt.Fprintf(&sb, "\n<td>%s\n</td>", defaultValue)
			fmt.Fprint(&sb, "\n</tr>")
		}
		fmt.Fprint(&sb, "\n</table>\n")
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
			if trimmed := strings.TrimSpace(str); trimmed != "" {
				str = trimmed
				str = fmt.Sprintf("<p>%s</p>", strings.ReplaceAll(str, "\n\n", "</p><p>"))
				str = strings.ReplaceAll(str, "\n", "<br>\n")
				str = strings.ReplaceAll(str, "<p>", "\n<p>\n\n")
				str = strings.ReplaceAll(str, "</p>", "\n\n</p>")
			}
		default:
			continue
		}
		sb.WriteString(str)
	}

	return sb.String()
}

func renderStringAsYamlCodeBlock(str string) string {
	str = fmt.Sprintf("\n\n<pre lang=\"yaml\">%s</pre>\n", str)
	return str
}
