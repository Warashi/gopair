package gopair

import (
	"cmp"
	"math"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type (
	Seeds                            map[string]int
	KV[K cmp.Ordered, V cmp.Ordered] struct {
		key   K
		value V
	}
	Candidate []KV[string, int]
)

func keyCompareFunc[K cmp.Ordered, V cmp.Ordered](a, b KV[K, V]) int {
	return cmp.Compare(a.key, b.key)
}

func kvCompareFunc[K cmp.Ordered, V cmp.Ordered](a, b KV[K, V]) int {
	return cmp.Compare(a.key, b.key)*10 + cmp.Compare(a.value, b.value)
}

func kvSortFunc[K cmp.Ordered, V cmp.Ordered](a, b KV[K, V]) bool {
	return kvCompareFunc(a, b) < 0
}

func candidateCompareFunc(a, b Candidate) int {
	slices.SortFunc(a, kvSortFunc)
	slices.SortFunc(b, kvSortFunc)
	return slices.CompareFunc(a, b, kvCompareFunc)
}

func candidateSortFunc(a, b Candidate) bool {
	return candidateCompareFunc(a, b) < 0
}

func (s Seeds) Generate(order int) []Candidate {
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
	slices.SortStableFunc(keycombs, func(a, b []string) bool {
		slices.Sort(a)
		slices.Sort(b)
		return slices.Compare(a, b) < 0
	})

	var combs []Candidate
	for _, k := range keycombs {
		combs = append(combs, s.comb(k)...)
	}
	slices.SortStableFunc(combs, candidateSortFunc)

	for {
		newc := compact(combs, keys)
		if len(newc) == len(combs) {
			break
		}
		combs = newc
	}

	for _, k := range keys {
		for _, c := range combs {
			if _, ok := slices.BinarySearchFunc(c, KV[string, int]{key: k}, keyCompareFunc); !ok {
				c = append(c, KV[string, int]{key: k, value: 0})
			}
			slices.SortFunc(c, kvSortFunc)
		}
	}
	return combs
}

func (s Seeds) comb(keys []string) []Candidate {
	if len(keys) == 0 {
		return make([]Candidate, 1)
	}
	var combinations []Candidate

	mid := s.comb(keys[1:])
	k := keys[0]
	for _, v := range indices(s[k]) {
		for _, m := range mid {
			m := slices.Clone(m)
			m = append(m, KV[string, int]{key: k, value: v})
			combinations = append(combinations, m)
		}
	}
	return combinations
}

func compact(s []Candidate, keys []string) []Candidate {
	maxScore, maxScoreI, maxScoreJ := math.MinInt, 0, 0
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if !mergable(s[i], s[j]) {
				continue
			}
			if score := score(s, i, j); maxScore < score {
				maxScore, maxScoreI, maxScoreJ = score, i, j
			}
		}
	}
	if maxScore == math.MinInt {
		return s
	}

	merged := merge(s[maxScoreI], s[maxScoreJ])
	for i := 0; i < len(s); {
		if contains(merged, s[i]) {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
			continue
		}
		i++
	}
	s = append(s, merged)

	slices.SortFunc(s, candidateSortFunc)
	return s
}

func score(s []Candidate, i, j int) (score int) {
	merged := merge(s[i], s[j])
	for _, s := range s {
		if contains(merged, s) {
			score++
		}
	}

	return score
}

func contains[K cmp.Ordered, V cmp.Ordered, M ~[]KV[K, V]](large, small M) bool {
	if len(large) < len(small) {
		return false
	}
	slices.SortFunc(small, kvSortFunc)
	slices.SortFunc(large, kvSortFunc)
	if slices.CompareFunc(large, small, kvCompareFunc) == 0 {
		return true
	}
	for _, kv := range small {
		if _, ok := slices.BinarySearchFunc(large, kv, kvCompareFunc); !ok {
			return false
		}
	}
	return true
}

func mergable[K cmp.Ordered, V cmp.Ordered, M ~[]KV[K, V]](a, b M) (ok bool) {
	if len(a) > len(b) {
		a, b = b, a
	}
	for _, kv := range a {
		idx, ok := slices.BinarySearchFunc(b, KV[K, V]{key: kv.key}, keyCompareFunc)
		// 同じ key で value が異なる場合はマージできない
		if ok && kv.value != b[idx].value {
			return false
		}
	}
	return true
}

func merge[K cmp.Ordered, V cmp.Ordered, M ~[]KV[K, V]](a, b M) M {
	merged := append(slices.Clone(a), b...)
	slices.SortFunc(merged, kvSortFunc)
	slices.CompactFunc(merged, func(a, b KV[K, V]) bool {
		return kvCompareFunc(a, b) == 0
	})
	return merged
}
