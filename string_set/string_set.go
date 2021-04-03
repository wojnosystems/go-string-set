package string_set

const (
	defaultCapacity = 10
)

// Empty is a convenience declaration: it's an empty set you can use to compare
// to other sets if you want to use IsEqualTo instead of testing with Len
var Empty = (Immutable)(NewWithCapacity(0))

// New creates a new String set, with a small default capacity
func New() Interface {
	return NewWithCapacity(defaultCapacity)
}

// NewOf is a convenience method to create a string set containing the items you specify
func NewOf(items ...string) Interface {
	ret := NewWithCapacity(len(items))
	ret.AddMany(items...)
	return ret
}

// NewWithCapacity creates a new, empty, string set with the provided capacity
func NewWithCapacity(capacity int) Interface {
	return &collection{
		items: make(map[string]bool, capacity),
	}
}

type collection struct {
	items map[string]bool
}

func (c *collection) Add(v string) {
	c.items[v] = true
}

func (c *collection) AddMany(v ...string) {
	for _, s := range v {
		c.Add(s)
	}
}

func (c *collection) Remove(v string) {
	delete(c.items, v)
}

func (c *collection) RemoveMany(v ...string) {
	for _, s := range v {
		c.Remove(s)
	}
}

func (c *collection) Includes(v string) bool {
	_, ok := c.items[v]
	return ok
}

func (c *collection) IsEmpty() bool {
	return c.Len() == 0
}

func (c *collection) Len() int {
	return len(c.items)
}

func (c *collection) IsEqualTo(o Immutable) (equal bool) {
	// short-circuit test for speed
	if c.Len() != o.Len() {
		return false
	}
	equal = true
	c.EachCancelable(func(v string) NextAction {
		if !o.Includes(v) {
			equal = false
			return Stop
		}
		return Continue
	})
	return
}

func (c *collection) Union(o Immutable) (out Interface) {
	out = c.Copy()
	o.Each(func(v string) {
		out.Add(v)
	})
	return
}

func (c *collection) Subtract(o Immutable) (out Interface) {
	out = NewWithCapacity(c.Len())
	c.Each(func(v string) {
		if !o.Includes(v) {
			out.Add(v)
		}
	})
	return
}

func (c *collection) Intersection(o Immutable) (out Interface) {
	out = NewWithCapacity(c.Len())
	o.Each(func(v string) {
		if c.Includes(v) {
			out.Add(v)
		}
	})
	return
}

func (c *collection) ToSlice() (out []string) {
	out = make([]string, c.Len())
	i := 0
	for s := range c.items {
		out[i] = s
		i++
	}
	return
}

func (c *collection) Each(item func(v string)) {
	for value := range c.items {
		item(value)
	}
}

func (c *collection) EachCancelable(item func(v string) (next NextAction)) {
	for value := range c.items {
		action := item(value)
		if action == Stop {
			break
		}
	}
}

func (c *collection) Any(item func(v string) (didMatch bool)) (anyFound bool) {
	c.EachCancelable(func(v string) (a NextAction) {
		if item(v) {
			anyFound = true
			return Stop
		}
		return
	})
	return
}

func (c *collection) None(item func(v string) (didMatch bool)) (noneFound bool) {
	noneFound = true
	c.EachCancelable(func(v string) (a NextAction) {
		if item(v) {
			noneFound = false
			return Stop
		}
		return
	})
	return
}

func (c *collection) Copy() Interface {
	outItems := NewWithCapacity(c.Len())
	for s := range c.items {
		outItems.Add(s)
	}
	return outItems
}
