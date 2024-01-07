package parser

import (
	"fmt"
	"io"
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
