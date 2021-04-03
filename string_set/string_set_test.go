package string_set

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestCollection_Add(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected Immutable
	}{
		"empty": {
			input:    []string{},
			expected: Empty,
		},
		"one": {
			input:    []string{"a"},
			expected: NewOf("a"),
		},
		"two": {
			input:    []string{"a", "b"},
			expected: NewOf("a", "b"),
		},
		"two with duplicate": {
			input:    []string{"a", "b", "a"},
			expected: NewOf("a", "b"),
		},
		"three with duplicates": {
			input:    []string{"a", "a", "a"},
			expected: NewOf("a"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := NewOf(c.input...)
			assert.True(t, c.expected.IsEqualTo(actual))
		})
	}
}

func TestCollection_Remove(t *testing.T) {
	cases := map[string]struct {
		input    Interface
		rm       []string
		expected Immutable
	}{
		"empty": {
			input:    New(),
			rm:       []string{},
			expected: Empty,
		},
		"no items removed": {
			input:    NewOf("a", "b", "c"),
			rm:       []string{},
			expected: NewOf("a", "b", "c"),
		},
		"single item removed": {
			input:    NewOf("a", "b", "c"),
			rm:       []string{"b"},
			expected: NewOf("a", "c"),
		},
		"all removed": {
			input:    NewOf("a", "b", "c"),
			rm:       []string{"a", "b", "c"},
			expected: Empty,
		},
		"non-existent item removed": {
			input:    NewOf("a", "b", "c"),
			rm:       []string{"x"},
			expected: NewOf("a", "b", "c"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.input.RemoveMany(c.rm...)
			assert.True(t, c.expected.IsEqualTo(c.input))
		})
	}
}

func TestCollection_Includes(t *testing.T) {
	set := New()
	assert.False(t, set.Includes("something"))
	set.Add("something")
	assert.True(t, set.Includes("something"))
	assert.False(t, set.Includes("missing"))
}

func TestCollection_IsEmpty(t *testing.T) {
	set := New()
	assert.True(t, set.IsEmpty())
	assert.True(t, set.IsEqualTo(Empty))
	set.Add("something")
	assert.False(t, set.IsEmpty())
	assert.False(t, set.IsEqualTo(Empty))
}

func TestCollection_Len(t *testing.T) {
	set := New()
	assert.Equal(t, 0, set.Len())
	set.Add("something")
	assert.Equal(t, 1, set.Len())
	set.Add("something")
	assert.Equal(t, 1, set.Len())
	set.Add("something else")
	assert.Equal(t, 2, set.Len())
	set.Remove("something")
	assert.Equal(t, 1, set.Len())
}

func TestCollection_IsEqualTo(t *testing.T) {
	cases := map[string]struct {
		a        Immutable
		b        Immutable
		expected bool
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: true,
		},
		"a items, b empty": {
			a: NewOf("a", "b"),
			b: Empty,
		},
		"a empty, b items": {
			a: Empty,
			b: NewOf("a", "b"),
		},
		"equal same order": {
			a:        NewOf("a", "b", "c", "d"),
			b:        NewOf("a", "b", "c", "d"),
			expected: true,
		},
		"equal different order": {
			a:        NewOf("a", "b", "c", "d"),
			b:        NewOf("d", "c", "b", "a"),
			expected: true,
		},
		"same length, different items": {
			a: NewOf("a", "b"),
			b: NewOf("c", "d"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.IsEqualTo(c.b)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCollection_Each(t *testing.T) {
	input := NewOf("a", "b", "c", "d")

	actual := make([]string, 0, input.Len())
	input.Each(func(v string) {
		actual = append(actual, v)
	})

	expected := input.ToSlice()
	sort.Strings(expected)
	sort.Strings(actual)
	assert.Equal(t, expected, actual)
}

func TestCollection_EachCancelable(t *testing.T) {
	input := NewOf("a", "b", "c", "d")

	actual := make([]string, 0, input.Len())
	input.EachCancelable(func(v string) NextAction {
		if len(actual) > 1 {
			return Stop
		}
		actual = append(actual, v)
		return Continue
	})

	assert.Equal(t, 2, len(actual))
}

func TestCollection_Any(t *testing.T) {

	cases := map[string]struct {
		input    Iterator
		test     func(v string) bool
		expected bool
	}{
		"empty": {
			input: Empty,
			test: func(v string) bool {
				return true
			},
		},
		"nothing matches": {
			input: NewOf("a", "b"),
			test: func(v string) bool {
				return v == "c"
			},
		},
		"one thing matches": {
			input: NewOf("a", "b", "c"),
			test: func(v string) bool {
				return v == "b"
			},
			expected: true,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, c.expected, c.input.Any(c.test))
		})
	}
}

func TestCollection_None(t *testing.T) {
	cases := map[string]struct {
		input    Iterator
		test     func(v string) bool
		expected bool
	}{
		"empty": {
			input: Empty,
			test: func(v string) bool {
				return true
			},
			expected: true,
		},
		"nothing matches": {
			input: NewOf("a", "b"),
			test: func(v string) bool {
				return v == "c"
			},
			expected: true,
		},
		"one thing matches": {
			input: NewOf("a", "b", "c"),
			test: func(v string) bool {
				return v == "b"
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, c.expected, c.input.None(c.test))
		})
	}
}

