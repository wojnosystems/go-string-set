package string_set

const (
	defaultCapacity = 10
)

// Empty is a convenience declaration: it's an empty set you can use to compare
// to other sets if you want to use IsEqualTo instead of testing with Len
var Empty = NewWithCapacity(0)

// New creates a new String set, with a small default capacity
func New() *T {
	return NewWithCapacity(defaultCapacity)
}

// NewOf is a convenience method to create a string set containing the items you specify
func NewOf(items ...string) *T {
	ret := NewWithCapacity(len(items))
	ret.AddMany(items...)
	return ret
}

// NewWithCapacity creates a new, empty, string set with the provided capacity
func NewWithCapacity(capacity int) *T {
	return &T{
		items: make(map[string]bool, capacity),
	}
}

// T holds the underlying string_set type, do not instantiate this yourself,
// Please use New, NewOf, or NewWithCapacity
type T struct {
	items map[string]bool
}

func (c *T) Add(v string) {
	c.items[v] = true
}

func (c *T) AddMany(v ...string) {
	for _, s := range v {
		c.Add(s)
	}
}

func (c *T) Remove(v string) {
	delete(c.items, v)
}

func (c *T) RemoveMany(v ...string) {
	for _, s := range v {
		c.Remove(s)
	}
}

func (c *T) Includes(v string) bool {
	_, ok := c.items[v]
	return ok
}

func (c *T) IsEmpty() bool {
	return c.Len() == 0
}

func (c *T) Len() int {
	return len(c.items)
}

func (c *T) IsEqualTo(o Immutable) (equal bool) {
	// short-circuit test for speed
	if c.Len() != o.Len() {
		return false
	}
	equal = true
	c.EachCancelable(func(v string) NextAction {
		if !o.Includes(v) {
			equal = false
			return Break
		}
		return Continue
	})
	return
}

func (c *T) Union(o Immutable) (out Interface) {
	out = c.Copy()
	o.Each(func(v string) {
		out.Add(v)
	})
	return
}

func (c *T) Subtract(o Immutable) (out Interface) {
	out = NewWithCapacity(c.Len())
	c.Each(func(v string) {
		if !o.Includes(v) {
			out.Add(v)
		}
	})
	return
}

func (c *T) Intersection(o Immutable) (out Interface) {
	out = NewWithCapacity(c.Len())
	o.Each(func(v string) {
		if c.Includes(v) {
			out.Add(v)
		}
	})
	return
}

func (c *T) ToSlice() (out []string) {
	out = make([]string, c.Len())
	i := 0
	for s := range c.items {
		out[i] = s
		i++
	}
	return
}

func (c *T) Each(item func(v string)) {
	for value := range c.items {
		item(value)
	}
}

func (c *T) EachCancelable(item func(v string) (next NextAction)) {
	for value := range c.items {
		action := item(value)
		if action == Break {
			break
		}
	}
}

// Any returns true if predicate returns true for any item. Short-circuits and stops iteration when didMatch
// returns true. Returns false if no item caused predicate to return true
func (c *T) Any(item func(v string) (didMatch bool)) (anyFound bool) {
	c.EachCancelable(func(v string) (a NextAction) {
		if item(v) {
			anyFound = true
			return Break
		}
		return
	})
	return
}

// None return true if predicate returned false for every item in the set. If predicate returns true, short-circuit
// and return false from this method, indicating that at least 1 item matched
func (c *T) None(item func(v string) (didMatch bool)) (noneFound bool) {
	noneFound = true
	c.EachCancelable(func(v string) (a NextAction) {
		if item(v) {
			noneFound = false
			return Break
		}
		return
	})
	return
}

func (c *T) Copy() Interface {
	outItems := NewWithCapacity(c.Len())
	for s := range c.items {
		outItems.Add(s)
	}
	return outItems
}
