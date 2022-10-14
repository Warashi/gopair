package gopair_test

import (
	"testing"

	"github.com/Warashi/gopair"
)

func TestGenerate(t *testing.T) {
	s := gopair.Seeds{
		"a": []any{"A", "B"},
		"b": []any{"C", "D"},
		"c": []any{"E", "F"},
	}
	// t.Log(s.Generate(1))
	t.Log(s.Generate(2))
	// t.Log(s.Generate(3))
}