func TestCollection_Copy(t *testing.T) {
	cases := map[string]struct {
		input Immutable
	}{
		"empty": {
			input: Empty,
		},
		"not empty": {
			input: NewOf("a", "b"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.input.Copy()
			assert.True(t, c.input.IsEqualTo(actual))
		})
	}
}

func TestCollection_ToSlice(t *testing.T) {
	cases := map[string]struct {
		input    Immutable
		expected []string
	}{
		"empty": {
			input:    Empty,
			expected: []string{},
		},
		"not empty": {
			input:    NewOf("a", "b"),
			expected: []string{"a", "b"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.input.ToSlice()
			sort.Strings(actual)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCollection_Union(t *testing.T) {
	cases := map[string]struct {
		a        Immutable
		b        Immutable
		expected Immutable
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: Empty,
		},
		"b empty": {
			a:        NewOf("a", "b"),
			b:        Empty,
			expected: NewOf("a", "b"),
		},
		"a empty": {
			a:        Empty,
			b:        NewOf("a", "b"),
			expected: NewOf("a", "b"),
		},
		"overlapping items": {
			a:        NewOf("a", "b"),
			b:        NewOf("a", "b"),
			expected: NewOf("a", "b"),
		},
		"unique items": {
			a:        NewOf("a", "b"),
			b:        NewOf("c", "d"),
			expected: NewOf("a", "b", "c", "d"),
		},
		"partial overlap items": {
			a:        NewOf("a", "b"),
			b:        NewOf("c", "b"),
			expected: NewOf("a", "b", "c"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.Union(c.b)
			assert.True(t, c.expected.IsEqualTo(actual))
		})
	}
}

func TestCollection_Subtract(t *testing.T) {
	cases := map[string]struct {
		a        Immutable
		b        Immutable
		expected Immutable
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: Empty,
		},
		"b empty": {
			a:        NewOf("a", "b"),
			b:        Empty,
			expected: NewOf("a", "b"),
		},
		"a empty": {
			a:        Empty,
			b:        NewOf("a", "b"),
			expected: Empty,
		},
		"overlapping items": {
			a:        NewOf("a", "b"),
			b:        NewOf("a", "b"),
			expected: Empty,
		},
		"unique items": {
			a:        NewOf("a", "b"),
			b:        NewOf("c", "d"),
			expected: NewOf("a", "b"),
		},
		"partial overlap items": {
			a:        NewOf("a", "b"),
			b:        NewOf("c", "b"),
			expected: NewOf("a"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.Subtract(c.b)
			assert.True(t, c.expected.IsEqualTo(actual))
		})
	}
}

func TestCollection_Intersection(t *testing.T) {
	cases := map[string]struct {
		a        Immutable
		b        Immutable
		expected Immutable
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: Empty,
		},
		"b empty": {
			a:        NewOf("a", "b"),
			b:        Empty,
			expected: Empty,
		},
		"a empty": {
			a:        Empty,
			b:        NewOf("a", "b"),
			expected: Empty,
		},
		"overlapping items": {
			a:        NewOf("a", "b"),
			b:        NewOf("a", "b"),
			expected: NewOf("a", "b"),
		},
		"unique items": {
			a:        NewOf("a", "b"),
			b:        NewOf("c", "d"),
			expected: Empty,
		},
		"partial overlap items": {
			a:        NewOf("a", "b"),
			b:        NewOf("b", "c"),
			expected: NewOf("b"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.Intersection(c.b)
			assert.True(t, c.expected.IsEqualTo(actual))
		})
	}
}
