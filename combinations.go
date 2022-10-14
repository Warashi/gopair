package gopair

func indices(n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = i
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
		reversed: reversed(indices(order)),
	}
}

func (c *Combinations[T]) Next() bool {
	if c.finished {
		return false
	}

	if c.indices == nil {
		c.indices = indices(c.order)
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
