package gomux

import "testing"

func TestNewPane(t *testing.T) {
	window := &Window{
		Number: 1,
		Name:   "foo",
	}

	testPane := NewPane(1, window)
	if testPane.Number != 1 {
		t.Error(testPane.Number)
	}

}
