package gopair_test

import (
	"testing"

	"github.com/Warashi/gopair"
)

func TestGenerate(t *testing.T) {
	s := gopair.Seeds{
		"a": 3,
		"b": 3,
		"c": 3,
	}
	got := s.Generate(2)
	wantIncludes := []map[string]int{
		{"a": 0, "b": 0},
		{"a": 0, "b": 1},
		{"a": 0, "b": 2},
		{"a": 1, "b": 0},
		{"a": 1, "b": 1},
		{"a": 1, "b": 2},
		{"a": 2, "b": 0},
		{"a": 2, "b": 1},
		{"a": 2, "b": 2},
		{"a": 0, "c": 0},
		{"a": 0, "c": 1},
		{"a": 0, "c": 2},
		{"a": 1, "c": 0},
		{"a": 1, "c": 1},
		{"a": 1, "c": 2},
		{"a": 2, "c": 0},
		{"a": 2, "c": 1},
		{"a": 2, "c": 2},
		{"b": 0, "c": 0},
		{"b": 0, "c": 1},
		{"b": 0, "c": 2},
		{"b": 1, "c": 0},
		{"b": 1, "c": 1},
		{"b": 1, "c": 2},
		{"b": 2, "c": 0},
		{"b": 2, "c": 1},
		{"b": 2, "c": 2},
	}

	for _, want := range wantIncludes {
		if !includes(got, want) {
			t.Errorf("want includes %v, got %v", want, got)
		}
	}
}

func includes(set []map[string]int, elem map[string]int) bool {
loop:
	for _, s := range set {
		for k, v := range elem {
			if s[k] != v {
				continue loop
			}
		}
		return true
	}
	return false
}
