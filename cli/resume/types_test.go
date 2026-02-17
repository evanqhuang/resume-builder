package resume

import "testing"

func TestResumeTypes(t *testing.T) {
	// Basic compilation test
	r := &Resume{}
	if r == nil {
		t.Fatal("Resume should not be nil")
	}
}
