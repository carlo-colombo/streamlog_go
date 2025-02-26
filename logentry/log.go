package logentry

import "fmt"

type Encoder interface {
	Encode(v any) error
}

type Log struct {
	Line string `json:"line"`
}

func (l Log) Encode(encoder Encoder) error {
	if err := encoder.Encode(l); err != nil {
		return fmt.Errorf("impossible to encode log entry: %w", err)
	}
	return nil
}
