package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type pathComponent interface {
	Append(idx int, path io.Writer)
}

type stringPathComponent string

func (s stringPathComponent) Append(idx int, w io.Writer) {
	if strings.Contains(string(s), ".") {
		fmt.Fprintf(w, "[%q]", s)
	} else if idx == 0 {
		fmt.Fprintf(w, "%s", s)
	} else {
		fmt.Fprintf(w, ".%s", s)
	}
}

type indexPathComponent int

func (i indexPathComponent) Append(idx int, w io.Writer) {
	fmt.Fprintf(w, "[%d]", i)
}

type Path []pathComponent

func ParsePath(pathString string) (Path, error) {
	scanner := bytes.NewBufferString(pathString)

	path := Path{}

	for {
		switch k, last, err := runesUntil(scanner, runeSet([]rune{'[', '.'})); {
		case err == io.EOF && len(k) > 0:
			return append(path, stringPathComponent(k)), nil
		case err == io.EOF && len(k) == 0:
			return path, nil
		case err != nil:
			return path, err
		case last == '[':
			path = append(path, stringPathComponent(k))

			betweenBrackets, _, err := runesUntil(scanner, runeSet([]rune{']'}))
			if err != nil {
				return path, fmt.Errorf("unexpected end of path")
			}

			if len(betweenBrackets) == 0 {
				return path, fmt.Errorf("unexpected empty array index")
			}

			if betweenBrackets[0] == '"' && betweenBrackets[len(betweenBrackets)-1] == '"' {
				// betweenBrackets is a string
				path = append(path, stringPathComponent(betweenBrackets[1:len(betweenBrackets)-1]))
			} else {
				idx, err := strconv.Atoi(string(betweenBrackets))
				if err != nil {
					return path, fmt.Errorf("unexpected array index: %q", betweenBrackets)
				}

				path = append(path, indexPathComponent(idx))
			}

			if r, _, e := scanner.ReadRune(); e == nil && r != '.' {
				return path, fmt.Errorf("unexpected token %q", r)
			}
		case last == '.':
			path = append(path, stringPathComponent(k))
		default:
			return nil, fmt.Errorf("parse error: unexpected token %v", last)
		}
	}
}

func (p Path) WithProperty(part string) Path {
	return append(p, stringPathComponent(part))
}

func (p Path) WithIndex(idx int) Path {
	return append(p, indexPathComponent(idx))
}

func (p Path) Parent() Path {
	l := len(p)
	if l == 0 {
		return nil
	}

	return append(Path{}, p[:l-1]...)
}

func (p Path) String() string {
	sb := strings.Builder{}
	for i, part := range p {
		part.Append(i, &sb)
	}
	return sb.String()
}

func runeSet(r []rune) map[rune]bool {
	s := make(map[rune]bool, len(r))
	for _, rr := range r {
		s[rr] = true
	}
	return s
}

func runesUntil(in io.RuneReader, stop map[rune]bool) ([]rune, rune, error) {
	v := []rune{}
	for {
		switch r, _, e := in.ReadRune(); {
		case e != nil:
			return v, r, e
		case inMap(r, stop):
			return v, r, nil
		case r == '\\':
			next, _, e := in.ReadRune()
			if e != nil {
				return v, next, e
			}
			v = append(v, next)
		default:
			v = append(v, r)
		}
	}
}

func inMap(r rune, m map[rune]bool) bool {
	_, ok := m[r]
	return ok
}
