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

package render

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/cert-manager/helm-tool/parser"

	"github.com/Masterminds/sprig/v3"
)

//go:embed markdown-plain
//go:embed markdown-table
//go:embed markdown-table-vertical
var templates embed.FS

func openTemplate(path string) (fs.File, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return templates.Open(path)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}

func Render(templateName string, document *parser.Document) (string, error) {
	tpl, err := openTemplate(templateName)
	if err != nil {
		return "", err
	}

	defer tpl.Close()

	templateBytes, err := io.ReadAll(tpl)
	if err != nil {
		return "", err
	}

	template, err := template.New(templateName).Funcs(sprig.TxtFuncMap()).Parse(string(templateBytes))
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err := template.Execute(&sb, document); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func Inject(path, templateName string, document *parser.Document, headerMatch, footerMatch *regexp.Regexp) error {
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
	startIdx := headerMatch.FindIndex(fileContents)
	if startIdx == nil {
		return errors.New("could not find parameters tag")
	}
	start := startIdx[1]

	// Find the end of where to inject
	endIdx := footerMatch.FindIndex(fileContents[start:])
	end := len(fileContents)
	if endIdx != nil {
		end = start + endIdx[0]
	}

	renderedDocument, err := Render(templateName, document)
	if err != nil {
		return errors.New("could not render documentation from template")
	}

	header := fileContents[:start]
	content := []byte(renderedDocument + "\n")
	footer := fileContents[end:]

	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(header)
	file.Write(content)
	file.Write(footer)

	return nil
}
