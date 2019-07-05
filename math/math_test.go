package math

import "testing"

func TestMin(t *testing.T) {

	if Min(1, 2) == 2 {
		t.Fail()
	}

	if Min(2, 1) == 2 {
		t.Fail()
	}
}
