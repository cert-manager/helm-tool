package markdown

import (
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/thatsmrtalbot/helm-docgen/parser"
)

var headerSearch = regexp.MustCompile(`(?m)^##\s+Parameters\s*$`)
var footerSearch = regexp.MustCompile(`(?m)^##?\s+.*$`)

func InjectIntoMarkdown(path string, document *parser.Document) error {
	// Open the file
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	// Read the contents
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Find the start of where to inject
	startIdx := headerSearch.FindIndex(fileContents)
	if startIdx == nil {
		return errors.New("could not find parameters tag")
	}
	start := startIdx[1]

	// Find the end of where to inject
	endIdx := footerSearch.FindIndex(fileContents[start:])
	end := len(fileContents)
	if endIdx != nil {
		end = start + endIdx[0]
	}

	header := fileContents[:start]
	content := []byte(RenderDocument(document) + "\n")
	footer := fileContents[end:]

	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(header)
	file.Write(content)
	file.Write(footer)

	return nil
}
