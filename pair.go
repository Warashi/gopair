package gopair

import (
	"log"
	"math"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type (
	Seeds     map[string]int
	Candidate map[string]int
)

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

	var combs []Candidate
	for _, k := range keycombs {
		combs = append(combs, s.comb(k)...)
	}

	for {
		log.Println(len(combs))
		newc := compact(combs, keys)
		if len(newc) == len(combs) {
			break
		}
		combs = newc
	}

	var count int
	for _, c := range combs {
		for _, k := range keys {
			if _, ok := c[k]; !ok {
				count++
				c[k] = 0
			}
		}
	}
	log.Println("ゼロ埋めの数", count)
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
			m := maps.Clone(m)
			if m == nil {
				m = make(Candidate)
			}
			m[k] = v
			combinations = append(combinations, m)
		}
	}
	return combinations
}

func compact(s []Candidate, keys []string) []Candidate {
	log.Println(len(s))

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

	slices.SortFunc(s, func(a, b Candidate) bool {
		for _, k := range keys {
			av, aok := a[k]
			bv, bok := b[k]
			if aok && bok && av != bv {
				return av < bv
			}
			if aok && !bok {
				return true
			}
			if !aok && bok {
				return false
			}
		}
		return true
	})

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

func contains[K comparable, V comparable, M ~map[K]V](large, small M) bool {
	if len(large) < len(small) {
		return false
	}
	for k, v := range small {
		if vv, ok := large[k]; !ok || v != vv {
			return false
		}
	}
	return true
}

func mergable[K comparable, V comparable, M ~map[K]V](a, b M) (ok bool) {
	for k, av := range a {
		if bv, ok := b[k]; ok && av != bv {
			return false
		}
	}
	return true
}

func merge[K comparable, V any, M ~map[K]V](a, b M) M {
	a = maps.Clone(a)
	maps.Copy(a, b)
	return a
}
