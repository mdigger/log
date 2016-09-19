package log

import "testing"

func TestCaller(t *testing.T) {
	caller := MakeCaller(0)
	if caller.String() != "caller_test.go:6" {
		t.Error("bad caller:", caller)
	}

	caller = MakeCaller(30)
	if caller.String() != "<unknown>" {
		t.Error("bad caller:", caller)
	}
}
