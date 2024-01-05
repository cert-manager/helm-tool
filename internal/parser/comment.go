package parser

import (
	"bufio"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type Comment struct {
	Sections []CommentSection
	Tags     Tags
}

type CommentSection struct {
	Type    CommentType
	Content []string
}

func (s CommentSection) String() string {
	if s.Type == CommentTypeCode {
		return printCodeTrimLeft(s.Content)
	}

	if s.Type != CommentTypeText {
		return strings.Join(s.Content, "\n")
	}

	// Comments use new lines to word wrap, we need to destinguish betweewn
	// desired and undesired new lines, we can try and do this by looking
	// at the line length
	const newLineThreshold = 100

	var sb strings.Builder

	nl := 0
	for _, line := range s.Content {
		// Preserve empty lines
		if len(line) != 0 {
			sb.WriteString(line)
			nl = 0
		} else {
			for i := nl; i < 2; i++ {
				sb.WriteRune('\n')
			}
			nl = 2
		}

		// If the line was shorter then the threshold, write a new
		// line
		if len(line) < newLineThreshold || countLeadingSpaces(line) > 2 {
			for i := nl; i < 1; i++ {
				sb.WriteRune('\n')
			}
			nl = 1
		}
	}

	return strings.TrimSpace(sb.String())
}

type CommentType uint8

const (
	CommentTypeUnknown CommentType = iota
	CommentTypeText
	CommentTypeCode
	CommentTypeTag
)

func (c *Comment) Empty() bool {
	return len(c.Sections) == 0
}

func (c *Comment) String() string {
	var sb strings.Builder

	for _, section := range c.Sections {
		if section.Type != CommentTypeTag {
			sb.WriteString(section.String())
		}
	}

	return strings.TrimSpace(sb.String())
}

type Comments []Comment

func (c *Comments) Pop() (Comment, bool) {
	l := len(*c)
	if l == 0 {
		return Comment{}, false
	}

	comment := (*c)[l-1]
	*c = (*c)[:l-1]
	return comment, true
}

func ParseComment(comment string) Comments {
	reader := strings.NewReader(comment)
	scanner := bufio.NewScanner(reader)

	comments := []Comment{}
	currentComment := Comment{}
	currentSection := CommentSection{}

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// If there is a line break, then the previous section has ended
		if trimmedLine == "" {
			if currentSection.Type != CommentTypeUnknown {
				currentComment.Sections = append(currentComment.Sections, currentSection)
				currentSection = CommentSection{}
			}
			if !currentComment.Empty() {
				comments = append(comments, currentComment)
				currentComment = Comment{}
			}
			continue
		}

		// Remove leading '#'
		trimmedLine = trimmedLine[1:]

		for {
			// If we are parsing a code block then we append this new line to
			// the existing lines to see if the yaml is still valid. If it is
			// we append to the code block, if not then the code block has
			// ended.
			if currentSection.Type == CommentTypeCode {
				content := append(currentSection.Content, trimmedLine)
				body := strings.Join(content, "\n")

				if isValidYaml(body) {
					currentSection.Content = content
					break
				} else {
					currentComment.Sections = append(currentComment.Sections, currentSection)
					currentSection = CommentSection{}
					continue
				}
			}

			// If this line is part of a code block then start parsing it as code
			if isValidNonScalarYaml(trimmedLine) {
				if currentSection.Type != CommentTypeUnknown {
					currentComment.Sections = append(currentComment.Sections, currentSection)
				}

				currentSection = CommentSection{Type: CommentTypeCode, Content: []string{trimmedLine}}
				break
			}

			// If the line is a tag, push it as a section on its own
			if isTag(trimmedLine) {
				if currentSection.Type != CommentTypeUnknown {
					currentComment.Sections = append(currentComment.Sections, currentSection)
					currentSection = CommentSection{}
				}

				currentComment.Sections = append(currentComment.Sections, CommentSection{
					Type:    CommentTypeTag,
					Content: []string{trimmedLine},
				})

				currentComment.Tags.Push(trimmedLine)
				break
			}

			if currentSection.Type == CommentTypeUnknown {
				currentSection.Type = CommentTypeText
			}

			if currentSection.Type != CommentTypeText {
				currentComment.Sections = append(currentComment.Sections, currentSection)
				currentSection = CommentSection{Type: CommentTypeText}
			}

			currentSection.Content = append(currentSection.Content, trimmedLine)
			break
		}
	}

	// Always include a comment section
	if currentSection.Type == CommentTypeUnknown {
		currentSection.Type = CommentTypeText
	}

	currentComment.Sections = append(currentComment.Sections, currentSection)

	if !currentComment.Empty() {
		comments = append(comments, currentComment)
	}

	return comments
}

func isTag(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return strings.HasPrefix(trimmedLine, "+docs:")
}

func isValidYaml(body string) bool {
	var node yaml.Node
	return yaml.Unmarshal([]byte(body), &node) == nil
}

func isValidNonScalarYaml(body string) bool {
	var node yaml.Node
	if yaml.Unmarshal([]byte(body), &node) != nil {
		return false
	}

	return isNonScalar(&node)
}

func isNonScalar(node *yaml.Node) bool {
	switch node.Kind {
	case yaml.DocumentNode:
		return len(node.Content) > 0 && isNonScalar(node.Content[0])
	case yaml.MappingNode:
		// We dont want false positives on lines like "for example: " or "note: ".
		// To help with this we add some extra rules, like keys can't contain
		// spaces or start with capital letters
		return len(node.Content) > 0 &&
			!strings.Contains(node.Content[0].Value, " ") &&
			len(node.Content[0].Value) > 0 &&
			!unicode.IsUpper(rune(node.Content[0].Value[0]))
	case yaml.SequenceNode:
		return true
	default:
		return false
	}
}

func printCodeTrimLeft(lines []string) string {
	minLeadingSpaces := 1000

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		leadingSpaces := countLeadingSpaces(line)

		if leadingSpaces < minLeadingSpaces {
			minLeadingSpaces = leadingSpaces
		}
	}

	prefix := strings.Repeat(" ", minLeadingSpaces)

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(strings.TrimPrefix(line, prefix))
		sb.WriteRune('\n')
	}

	return strings.TrimSpace(sb.String())
}

func countLeadingSpaces(line string) int {
	leadingSpaces := 0
	for _, r := range line {
		if r == ' ' {
			leadingSpaces++
		} else {
			break
		}
	}
	return leadingSpaces
}
