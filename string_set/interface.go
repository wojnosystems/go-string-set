package string_set

// Mutable allows you to alter the contents of the set
type Mutable interface {
	// Add an item to the set. If it already exists, this is just skipped and the item remains in the set
	Add(v string)

	// AddMany items to the set. Ignoring any duplicates
	AddMany(v ...string)

	// Remove an item from the set. If it doesn't exist, skip
	Remove(v string)

	// RemoveMany items from the set. Ignoring any non-existing items
	RemoveMany(v ...string)
}

// Tester contains read-only methods to query the metadata about the contents of the set
type Tester interface {
	// Includes returns true if the string is in the set, false if not found. Case-sensitive
	Includes(v string) bool

	// IsEmpty returns true if there are no items in the set, false if there is at least 1 item in the set
	IsEmpty() bool

	// Len returns the number of items in the set
	Len() int

	// IsEqualTo returns true both sets contain the same strings, false if they contain different numbers of items of
	// values of items differ. Sets don't care about item ordering, so you can add items to the sets in any order and
	// this will still be true.
	IsEqualTo(o Immutable) bool
}

// Iterator allows callers to loop over the contents of sets
type Iterator interface {
	// Each loops over each string in the set. The order is not guaranteed and can change between invocations
	Each(item func(v string))

	// EachCancelable is just like Each, but you can stop the iteration by returning
	// string_set.Break instead of string_set.Continue
	EachCancelable(item func(v string) (next NextAction))
}

// Copier allows new sets to be created from existing sets
type Copier interface {
	// Copy returns a mutable copy of the set. The returned set is shared-nothing, so you can treat this as a safe-copy
	Copy() Interface
}

// Slicer converts the set to a slice
type Slicer interface {
	// ToSlice returns a string slice with the contents of the set. There is no guarantee of the order the items will
	// be returned in. If you need them in lexical ordering, call sort.Strings() on the output of this method
	ToSlice() (out []string)
}

// Setter contains the set-specific methods
type Setter interface {
	// Union returns a new set containing all of the items from the callee and the parameter
	// union = left ∪ o
	Union(o Immutable) (out Interface)

	// Subtract returns a new set containing only items from the callee, but without the items in the parameter
	// subtracted = left - o
	Subtract(o Immutable) (out Interface)

	// Intersection returns a new set containing only items common to both the callee and parameter
	// intersection = left ∩ o
	Intersection(o Immutable) (out Interface)
}

// Immutable contains all of the read-only method calls that do not modify the set
type Immutable interface {
	Iterator
	Slicer
	Tester
	Setter
	Copier
}

// Interface contains all the methods for a set, Immutable and Mutable
type Interface interface {
	Mutable
	Immutable
}
