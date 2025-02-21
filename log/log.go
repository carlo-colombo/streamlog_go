package log

type Encoder interface {
	Encode(v any) error
}

type Log struct {
	Line string `json:"line"`
}

func (l Log) Encode(encoder Encoder) {
	encoder.Encode(l)
}
