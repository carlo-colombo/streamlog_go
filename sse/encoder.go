package sse

import (
	"fmt"
	"github.com/carlo-colombo/streamlog_go/log"
	"io"
)

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) Encoder {
	return Encoder{w}
}

func (e Encoder) Encode(v any) error {
	l, ok := v.(log.Log)
	if !ok {
		return fmt.Errorf("encoder can only encode a log object")
	}

	fmt.Fprintf(e.writer, "data: %s\n\n", l.Line)

	return nil
}
