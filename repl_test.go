package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"simple":              {input: " hello  world ", want: []string{"hello", "world"}},
		"nothing":             {input: "", want: nil},
		"mixed case":          {input: "HelLO WOrLd", want: []string{"hello", "world"}},
		"no whitespace":       {input: "helloworld", want: []string{"helloworld"}},
		"trailing whitespace": {input: "hello world    ", want: []string{"hello", "world"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := cleanInput(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
