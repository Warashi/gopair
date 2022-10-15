package gopair

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Seeds map[string]int

func (s Seeds) Generate(order int) []map[string]int {
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

	var combs []map[string]int
	for _, k := range keycombs {
		combs = append(combs, s.comb(k)...)
	}

	combs = compact(combs)

	for _, c := range combs {
		for _, k := range keys {
			if _, ok := c[k]; !ok {
				c[k] = 0
			}
		}
	}
	return combs
}

func copyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (s Seeds) comb(keys []string) []map[string]int {
	if len(keys) == 0 {
		return []map[string]int{{}}
	}
	var combinations []map[string]int

	mid := s.comb(keys[1:])
	k := keys[0]
	for _, v := range indices(s[k]) {
		for _, m := range mid {
			m := copyMap(m)
			m[k] = v
			combinations = append(combinations, m)
		}
	}
	return combinations
}

func compact(s []map[string]int) []map[string]int {
	for i := 0; i < len(s); i++ {
	inner:
		for j := len(s) - 1; i < j; j-- {
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

func merge[K comparable, V any](a, b map[K]V) map[K]V {
	r := copyMap(a)
	for k, v := range b {
		r[k] = v
	}
	return r
}
