package string_set_insensitive

import (
	"github.com/wojnosystems/go-string-set/string_set"
	"strings"
)

const (
	defaultCapacity = 10
)

// Empty is a convenience declaration: it's an empty set you can use to compare
// to other sets if you want to use IsEqualTo instead of testing with Len
var Empty = NewWithCapacity(0)

// New creates a new String set, with a small default capacity. Case-insensitive
func New() *T {
	return NewWithCapacity(defaultCapacity)
}

// NewOf is a convenience method to create a string set containing the items you specify. Case-insensitive
func NewOf(items ...string) *T {
	ret := NewWithCapacity(len(items))
	ret.AddMany(items...)
	return ret
}

// NewWithCapacity creates a new, empty, string set with the provided capacity. Case-insensitive
func NewWithCapacity(capacity int) *T {
	return &T{
		T: string_set.NewWithCapacity(capacity),
	}
}

// convert changes the parameter into the value within the underlying storage
func convert(v string) string {
	return strings.ToLower(v)
}

// T holds the underlying string_set_insensitive type, do not instantiate this yourself,
// Please use New, NewOf, or NewWithCapacity
type T struct {
	*string_set.T
}

func (c *T) Add(v string) {
	c.T.Add(convert(v))
}

func (c *T) AddMany(v ...string) {
	for _, s := range v {
		c.Add(s)
	}
}

func (c *T) Remove(v string) {
	c.T.Remove(convert(v))
}

func (c *T) RemoveMany(v ...string) {
	for _, s := range v {
		c.Remove(s)
	}
}

func (c *T) Includes(v string) bool {
	return c.T.Includes(convert(v))
}

func (c *T) IsEqualTo(o string_set.Immutable) (equal bool) {
	// short-circuit test for speed
	if c.Len() != o.Len() {
		return false
	}
	equal = true
	o.EachCancelable(func(v string) string_set.NextAction {
		if !c.Includes(v) {
			equal = false
			return string_set.Break
		}
		return string_set.Continue
	})
	return
}

func (c *T) Union(o string_set.Immutable) (out string_set.Interface) {
	out = c.Copy()
	o.Each(func(v string) {
		out.Add(v)
	})
	return
}

func (c *T) Subtract(o string_set.Immutable) (out string_set.Interface) {
	out = NewWithCapacity(c.Len())
	c.Each(func(v string) {
		if !o.Includes(v) {
			out.Add(v)
		}
	})
	return
}

func (c *T) Intersection(o string_set.Immutable) (out string_set.Interface) {
	out = NewWithCapacity(c.Len())
	o.Each(func(v string) {
		if c.Includes(v) {
			out.Add(v)
		}
	})
	return
}

func (c *T) Any(item func(v string, converter func(in string) string) (didMatch bool)) (anyFound bool) {
	c.EachCancelable(func(v string) (a string_set.NextAction) {
		if item(strings.ToLower(v), convert) {
			anyFound = true
			return string_set.Break
		}
		return
	})
	return
}

func (c *T) None(item func(v string, converter func(in string) string) (didMatch bool)) (noneFound bool) {
	noneFound = true
	c.EachCancelable(func(v string) (a string_set.NextAction) {
		if item(strings.ToLower(v), convert) {
			noneFound = false
			return string_set.Break
		}
		return
	})
	return
}

func (c *T) Copy() string_set.Interface {
	outItems := NewWithCapacity(c.Len())
	c.T.Each(func(v string) {
		outItems.Add(v)
	})
	return outItems
}
