package string_set

type NextAction uint8

const (
	// Continue: keep iterating over the contents of the set
	Continue NextAction = iota
	// Stop: do not iterate over the remaining contents of the set
	Stop
)
