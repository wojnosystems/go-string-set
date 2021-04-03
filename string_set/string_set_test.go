package string_set

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestCollection_Add(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected []string
	}{
		"empty": {
			input:    []string{},
			expected: []string{},
		},
		"one": {
			input:    []string{"a"},
			expected: []string{"a"},
		},
		"two": {
			input:    []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		"two with duplicate": {
			input:    []string{"a", "b", "a"},
			expected: []string{"a", "b"},
		},
		"three with duplicates": {
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := populateSetResult(func(set Interface) {
				for _, item := range c.input {
					set.Add(item)
				}
			})
			assert.Equal(t, c.expected, actual)
		})
	}
}

func populateSetResult(cb func(set Interface)) []string {
	working := New()
	cb(working)
	actual := working.ToSlice()
	sort.Strings(actual)
	return actual
}

func TestCollection_Remove(t *testing.T) {
	cases := map[string]struct {
		input    []string
		rm       []string
		expected []string
	}{
		"empty": {
			input:    []string{},
			rm:       []string{},
			expected: []string{},
		},
		"no items removed": {
			input:    []string{"a", "b", "c"},
			rm:       []string{},
			expected: []string{"a", "b", "c"},
		},
		"single item removed": {
			input:    []string{"a", "b", "c"},
			rm:       []string{"b"},
			expected: []string{"a", "c"},
		},
		"all removed": {
			input:    []string{"a", "b", "c"},
			rm:       []string{"a", "b", "c"},
			expected: []string{},
		},
		"non-existant item removed": {
			input:    []string{"a", "b", "c"},
			rm:       []string{"x"},
			expected: []string{"a", "b", "c"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := populateSetResult(func(set Interface) {
				for _, item := range c.input {
					set.Add(item)
				}

				for _, item := range c.rm {
					set.Remove(item)
				}
			})
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCollection_Exists(t *testing.T) {
	set := New()
	assert.False(t, set.Exists("something"))
	set.Add("something")
	assert.True(t, set.Exists("something"))
	assert.False(t, set.Exists("missing"))
}

func TestCollection_IsEmpty(t *testing.T) {
	set := New()
	assert.True(t, set.IsEmpty())
	set.Add("something")
	assert.False(t, set.IsEmpty())
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

func TestCollection_Equal(t *testing.T) {
	cases := map[string]struct {
		a        []string
		b        []string
		expected bool
	}{
		"empty": {
			a:        []string{},
			b:        []string{},
			expected: true,
		},
		"a items, b empty": {
			a: []string{"a", "b"},
			b: []string{},
		},
		"a empty, b items": {
			a: []string{},
			b: []string{"a", "b"},
		},
		"equal same order": {
			a:        []string{"a", "b", "c", "d"},
			b:        []string{"a", "b", "c", "d"},
			expected: true,
		},
		"equal different order": {
			a:        []string{"a", "b", "c", "d"},
			b:        []string{"d", "c", "b", "a"},
			expected: true,
		},
		"same length, different items": {
			a: []string{"a", "b"},
			b: []string{"c", "d"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			a := New()
			for _, s := range c.a {
				a.Add(s)
			}
			b := New()
			for _, s := range c.b {
				b.Add(s)
			}
			actual := a.Equal(b)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCollection_Each(t *testing.T) {
	input := []string{"a", "b", "c", "d"}

	working := New()
	for _, s := range input {
		working.Add(s)
	}

	actual := make([]string, 0, working.Len())
	working.Each(func(v string) {
		actual = append(actual, v)
	})

	sort.Strings(actual)
	assert.Equal(t, input, actual)
}

func TestCollection_EachCancelable(t *testing.T) {
	input := []string{"a", "b", "c", "d"}

	working := New()
	for _, s := range input {
		working.Add(s)
	}

	actual := make([]string, 0, working.Len())
	working.EachCancelable(func(v string) NextAction {
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
		input    []string
		test     func(v string) bool
		expected bool
	}{
		"empty": {
			input: []string{},
			test: func(v string) bool {
				return true
			},
		},
		"nothing matches": {
			input: []string{"a", "b"},
			test: func(v string) bool {
				return v == "c"
			},
		},
		"one thing matches": {
			input: []string{"a", "b", "c"},
			test: func(v string) bool {
				return v == "b"
			},
			expected: true,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := New()
			for _, s := range c.input {
				actual.Add(s)
			}
			assert.Equal(t, c.expected, actual.Any(c.test))
		})
	}
}

func TestCollection_None(t *testing.T) {
	cases := map[string]struct {
		input    []string
		test     func(v string) bool
		expected bool
	}{
		"empty": {
			input: []string{},
			test: func(v string) bool {
				return true
			},
			expected: true,
		},
		"nothing matches": {
			input: []string{"a", "b"},
			test: func(v string) bool {
				return v == "c"
			},
			expected: true,
		},
		"one thing matches": {
			input: []string{"a", "b", "c"},
			test: func(v string) bool {
				return v == "b"
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := New()
			for _, s := range c.input {
				actual.Add(s)
			}
			assert.Equal(t, c.expected, actual.None(c.test))
		})
	}
}

func TestCollection_Copy(t *testing.T) {
	cases := map[string]struct {
		input []string
	}{
		"empty": {
			input: []string{},
		},
		"not empty": {
			input: []string{"a", "b"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			working := New()
			for _, s := range c.input {
				working.Add(s)
			}
			actual := working.Copy()
			assert.True(t, working.Equal(actual))
		})
	}
}

func TestCollection_ToSlice(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected []string
	}{
		"empty": {
			input:    []string{},
			expected: []string{},
		},
		"not empty": {
			input:    []string{"a", "b"},
			expected: []string{"a", "b"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := populateSetResult(func(set Interface) {
				for _, s := range c.input {
					set.Add(s)
				}
			})
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
			a:        New(),
			b:        New(),
			expected: New(),
		},
		"b empty": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: New(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
		},
		"a empty": {
			a: New(),
			b: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
		},
		"overlapping items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
		},
		"unique items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("c")
				s.Add("d")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				s.Add("c")
				s.Add("d")
				return s
			}(),
		},
		"partial overlap items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("b")
				s.Add("c")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				s.Add("c")
				return s
			}(),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.Union(c.b)
			assert.True(t, c.expected.Equal(actual))
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
			a:        New(),
			b:        New(),
			expected: New(),
		},
		"b empty": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: New(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
		},
		"a empty": {
			a: New(),
			b: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			expected: New(),
		},
		"overlapping items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			expected: New(),
		},
		"unique items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("c")
				s.Add("d")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
		},
		"partial overlap items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("b")
				s.Add("c")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				return s
			}(),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.Subtract(c.b)
			assert.True(t, c.expected.Equal(actual))
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
			a:        New(),
			b:        New(),
			expected: New(),
		},
		"b empty": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b:        New(),
			expected: New(),
		},
		"a empty": {
			a: New(),
			b: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			expected: New(),
		},
		"overlapping items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
		},
		"unique items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("c")
				s.Add("d")
				return s
			}(),
			expected: New(),
		},
		"partial overlap items": {
			a: func() Immutable {
				s := New()
				s.Add("a")
				s.Add("b")
				return s
			}(),
			b: func() Immutable {
				s := New()
				s.Add("b")
				s.Add("c")
				return s
			}(),
			expected: func() Immutable {
				s := New()
				s.Add("b")
				return s
			}(),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.a.Intersection(c.b)
			assert.True(t, c.expected.Equal(actual))
		})
	}
}
