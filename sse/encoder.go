package sse

import (
	"fmt"
	"github.com/carlo-colombo/streamlog_go/log"
	"html/template"
	"io"
	"strings"
)

type Encoder struct {
	writer io.Writer
	line   string
}

func NewEncoder(w io.Writer, line string) Encoder {
	return Encoder{
		w,
		strings.ReplaceAll(line, "\n", ""),
	}
}

func (e Encoder) Encode(v any) error {
	l, ok := v.(log.Log)
	if !ok {
		return fmt.Errorf("encoder can only encode a log object")
	}

	logTmpl := template.Must(template.New("log").Parse(e.line))

	fmt.Fprintf(e.writer, "data: ")
	logTmpl.Execute(e.writer, l)
	fmt.Fprint(e.writer, "\n\n")

	return nil
}
