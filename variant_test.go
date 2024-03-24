package variant_test

import (
	"strings"
	"testing"

	"github.com/jakew/variant"
)

func TestAddCombination(t *testing.T) {
	vars := [][]int{{1, 2}, {3, 4}}
	n := 3
	expected := [][]int{
		{1, 2, 0},
		{1, 2, 1},
		{1, 2, 2},
		{3, 4, 0},
		{3, 4, 1},
		{3, 4, 2},
	}
	result := variant.AddCombination(vars, n)
	if len(result) != len(expected) {
		t.Errorf("Expected %d combinations, but got %d", len(expected), len(result))
	}
	for i, r := range result {
		for j, v := range r {
			if v != expected[i][j] {
				t.Errorf("Expected %v, but got %v", expected[i], result[i])
				break
			}
		}
	}
}

func Test(t *testing.T) {
	var names []string

	variant.New(func(v *variant.Variant) {
		var name string

		v.Variants(variant.M{
			"first": func(v *variant.Variant) {
				name = "1"
			},
			"second": func(v *variant.Variant) {
				name = "2"
			},
		})

		v.Variants(variant.M{
			"first suffix": func(v *variant.Variant) {
				name += "A"
			},
			"second prefix": func(v *variant.Variant) {
				name += "B"
			},
		})

		names = append(names, name)

		t.Logf("name: %s", name)
	})

	got := strings.Join(names, ",")

	want := "1A,2A,1B,2B"
	if got != want {
		t.Errorf("expected names to be %q, got %q", want, got)
	}

	t.Fail()
}
