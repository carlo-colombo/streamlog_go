package logentry

import (
	"fmt"
	"time"
)

type Encoder interface {
	Encode(v any) error
}

type Log struct {
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`
}

func NewLog(line string) Log {
	return Log{
		Line:      line,
		Timestamp: time.Now(),
	}
}

func (l Log) Encode(encoder Encoder) error {
	if err := encoder.Encode(l); err != nil {
		return fmt.Errorf("impossible to encode log entry: %w", err)
	}
	return nil
}
