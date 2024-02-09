/*
Copyright 2021 The cert-manager Authors.

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

package paths

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

type mapPathComponent string

func (s mapPathComponent) Append(idx int, w io.Writer) {
	if strings.Contains(string(s), ".") {
		fmt.Fprintf(w, "[%q]", s)
	} else if idx == 0 {
		fmt.Fprintf(w, "%s", s)
	} else {
		fmt.Fprintf(w, ".%s", s)
	}
}

type arrayPathComponent int

func (i arrayPathComponent) Append(idx int, w io.Writer) {
	fmt.Fprintf(w, "[%d]", i)
}

func IsArrayPathComponent(pc pathComponent) bool {
	_, ok := pc.(arrayPathComponent)
	return ok
}

func SegmentString(pc pathComponent) string {
	sb := strings.Builder{}
	pc.Append(0, &sb)
	return sb.String()
}

type Path []pathComponent

func Parse(pathString string) (Path, error) {
	scanner := bytes.NewBufferString(pathString)

	path := Path{}

	for {
		switch k, last, err := runesUntil(scanner, runeSet([]rune{'[', '.'})); {
		case err == io.EOF && len(k) > 0:
			return append(path, mapPathComponent(k)), nil
		case err == io.EOF && len(k) == 0:
			return path, nil
		case err != nil:
			return path, err
		case last == '[':
			path = append(path, mapPathComponent(k))

			betweenBrackets, _, err := runesUntil(scanner, runeSet([]rune{']'}))
			if err != nil {
				return path, fmt.Errorf("unexpected end of path")
			}

			if len(betweenBrackets) == 0 {
				return path, fmt.Errorf("unexpected empty array index")
			}

			if betweenBrackets[0] == '"' && betweenBrackets[len(betweenBrackets)-1] == '"' {
				// betweenBrackets is a string
				path = append(path, mapPathComponent(betweenBrackets[1:len(betweenBrackets)-1]))
			} else {
				idx, err := strconv.Atoi(string(betweenBrackets))
				if err != nil {
					return path, fmt.Errorf("unexpected array index: %q", betweenBrackets)
				}

				path = append(path, arrayPathComponent(idx))
			}

			if r, _, e := scanner.ReadRune(); e == nil && r != '.' {
				return path, fmt.Errorf("unexpected token %q", r)
			}
		case last == '.':
			path = append(path, mapPathComponent(k))
		default:
			return nil, fmt.Errorf("parse error: unexpected token %v", last)
		}
	}
}

func (p Path) WithProperty(part string) Path {
	nePath := append(Path{}, p...)
	nePath = append(nePath, mapPathComponent(part))
	return nePath
}

func (p Path) WithIndex(idx int) Path {
	nePath := append(Path{}, p...)
	nePath = append(nePath, arrayPathComponent(idx))
	return nePath
}

func (p Path) Property() pathComponent {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1]
}

func (p Path) Parent() Path {
	l := len(p)
	if l == 0 {
		return nil
	}

	return append(Path{}, p[:l-1]...)
}

func (p Path) IsSubPathOf(other Path) bool {
	if len(p) > len(other) {
		return false
	}

	for i, part := range p {
		if part != other[i] {
			return false
		}
	}

	return true
}

// Expand the path by n levels, following the path of other.
func (p Path) Expand(other Path, n int) Path {
	// p must be a subpath of other
	if !p.IsSubPathOf(other) {
		return nil
	}

	return append(Path{}, other[:len(p)+n]...)
}

func (p Path) Equal(other Path) bool {
	if len(p) != len(other) {
		return false
	}

	return p.IsSubPathOf(other)
}

func (p Path) String() string {
	sb := strings.Builder{}
	for i, part := range p {
		part.Append(i, &sb)
	}
	return sb.String()
}

func (p Path) PatternString() string {
	sb := strings.Builder{}
	for i, part := range p {
		if _, ok := part.(arrayPathComponent); ok {
			sb.WriteString("[*]")
			continue
		}
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
