package autoscaling

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestStringInSlice(t *testing.T) {
	var cases = []struct {
		s      string
		slice  []string
		expect bool
	}{
		{"", []string{}, false},
		{"a", []string{}, false},
		{"a", []string{"a"}, true},
		{"a", []string{"a", "b"}, true},
		{"a", []string{"b", "c"}, false},
		{"tom", []string{"tom", "cat"}, true},
	}

	for _, c := range cases {
		//call
		got := stringInSlice(c.s, c.slice)
		//assert
		assert.Equal(t, c.expect, got)
	}
}
