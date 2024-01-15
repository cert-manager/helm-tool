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

package heuristics

import (
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type ContentType string

const (
	ContentTypeUnknown ContentType = ""
	ContentTypeText    ContentType = "text"
	ContentTypeYaml    ContentType = "yaml"
	ContentTypeTag     ContentType = "tag"
)

// ContentSniffer is used to parse lines of text to determine the content, this
// is used to try and format the comments in a values.yaml file in a sensible
// way, by formatting yaml blocks correctly for example.
type ContentSniffer struct {
	previousLineBuffer []string
	leadingSpaces      int
	currentType        ContentType
}

func (c *ContentSniffer) SniffContentType(line string) (ContentType, bool) {
	switch c.currentType {
	case ContentTypeUnknown, ContentTypeText, ContentTypeTag:
		return c.sniffBasic(line)
	case ContentTypeYaml:
		return c.sniffYamlContinuation(line)
	default:
		panic("unreachable")
	}
}

func (c *ContentSniffer) sniffBasic(line string) (ContentType, bool) {
	previousType := c.currentType

	switch {
	case isLineTag(line):
		c.previousLineBuffer = nil
		c.currentType = ContentTypeTag
		return ContentTypeTag, true
	case isLineYamlRestrictive(line):
		c.previousLineBuffer = []string{line}
		c.currentType = ContentTypeYaml
		c.leadingSpaces = countLeadingSpaces(line)
		return ContentTypeYaml, true
	default:
		c.previousLineBuffer = append(c.previousLineBuffer, line)
		c.currentType = ContentTypeText
		return ContentTypeText, previousType != ContentTypeText
	}
}

func (c *ContentSniffer) sniffYamlContinuation(line string) (ContentType, bool) {
	candidateCodeLines := append(c.previousLineBuffer, line)
	candidateCode := strings.Join(candidateCodeLines, "\n")

	// If we are less indented then the first line of yaml, start sniffing again
	if countLeadingSpaces(line) < c.leadingSpaces {
		return c.sniffBasic(line)
	}

	var node yaml.Node
	if yaml.Unmarshal([]byte(candidateCode), &node) != nil {
		return c.sniffBasic(line)
	}

	c.previousLineBuffer = candidateCodeLines
	return ContentTypeYaml, false
}

func isLineTag(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}

	return strings.HasPrefix(trimmed, "+docs:")
}

// isLineYamlRestrictive determine if the line is yaml(ish). It parses the line as
// yaml and returns true only if the following criteria is met:
//   - It is a yaml map
//   - The map key has no spaces
//   - The map key starts with a lowercase letter
//   - If the map value contains spaces, the value must be quoted
//   - If the key is "ref", and the value is a URL then skip it (it's likely a
//     comment referencing something)
func isLineYamlRestrictive(line string) bool {
	var node yaml.Node
	if yaml.Unmarshal([]byte(line), &node) != nil {
		return false
	}

	return isNodeYamlRestrictive(&node)
}

func isNodeYamlRestrictive(node *yaml.Node) bool {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, node := range node.Content {
			return isNodeYamlRestrictive(node)
		}
	case yaml.MappingNode:
		if len(node.Content) != 2 {
			return false
		}

		keyNode := node.Content[0]
		valueNode := node.Content[1]

		if strings.Contains(keyNode.Value, " ") {
			return false
		}

		if len(keyNode.Value) == 0 || unicode.IsUpper(rune(keyNode.Value[0])) {
			return false
		}

		if strings.Contains(valueNode.Value, " ") && valueNode.Style != yaml.DoubleQuotedStyle && valueNode.Style != yaml.SingleQuotedStyle {
			return false
		}

		if keyNode.Value == "ref" && strings.HasPrefix(valueNode.Value, "http") {
			return false
		}

		return true
	}
	return false
}

type CommentBlock struct {
	Segments []CommentBlockSegment
}

func (c *CommentBlock) String() string {
	var sb strings.Builder
	for _, segment := range c.Segments {
		if segment.Type == ContentTypeTag {
			continue
		}
		sb.WriteString(segment.String())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

type CommentBlockSegment struct {
	Type     ContentType
	Contents []string
}

func (c *CommentBlockSegment) String() string {
	switch c.Type {
	case ContentTypeTag:
		return c.Contents[0]
	case ContentTypeYaml:
		return strings.Join(trimLeadingSpaces(c.Contents), "\n")
	case ContentTypeText:
		return strings.Join(RecutNewLines(c.Contents), "\n")
	default:
		panic("unreachable")
	}
}

func ParseCommentIntoBlocks(comment string) []CommentBlock {
	var sniffer ContentSniffer
	var parsedBlocks []CommentBlock
	var currentBlock CommentBlock
	var currentSegment CommentBlockSegment

	for _, line := range strings.Split(comment, "\n") {
		// Get the line without leading and training spaces, this means the
		// first character should be '#'
		trimmedLine := strings.TrimSpace(line)

		// Empty lines are used to break up blocks
		//
		// For example:
		//
		// # This is a comment block
		//
		// # This is a different comment block
		//
		if trimmedLine == "" {
			// Append any last segments before we are done with the block
			if currentSegment.Type != ContentTypeUnknown {
				completeSegment(&currentSegment)
				currentBlock.Segments = append(currentBlock.Segments, currentSegment)
				currentSegment = CommentBlockSegment{}
			}

			// If we have segments in the current block add them
			if len(currentBlock.Segments) > 0 {
				parsedBlocks = append(parsedBlocks, currentBlock)
				currentBlock = CommentBlock{}
			}

			continue
		}

		// Get the line without the leading comment characters
		trimmedLineWithoutCommentCharacter := strings.TrimLeft(trimmedLine, "#")

		// Sniffer tells us when to break up blocks and the type of those blocks
		typ, isNewBlock := sniffer.SniffContentType(trimmedLineWithoutCommentCharacter)
		if isNewBlock && currentSegment.Type != ContentTypeUnknown {
			completeSegment(&currentSegment)
			currentBlock.Segments = append(currentBlock.Segments, currentSegment)
			currentSegment = CommentBlockSegment{}
		}

		currentSegment.Type = typ
		currentSegment.Contents = append(currentSegment.Contents, trimmedLineWithoutCommentCharacter)
	}

	// Ensure any last segments get appended to the block
	if currentSegment.Type != ContentTypeUnknown {
		completeSegment(&currentSegment)
		currentBlock.Segments = append(currentBlock.Segments, currentSegment)
	}

	// We always append the currentBlock to the parsedBlocks even if it's empty.
	// This is because we don't want a comment block above a value, but
	// separated by a new line, to be attributed to the value.
	//
	// For example:
	//
	//   # This comment should not be attributed to the below value as there is
	//   # a new line separating them
	//
	//   a: b
	//
	return append(parsedBlocks, currentBlock)
}

func completeSegment(segment *CommentBlockSegment) {
	switch segment.Type {
	case ContentTypeText:
		segment.Contents = RecutNewLines(segment.Contents)
	case ContentTypeYaml:
		segment.Contents = trimLeadingSpaces(segment.Contents)
	case ContentTypeTag:
		segment.Contents = trimLeadingSpaces(segment.Contents)
	}
}
