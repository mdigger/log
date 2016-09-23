package log

// Level of severity.
type Level int8

// String implements io.Stringer.
func (l Level) String() string {
	switch {
	case l < 0:
		return "debug"
	case l == 1:
		return "warning"
	case l > 1:
		return "error"
	default:
		return "info"
	}
}

// MarshalText returns the level string.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

const (
	DebugLevel   Level = -1
	InfoLevel    Level = 0
	WarningLevel Level = 1
	ErrorLevel   Level = 2
)
