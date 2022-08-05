package automaton

// State represents a state of automata.
type State int

// StateSet represents a set of states.
type StateSet map[State]struct{}

// NewStateSet creates a set of states.
func NewStateSet(s ...State) StateSet {
	ret := make(StateSet, len(s))
	for _, v := range s {
		ret[v] = struct{}{}
	}
	return ret
}
