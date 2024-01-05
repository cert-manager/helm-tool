package parser

import (
	"fmt"
	"io"
	"strings"
)

type PathComponent interface {
	Append(idx int, path io.Writer)
}

type StringPathComponent string

func (s StringPathComponent) Append(idx int, w io.Writer) {
	if strings.Contains(string(s), ".") {
		fmt.Fprintf(w, "[%q]", s)
	} else if idx == 0 {
		fmt.Fprintf(w, "%s", s)
	} else {
		fmt.Fprintf(w, ".%s", s)
	}
}

type IndexPathComponent int

func (i IndexPathComponent) Append(idx int, w io.Writer) {
	fmt.Fprintf(w, "[%d]", i)
}

type Path []PathComponent

func (p Path) WithProperty(part string) Path {
	return append(p, StringPathComponent(part))
}

func (p Path) WithIndex(idx int) Path {
	return append(p, IndexPathComponent(idx))
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
