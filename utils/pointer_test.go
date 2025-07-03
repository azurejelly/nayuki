package utils

import "testing"

func TestPtr(t *testing.T) {
	val := 5
	ptr := Ptr(val)
	if ptr == nil || *ptr != val {
		t.Errorf("Ptr(%d) = %v; want pointer to %d", val, ptr, val)
	}

	str := "test"
	strPtr := Ptr(str)
	if strPtr == nil || *strPtr != str {
		t.Errorf("Ptr(%q) = %v; want pointer to %q", str, strPtr, str)
	}
}
