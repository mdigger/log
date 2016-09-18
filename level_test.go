package log

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLevel(t *testing.T) {
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(Level(4))
	enc.Encode(struct {
		Level Level `json:"level"`
	}{4})
}
