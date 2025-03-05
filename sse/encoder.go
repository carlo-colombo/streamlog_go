package sse

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/carlo-colombo/streamlog_go/logentry"
)

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) Encoder {
	return Encoder{
		writer: w,
	}
}

func (e Encoder) Encode(v any) error {
	l, ok := v.(logentry.Log)
	if !ok {
		return fmt.Errorf("encoder can only encode a log object")
	}

	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	fmt.Fprintf(e.writer, "data: %s\n\n", string(data))
	return nil
}
