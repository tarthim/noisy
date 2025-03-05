package noisy

import (
	"testing"
)

func TestExpectErrorOnNegativeHeight(t *testing.T) {
	n := noisy{
		width:  -10,
		height: 10,
	}
	err := n.validate()
	if err == nil {
		t.Fatalf("want: validation error; got: nothing")
	}
}
