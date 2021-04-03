package string_set

func New() Interface {
	return &collection{
		items: make(map[string]bool),
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

func (c *collection) Equal(o Immutable) (equal bool) {
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
	out = New()
	c.Each(func(v string) {
		if !o.Exists(v) {
			out.Add(v)
		}
	})
	return
}

func (c *collection) Intersection(o Immutable) (out Interface) {
	out = New()
	o.Each(func(v string) {
		if c.Exists(v) {
			out.Add(v)
		}
	})
	return
}

func (c *collection) ToSlice() (out []string) {
	out = make([]string, 0, len(c.items))
	for s, _ := range c.items {
		out = append(out, s)
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
	outItems := make(map[string]bool, len(c.items))
	for s := range c.items {
		outItems[s] = true
	}
	return &collection{items: outItems}
}
