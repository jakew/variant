/*
Package variant is used for creating complex tests in the easiest way possible.
Instead of refactoring the flow of the test into individual functions and
creating a new `t.Run` instance for each, just write the function once and put
the varying pieces into a `t.Variant` call.
*/

package variant

import (
	"fmt"
)

/*
TODO:
- Update variants to be nodes in a tree, rather than a collection.
- Each node must be able to provide its full list of sub-combinations.
- Each full list of sub-combinations must provide all combinations of the
  variations within it.
- All nodes must be able to provide a name for each unique combination of
  variations.
- Allow for missing variants: If a test has variatns A1, A2, A3, and only A1 and
  A2 leads to B, then this should be accounted for.

TODO: Add in testing.T
- All iterations should use a named testing.T.Run instance.
- The Variation should be able to run only the specific test by reviewing
  testing.T.Name and acting accordingly.
*/

type Variant struct {
	invocationI int
	Collection  *Collection
	// The indexes of the invocations to follow.
	Permutations []int
}

type Collection struct {
	// This a 2d slice, sorted by the invocation #, followed by the options.
	Variants [][]string
}

// Combinations returns back all possible combinations of the variants.
func (c *Collection) Combinations() [][]int {
	var vars [][]int
	for _, v := range c.Variants {
		vars = AddCombination(vars, len(v))
	}

	return vars
}

func (c *Collection) NameForPermutation(perm []int) string {
	var name string
	for i, v := range perm {
		name += c.Variants[i][v]
	}
	return name
}

// AddCombination takes a 2d slice of ints vars and an int n, and, for i, each
// of the numbers between 0 to n, copies the slice and appends i to the end of
// the slice.
func AddCombination(vars [][]int, n int) [][]int {
	if len(vars) < 1 {
		vars = [][]int{{}}
	}

	var newVars [][]int
	for _, v := range vars {
		for i := 0; i < n; i++ {
			newV := make([]int, len(v))
			copy(newV, v)
			newV = append(newV, i)
			newVars = append(newVars, newV)
		}
	}

	return newVars
}

type NamedVariant struct {
	name string
	fn   func(v *Variant)
}

type M map[string]func(v *Variant)

func New(root func(v *Variant)) {
	v := &Variant{
		invocationI: 0,
		Collection: &Collection{
			Variants: [][]string{},
		},
	}

	root(v)

	combos := v.Collection.Combinations()
	fmt.Printf("Combinations: %v", combos)

	for i, perm := range combos {
		if i == 0 {
			continue
		}
		name := v.Collection.NameForPermutation(perm)
		fmt.Printf("Running test %s\n", name)

		// t.Run goes here. w/ name.
		root(&Variant{
			invocationI:  0,
			Collection:   v.Collection,
			Permutations: perm,
		})
	}
}

func (v *Variant) Variants(m M) {
	defer func() { v.invocationI++ }()

	if len(v.Collection.Variants) <= v.invocationI {
		namedVariants := make([]string, 0, len(m))
		for k := range m {
			namedVariants = append(namedVariants, k)
		}

		v.Collection.Variants = append(v.Collection.Variants, namedVariants)
	}

	if len(v.Collection.Variants[v.invocationI]) < 1 {
		return
	}

	// if we don't know which permuation this is, run the first one.
	perm := 0
	if len(v.Permutations) > v.invocationI {
		perm = v.Permutations[v.invocationI]
	}

	fmt.Printf("invoking %d %d\n", v.invocationI, perm)

	m[v.Collection.Variants[v.invocationI][perm]](v)
}
