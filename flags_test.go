package log

import (
	"fmt"
	"os"
	"testing"
)

func TestFlags(t *testing.T) {
	json := new(JSON)
	json.SetOutput(os.Stdout)
	console := new(Console)
	console.SetOutput(os.Stdout)
	log := New(json, console)

	for flag := 0; flag < Lindent<<1; flag++ {
		json.SetFlags(flag)
		console.SetFlags(flag)
		log.WithField("flag", flag).Info("")
		fmt.Println()
	}
}
