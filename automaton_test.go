package automaton

import (
	"bytes"
	"testing"
)

func TestAutomaton_Dot(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var a Automaton[byte]
		want := `digraph finite_state_automaton {
  rankdir = LR;
  fontname="sans-serif";
  node [shape=circle, fontname="sans-serif"];
  edge [fontname="sans-serif"];
  // start
  start [shape=point];
  start -> "0";
}
`
		var b bytes.Buffer
		if err := a.Dot(&b); err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		if got := b.String(); want != got {
			t.Errorf("want %s, got %s", want, got)
		}
	})
	t.Run("State x string -> State", func(t *testing.T) {
		transition := Transition[string]{
			InputPair[string]{State: 1, Input: "a"}: 2,
		}
		a := Automaton[string]{
			start:      0,
			transition: transition,
			finalState: NewStateSet(2),
		}
		want := `digraph finite_state_automaton {
  rankdir = LR;
  fontname="sans-serif";
  node [shape=circle, fontname="sans-serif"];
  edge [fontname="sans-serif"];
  // transition
  "1" -> "2" [label="a"];
  // start
  start [shape=point];
  start -> "0";
  // final states
  "2" [shape=ellipse, peripheries=2];
}
`
		var b bytes.Buffer
		if err := a.Dot(&b); err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		if got := b.String(); want != got {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("ddsk", func(t *testing.T) {
		transition := Transition[string]{
			{State: 0, Input: "ドド"}:  1,  // o
			{State: 0, Input: "スコ"}:  0,  // x
			{State: 1, Input: "スコ"}:  2,  // o
			{State: 2, Input: "スコ"}:  3,  // o
			{State: 3, Input: "スコ"}:  4,  // o
			{State: 4, Input: "ドド"}:  5,  // o
			{State: 4, Input: "スコ"}:  0,  // x
			{State: 5, Input: "スコ"}:  6,  // o
			{State: 6, Input: "スコ"}:  7,  // o
			{State: 7, Input: "スコ"}:  8,  // o
			{State: 8, Input: "ドド"}:  9,  // o
			{State: 8, Input: "スコ"}:  0,  // x
			{State: 9, Input: "スコ"}:  10, // o
			{State: 10, Input: "スコ"}: 11, // o
			{State: 11, Input: "スコ"}: 12, // o
		}
		a := Automaton[string]{
			start:      0,
			transition: transition,
			finalState: NewStateSet(12),
		}
		want := `digraph finite_state_automaton {
  rankdir = LR;
  fontname="sans-serif";
  node [shape=circle, fontname="sans-serif"];
  edge [fontname="sans-serif"];
  // transition
  "0" -> "0" [label="スコ"];
  "0" -> "1" [label="ドド"];
  "1" -> "2" [label="スコ"];
  "2" -> "3" [label="スコ"];
  "3" -> "4" [label="スコ"];
  "4" -> "0" [label="スコ"];
  "4" -> "5" [label="ドド"];
  "5" -> "6" [label="スコ"];
  "6" -> "7" [label="スコ"];
  "7" -> "8" [label="スコ"];
  "8" -> "0" [label="スコ"];
  "8" -> "9" [label="ドド"];
  "9" -> "10" [label="スコ"];
  "10" -> "11" [label="スコ"];
  "11" -> "12" [label="スコ"];
  // start
  start [shape=point];
  start -> "0";
  // final states
  "12" [shape=ellipse, peripheries=2];
}
`
		var b bytes.Buffer
		if err := a.Dot(&b); err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		if got := b.String(); want != got {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}
