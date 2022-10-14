package gopair

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Seeds map[string][]any

func indices[T any](s []T) []int {
	r := make([]int, 0, len(s))
	for i := range s {
		r = append(r, i)
	}
	return r
}

func reversed[T any](s []T) []T {
	r := make([]T, len(s))
	copy(r, s)
	for i := 0; i < len(r)/2; i++ {
		left, right := i, len(r)-1-i
		r[left], r[right] = r[right], r[left]
	}
	return r
}

func pick[T any](s []T, indices []int) []T {
	r := make([]T, 0, len(indices))
	for _, i := range indices {
		r = append(r, s[i])
	}
	return r
}

type Combinations[T any] struct {
	from     []T
	indices  []int
	finished bool
	order    int

	// reversed : reversed(range(order))
	reversed []int
}

func NewCombinations[T any](from []T, order int) *Combinations[T] {
	return &Combinations[T]{
		from:     from,
		indices:  nil,
		finished: false,
		order:    order,
		reversed: reversed(indices(make([]struct{}, order))),
	}
}

func (c *Combinations[T]) Next() bool {
	if c.finished {
		return false
	}

	if c.indices == nil {
		c.indices = indices(make([]struct{}, c.order))
		return true
	}

	for _, i := range c.reversed {
		if c.indices[i] != i+len(c.from)-c.order {
			c.indices[i]++
			for j := i + 1; j < c.order; j++ {
				c.indices[j] = c.indices[j-1] + 1
			}
			return true
		}
	}

	c.finished = true
	return false
}

func (c *Combinations[T]) Value() []T {
	return pick(c.from, c.indices)
}

func (s Seeds) Generate(order int) []map[string]any {
	if len(s) < order {
		return nil
	}

	keys := maps.Keys(s)
	slices.Sort(keys)

	var keycombs [][]string
	c := NewCombinations(keys, order)
	for c.Next() {
		keycombs = append(keycombs, c.Value())
	}

	var combs []map[string]any
	for _, k := range keycombs {
		combs = append(combs, s.comb(k)...)
	}

	combs = compact(combs)

	for _, c := range combs {
		for _, k := range keys {
			if _, ok := c[k]; !ok {
				c[k] = s[k][0]
			}
		}
	}
	return combs
}

func copyMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (s Seeds) comb(keys []string) []map[string]any {
	if len(keys) == 0 {
		return []map[string]any{{}}
	}
	var combinations []map[string]any

	mid := s.comb(keys[1:])
	k := keys[0]
	for _, v := range s[k] {
		for _, m := range mid {
			m := copyMap(m)
			m[k] = v
			combinations = append(combinations, m)
		}
	}
	return combinations
}

func compact(s []map[string]any) []map[string]any {
	for i := 0; i < len(s); i++ {
	inner:
		for j := i + 1; j < len(s); j++ {
			im := s[i]
			jm := s[j]
			for ik, iv := range im {
				jv, ok := jm[ik]
				if ok && iv != jv {
					// not margeable
					continue inner
				}
			}
			s[i] = merge(im, jm)
			s = slices.Delete(s, j, j+1)
		}
	}
	return s
}

func merge(a, b map[string]any) map[string]any {
	r := copyMap(a)
	for k, v := range b {
		r[k] = v
	}
	return r
}
