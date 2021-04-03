package string_set_insensitive

import (
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-string-set/string_set"
	"testing"
)

func TestCollection_Add(t *testing.T) {
	cases := map[string]struct {
		set       string_set.Immutable
		existTest string
		expected  bool
	}{
		"empty": {
			set:       Empty,
			existTest: "missing",
		},
		"case is the same": {
			set:       NewOf("a"),
			existTest: "a",
			expected:  true,
		},
		"case is different exists": {
			set:       NewOf("a"),
			existTest: "A",
			expected:  true,
		},
		"case is different when added": {
			set:       NewOf("A"),
			existTest: "a",
			expected:  true,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, c.expected, c.set.Includes(c.existTest))
		})
	}
}

func TestCollection_Remove(t *testing.T) {
	cases := map[string]struct {
		set      string_set.Interface
		remove   []string
		expected string_set.Immutable
	}{
		"empty": {
			set:      Empty,
			remove:   []string{"missing"},
			expected: Empty,
		},
		"case is the same": {
			set:      NewOf("a"),
			remove:   []string{"a"},
			expected: Empty,
		},
		"case is different": {
			set:      NewOf("a"),
			remove:   []string{"A"},
			expected: Empty,
		},
		"multiple case is different": {
			set:      NewOf("a", "B", "c", "D"),
			remove:   []string{"A", "b", "C", "d"},
			expected: Empty,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.set.RemoveMany(c.remove...)
			assert.True(t, c.set.IsEqualTo(c.expected))
		})
	}
}

func TestCollection_IsEmpty(t *testing.T) {
	set := New()
	assert.True(t, set.IsEmpty())
	set.Add("something")
	assert.False(t, set.IsEmpty())
}

func TestCollection_IsEqualTo(t *testing.T) {
	cases := map[string]struct {
		a        string_set.Immutable
		b        string_set.Immutable
		expected bool
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: true,
		},
		"case is the same": {
			a:        NewOf("a"),
			b:        NewOf("a"),
			expected: true,
		},
		"case is different": {
			a:        NewOf("a"),
			b:        NewOf("A"),
			expected: true,
		},
		"multiple case is different": {
			a:        NewOf("a", "B", "c", "D"),
			b:        NewOf("A", "b", "C", "d"),
			expected: true,
		},
		"unequal lengths": {
			a: NewOf("a", "B", "c", "D"),
			b: NewOf("A", "b"),
		},
		"equal lengths with different values": {
			a: NewOf("c", "D"),
			b: NewOf("A", "b"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, c.expected, c.a.IsEqualTo(c.b))
		})
	}
}

func TestCollection_Any(t *testing.T) {
	cases := map[string]struct {
		input    *T
		test     func(v string, convert func(in string) string) bool
		expected bool
	}{
		"empty": {
			input: Empty,
			test: func(v string, convert func(in string) string) bool {
				return true
			},
		},
		"nothing matches": {
			input: NewOf("a", "b"),
			test: func(v string, convert func(in string) string) bool {
				return v == convert("c")
			},
		},
		"one thing matches": {
			input: NewOf("a", "b", "c"),
			test: func(v string, convert func(in string) string) bool {
				return v == convert("B")
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
		input    *T
		test     func(v string, convert func(in string) string) bool
		expected bool
	}{
		"empty": {
			input: Empty,
			test: func(v string, convert func(in string) string) bool {
				return true
			},
			expected: true,
		},
		"nothing matches": {
			input: NewOf("a", "b"),
			test: func(v string, convert func(in string) string) bool {
				return v == convert("c")
			},
			expected: true,
		},
		"one thing matches": {
			input: NewOf("a", "b", "c"),
			test: func(v string, convert func(in string) string) bool {
				return v == convert("B")
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, c.expected, c.input.None(c.test))
		})
	}
}

func TestCollection_Union(t *testing.T) {
	cases := map[string]struct {
		a        string_set.Immutable
		b        string_set.Immutable
		expected string_set.Immutable
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: Empty,
		},
		"b empty": {
			a:        NewOf("a", "B"),
			b:        Empty,
			expected: NewOf("a", "b"),
		},
		"a empty": {
			a:        Empty,
			b:        NewOf("a", "B"),
			expected: NewOf("a", "b"),
		},
		"overlapping items": {
			a:        NewOf("a", "B"),
			b:        NewOf("A", "b"),
			expected: NewOf("a", "b"),
		},
		"unique items": {
			a:        NewOf("a", "B"),
			b:        NewOf("C", "d"),
			expected: NewOf("a", "b", "c", "d"),
		},
		"partial overlap items": {
			a:        NewOf("A", "b"),
			b:        NewOf("c", "B"),
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
		a        string_set.Immutable
		b        string_set.Immutable
		expected string_set.Immutable
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
			b:        NewOf("a", "B"),
			expected: Empty,
		},
		"overlapping items": {
			a:        NewOf("A", "b"),
			b:        NewOf("a", "B"),
			expected: Empty,
		},
		"unique items": {
			a:        NewOf("A", "b"),
			b:        NewOf("c", "D"),
			expected: NewOf("a", "b"),
		},
		"partial overlap items": {
			a:        NewOf("A", "b"),
			b:        NewOf("c", "B"),
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
		a        string_set.Immutable
		b        string_set.Immutable
		expected string_set.Immutable
	}{
		"empty": {
			a:        Empty,
			b:        Empty,
			expected: Empty,
		},
		"b empty": {
			a:        NewOf("A", "b"),
			b:        Empty,
			expected: Empty,
		},
		"a empty": {
			a:        Empty,
			b:        NewOf("A", "b"),
			expected: Empty,
		},
		"overlapping items": {
			a:        NewOf("A", "b"),
			b:        NewOf("a", "B"),
			expected: NewOf("a", "b"),
		},
		"unique items": {
			a:        NewOf("A", "b"),
			b:        NewOf("c", "D"),
			expected: Empty,
		},
		"partial overlap items": {
			a:        NewOf("A", "b"),
			b:        NewOf("B", "c"),
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
