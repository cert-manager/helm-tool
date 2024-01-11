package heuristics

import (
	"regexp"
	"strings"
	"unicode"
)

var numericListExp = regexp.MustCompile(`^[0-9]+\s*[-\.]$`)

// RecutNewLines attempts to recut new lines that may have been wrapped to
// remove the wrapping, while trying to preserve intentional new lines for
// lists, code etc.
func RecutNewLines(lines []string) []string {
	// Return early for an empty slice
	if len(lines) == 0 {
		return nil
	}

	leadingSpacesCount := countLeadingSpaces(lines[0])
	leadingSpaces := strings.Repeat(" ", leadingSpacesCount)

	var parsedLines []string
	var currentLine []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Get the previous line
		previousLine := ""
		if l := len(parsedLines); l != 0 {
			previousLine = parsedLines[l-1]
		}

		// If the line is empty, then the new line is considered intentional
		if trimmedLine == "" {
			if len(currentLine) != 0 {
				previousLine = strings.Join(currentLine, " ")
				parsedLines = append(parsedLines, previousLine)
			}
			if strings.TrimSpace(previousLine) != "" {
				parsedLines = append(parsedLines, "")
			}
			currentLine = nil
			continue
		}

		lineWithoutLeadingSpaces := trimSpaceRight(strings.TrimPrefix(line, leadingSpaces))
		previousLineLeadingSpacesCount := 0
		currentLineLeadingSpacesCount := countLeadingSpaces(lineWithoutLeadingSpaces)

		// If we are further indented then the start of the comment, we assume
		// the new line was intended
		if c := firstChar(lineWithoutLeadingSpaces); unicode.IsSpace(c) {
			if len(currentLine) != 0 {
				parsedLines = append(parsedLines, strings.Join(currentLine, " "))
			}
			parsedLines = append(parsedLines, lineWithoutLeadingSpaces)
			currentLine = nil
			continue
		}

		// If we start with a non alphanumeric character then assume the
		// new line was intended (it could be a list for example)
		if c := firstChar(trimmedLine); !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			if len(currentLine) != 0 {
				parsedLines = append(parsedLines, strings.Join(currentLine, " "))
			}
			parsedLines = append(parsedLines, lineWithoutLeadingSpaces)
			currentLine = nil
			continue
		}

		// If we are probably part of a numeric list, the new line is intended
		if numericListExp.MatchString(trimmedLine) {
			if len(currentLine) != 0 {
				parsedLines = append(parsedLines, strings.Join(currentLine, " "))
			}
			parsedLines = append(parsedLines, lineWithoutLeadingSpaces)
			currentLine = nil
			continue
		}

		// If the previous line ends with certain characters
		// (a colon, :, for example) we assume the new line was intentional and
		// not just text wrapping
		if lastChar(previousLine) == ':' {
			if len(currentLine) != 0 {
				parsedLines = append(parsedLines, strings.Join(currentLine, " "))
			}
			currentLine = []string{trimSpaceRight(lineWithoutLeadingSpaces)}
			continue
		}

		// If the line is intentionally short, then the new line is probably
		// intentional
		if c := firstChar(trimmedLine); unicode.IsUpper(c) && len(trimmedLine) < 50 {
			if len(currentLine) != 0 {
				parsedLines = append(parsedLines, strings.Join(currentLine, " "))
			}
			parsedLines = append(parsedLines, lineWithoutLeadingSpaces)
			currentLine = nil
			continue
		}

		if currentLineLeadingSpacesCount != previousLineLeadingSpacesCount {
			if len(currentLine) != 0 {
				parsedLines = append(parsedLines, strings.Join(currentLine, " "))
			}
			parsedLines = append(parsedLines, lineWithoutLeadingSpaces)
			currentLine = nil
			continue
		}

		// Line is likely wrapping, keep on current line
		currentLine = maybeAppendAddingFullStop(currentLine, trimmedLine)
	}

	if len(currentLine) != 0 {
		parsedLines = append(parsedLines, strings.Join(currentLine, " "))
	}

	// Trim trailing new line
	if l := len(parsedLines); l != 0 && strings.TrimSpace(parsedLines[l-1]) == "" {
		parsedLines = parsedLines[:l-1]
	}

	return parsedLines
}

func maybeAppendAddingFullStop(line []string, extra string) []string {
	// Ensure the leading an trailing spaces are removed first
	extra = strings.TrimSpace(extra)

	// Don't append empty lines
	if extra == "" {
		return line
	}

	// If lines is empty, just append
	l := len(line)
	if l == 0 {
		return append(line, extra)
	}

	previousLine := trimSpaceRight(line[l-1])

	// If the previous line is empty, just append
	if previousLine == "" {
		return append(line, extra)
	}

	// If the previous line is not empty, and does not end in a full stop,
	// add a full stop.
	if unicode.IsUpper(firstChar(extra)) && lastChar(previousLine) != '.' && unicode.IsLetter(lastChar(previousLine)) {
		line[l-1] = previousLine + "."
	}

	return append(line, extra)
}

func countLeadingSpaces(line string) int {
	for i, char := range line {
		if char != ' ' {
			return i
		}
	}
	return 0
}

func trimLeadingSpaces(lines []string) []string {
	// If there are no lines, there is nothing to trim
	if len(lines) == 0 {
		return lines
	}

	// Get the leading spaces from the first line
	leadingSpacesCount := countLeadingSpaces(lines[0])
	leadingSpaces := strings.Repeat(" ", leadingSpacesCount)

	trimmedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmedLines[i] = strings.TrimPrefix(line, leadingSpaces)
	}

	return trimmedLines
}

func lastChar(str string) rune {
	if str == "" {
		return 0
	}

	return rune(str[len(str)-1])
}

func firstChar(str string) rune {
	if str == "" {
		return 0
	}

	return rune(str[0])
}

func trimSpaceRight(str string) string {
	return strings.TrimRightFunc(str, unicode.IsSpace)
}
