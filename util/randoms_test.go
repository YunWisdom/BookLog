

package util

import "testing"

func TestRandInts(t *testing.T) {
	ints := RandInts(0, 50, 5)
	if 5 != len(ints) {
		t.Errorf("generate random integers failed")
	}
}

func TestRandString(t *testing.T) {
	a := RandString(16)
	if 16 != len(a) {
		t.Error("generate random string failed")
	}
	b := RandString(16)
	if a == b {
		t.Error("generate random string failed")
	}
	t.Log(a, b)
}
