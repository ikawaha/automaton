package automaton

type InputPair[T comparable] struct {
	State State
	Input T
}

type Transition[T comparable] map[InputPair[T]]State
