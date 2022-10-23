package gopair

import (
	"encoding/json"
	"log"
	"math"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type (
	Seeds     map[string]int
	Candidate map[string]int
)

func (c Candidate) String() string {
	b, _ := json.Marshal(c)
	return string(b)
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

	var combs []Candidate
	for _, k := range keycombs {
		combs = append(combs, s.comb(k)...)
	}

	bias := bias(combs)

	for {
		log.Println(len(combs))
		var newc []Candidate
		newc, bias = compact(combs, bias)
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

func copyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
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
			m := copyMap(m)
			m[k] = v
			combinations = append(combinations, m)
		}
	}
	return combinations
}

func compact(s []Candidate, b map[string]int) ([]Candidate, map[string]int) {
	log.Println(len(s))
	maxScore, maxScoreI, maxScoreJ := math.MinInt, 0, 0
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if !mergable(s[i], s[j]) {
				continue
			}
			if score := score(s, b, i, j); maxScore < score {
				maxScore, maxScoreI, maxScoreJ = score, i, j
			}
		}
	}
	if maxScore == math.MinInt {
		return s, b
	}
	i, j := maxScoreI, maxScoreJ
	merged := merge(s[i], s[j])
	for i := 0; i < len(s); {
		if contains(merged, s[i]) {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
			continue
		}
		i++
	}
	s = append(s, merged)
	return s, bias(s)
}

func score(s []Candidate, b map[string]int, i, j int) int {
	if len(s) < 100 {
		return scoreHeavy(s, i, j)
	}
	var score int
	for k := range s[i] {
		score += b[k]
	}
	for k := range s[j] {
		score += b[k]
	}
	return score
}

func scoreHeavy(s []Candidate, i, j int) int {
	s = slices.Clone(s)
	s[i] = merge(s[i], s[j])
	s[j] = s[len(s)-1]
	s = s[:len(s)-1]

	var score int
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if mergable(s[i], s[j]) {
				score++
			}
		}
	}
	return score
}

func bias(s []Candidate) map[string]int {
	score := make(map[string]int)
	for _, ss := range s {
		for k := range ss {
			score[k]++
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
		bv, ok := b[k]
		if ok && av != bv {
			return false
		}
	}
	return true
}

func merge[K comparable, V any](a, b map[K]V) map[K]V {
	r := copyMap(a)
	for k, v := range b {
		r[k] = v
	}
	return r
}
