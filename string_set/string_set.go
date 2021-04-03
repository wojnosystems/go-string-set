package string_set

const (
	defaultCapacity = 10
)

func New() Interface {
	return NewWithCapacity(defaultCapacity)
}

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

func (c *collection) Remove(v string) {
	delete(c.items, v)
}

func (c *collection) Exists(v string) bool {
	_, ok := c.items[v]
	return ok
}

func (c *collection) IsEmpty() bool {
	return c.Len() == 0
}

func (c *collection) Len() int {
	return len(c.items)
}

func (c *collection) IsEqual(o Immutable) (equal bool) {
	// short-circuit test for speed
	if c.Len() != o.Len() {
		return false
	}
	equal = true
	c.EachCancelable(func(v string) NextAction {
		if !o.Exists(v) {
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
		if !o.Exists(v) {
			out.Add(v)
		}
	})
	return
}

func (c *collection) Intersection(o Immutable) (out Interface) {
	out = NewWithCapacity(c.Len())
	o.Each(func(v string) {
		if c.Exists(v) {
			out.Add(v)
		}
	})
	return
}

func (c *collection) ToSlice() (out []string) {
	out = make([]string, c.Len())
	i := 0
	for s, _ := range c.items {
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
