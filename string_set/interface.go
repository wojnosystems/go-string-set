package string_set

type Mutable interface {
	Add(v string)
	Remove(v string)
}

type Tester interface {
	Exists(v string) bool
	IsEmpty() bool
	Len() int
	Equal(o Immutable) bool
}

type Iterator interface {
	Each(item func(v string))
	EachCancelable(item func(v string) (next NextAction))
	Any(item func(v string) (didMatch bool)) bool
	None(item func(v string) (didMatch bool)) bool
}

type Copier interface {
	Copy() Interface
}

type Slicer interface {
	ToSlice() (out []string)
}

type Setter interface {
	Union(o Immutable) (out Interface)
	Subtract(o Immutable) (out Interface)
	Intersection(o Immutable) (out Interface)
}

type Immutable interface {
	Iterator
	Slicer
	Tester
	Setter
	Copier
}

type Interface interface {
	Mutable
	Immutable
}
