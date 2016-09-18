package log

type Level int8

const (
	Debug Level = -1
	Info  Level = 0
	Error Level = 1
)

func (l Level) String() string {
	switch {
	case l < 0:
		return "debug"
	case l > 0:
		return "error"
	default:
		return "info"
	}
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
