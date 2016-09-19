package log

// Level is the log severity level.
type Level int8

// String converts a severity level to a string.
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

// MarshalText return the level string for text based encoding.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// Predefined log severity levels.
const (
	DebugLevel Level = -1 // debug level
	InfoLevel  Level = 0  // default severity level
	ErrorLevel Level = 1  // error level
)
