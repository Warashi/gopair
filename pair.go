package gopair

import (
	"log"
	"math"

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

	bias := bias(combs)

	for {
		log.Println(len(combs))
		var newc []map[string]int
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

func compact(s []map[string]int, b map[string]int) ([]map[string]int, map[string]int) {
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
	s[i] = merge(s[i], s[j])
	s = slices.Delete(s, j, j+1)
	return s, bias(s)
}

func score(s []map[string]int, b map[string]int, i, j int) int {
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

func scoreHeavy(s []map[string]int, i, j int) int {
	s = slices.Clone(s)
	s[i] = merge(s[i], s[j])
	s = slices.Delete(s, j, j+1)

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

func bias(s []map[string]int) map[string]int {
	score := make(map[string]int)
	for _, ss := range s {
		for k := range ss {
			score[k]++
		}
	}
	return score
}

func mergable[K comparable, V comparable](a, b map[K]V) (ok bool) {
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
