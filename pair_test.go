package gopair

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	s := Seeds{
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

func includes(set []Candidate, elem map[string]int) bool {
	kv := make([]KV[string, int], 0, len(elem))
	for k, v := range elem {
		kv = append(kv, KV[string, int]{k, v})
	}
	for _, s := range set {
		if contains(s, kv) {
			return true
		}
	}
	return false
}

func TestGenerate2(t *testing.T) {
	//	t.Skip()
	s := Seeds{
		"a": 5,
		"b": 5,
		"c": 5,
		"d": 5,
		"e": 5,
	}
	got := s.Generate(2)
	t.Log(len(got))
}
