package automaton

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

// Option sets an option to the automaton.
type Option[T comparable] func(*Automaton[T])

// Fallback sets the fallback state when no transitions are not defined for a given state and input.
func Fallback[T comparable](p State) Option[T] {
	return func(a *Automaton[T]) {
		a.fallback = &p
	}
}

func Loop[T comparable]() Option[T] {
	return func(a *Automaton[T]) {
		a.loop = true
	}
}

// Automaton represents an automaton.
type Automaton[T comparable] struct {
	start      State
	transition Transition[T]
	finalState StateSet
	fallback   *State
	loop       bool
}

// New generates an automaton with the initial state is start, the transition function is transition, the set of final states is final and several options.
func New[T comparable](start State, transition Transition[T], final StateSet, opts ...Option[T]) *Automaton[T] {
	ret := &Automaton[T]{
		start:      start,
		transition: transition,
		finalState: final,
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

// InitialState returns the initial state of the automaton.
func (a *Automaton[T]) InitialState() State {
	return a.start
}

// Next transitions from the state 'p' and the input 'in' to the next state.
func (a *Automaton[T]) Next(p State, in T) (State, bool) {
	q, ok := a.transition[InputPair[T]{State: p, Input: in}]
	return q, ok
}

// Final returns true if the state is in the set of the final states.
func (a *Automaton[T]) Final(p State) bool {
	_, ok := a.finalState[p]
	return ok
}

// Run executes the automaton while reading input from the channel. If you want
// to insert some action while the automaton is running, specify an insert
// function. If the state of the automaton goes to the accepting states, it
// returns true without reading the rest of the input. If there is no more
// input (the channel is closed), returns false.
// The automaton returns an error when there are no more transitions.
func (a *Automaton[T]) Run(ctx context.Context, in <-chan T, insertion func(p State, in T, q State)) (bool, error) {
	p := a.InitialState()
	for {
		select {
		case u, ok := <-in:
			if !ok {
				return false, nil
			}
			q, ok := a.Next(p, u)
			if !ok {
				if a.fallback == nil {
					return false, fmt.Errorf("undefined transition, state: %v, input: %v", p, u)
				}
				q = *a.fallback
			}
			insertion(p, u, q)
			if a.Final(q) && !a.loop {
				return true, nil
			}
			p = q
		case <-ctx.Done():
			return false, nil
		}
	}
}

//go:embed graph.tpl
var graphT string
var graphTpl = template.Must(template.New("fsa").Parse(graphT))

func (a *Automaton[T]) Dot(w io.Writer) error {
	return graphTpl.Execute(w, struct {
		Start      State
		Transition Transition[T]
		FinalState StateSet
	}{
		Start:      a.start,
		Transition: a.transition,
		FinalState: a.finalState,
	})
}
