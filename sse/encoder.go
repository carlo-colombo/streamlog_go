package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

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

type rawMessage struct {
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`
}

func (e Encoder) Encode(v any) error {
	l, ok := v.(logentry.Log)
	if !ok {
		return fmt.Errorf("encoder can only encode a log object")
	}

	// Use a custom type to avoid escaping in the Line field
	raw := rawMessage{
		Line:      l.Line,
		Timestamp: l.Timestamp,
	}

	// Use json.Marshal with HTMLEscape disabled
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(raw); err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Remove the trailing newline that json.Encoder adds
	data := buffer.Bytes()
	data = data[:len(data)-1]

	fmt.Fprintf(e.writer, "data: %s\n\n", data)
	return nil
}
