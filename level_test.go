package log

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLevel(t *testing.T) {
	enc := json.NewEncoder(os.Stdout)
	for _, level := range []Level{InfoLevel, DebugLevel, WarningLevel, ErrorLevel} {
		enc.Encode(level)
	}
}
