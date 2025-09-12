package indenter

import (
	"fmt"
	"strings"
)

type T = Indenter

type Indenter struct {
	indent  int
	builder strings.Builder
}

func New(indent int) Indenter {
	return Indenter{
		indent: indent,
	}
}

func (i *Indenter) Sprintf(level int, format string, args ...any) {
	i.builder.WriteString(strings.Repeat(" ", level*i.indent))
	i.builder.WriteString(fmt.Sprintf(format, args...))
}

func (i *Indenter) String() string {
	return i.builder.String()
}
