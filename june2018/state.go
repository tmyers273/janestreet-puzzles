package june2018

import "strconv"

type State int

const (
	StateValid = iota + 1
	StateInvalid
	StateIndeterminate
)

func (s State) String() string {
	switch s {
	case StateValid:
		return "Valid"
	case StateInvalid:
		return "Invalid"
	case StateIndeterminate:
		return "Unknown"
	default:
		panic("Unknown state: " + strconv.Itoa(int(s)))
	}
}
