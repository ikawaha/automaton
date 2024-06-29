package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ikawaha/automaton"
)

const (
	s0 automaton.State = iota
	s1
	s2
	s3
	s4
	s5
	s6
	s7
	s8
	s9
	s10
	s11
	s12
	s13
	s14
	s15
	s16
	s17
	s18
	s19
	s20
	s21
)

var transition = automaton.Transition[string]{
	{State: s0, Input: "しか"}: s1, // o
	{State: s1, Input: "のこ"}: s2, // o
	{State: s1, Input: "しか"}: s1, // x
	{State: s2, Input: "のこ"}: s3, // o
	{State: s2, Input: "しか"}: s1, // x
	{State: s3, Input: "のこ"}: s4, // o
	{State: s3, Input: "しか"}: s1, // x
	{State: s4, Input: "こし"}: s5, // o
	{State: s4, Input: "しか"}: s1, // x
	{State: s5, Input: "たん"}: s6, // o
	{State: s5, Input: "しか"}: s1, // x
	{State: s6, Input: "たん"}: s7, // o

	{State: s7, Input: "しか"}:  s8,  // o
	{State: s8, Input: "のこ"}:  s9,  // o
	{State: s8, Input: "しか"}:  s1,  // x
	{State: s9, Input: "のこ"}:  s10, // o
	{State: s9, Input: "しか"}:  s1,  // x
	{State: s10, Input: "のこ"}: s11, // o
	{State: s10, Input: "しか"}: s1,  // x
	{State: s11, Input: "こし"}: s12, // o
	{State: s11, Input: "しか"}: s1,  // x
	{State: s12, Input: "たん"}: s13, // o
	{State: s12, Input: "しか"}: s1,  // x
	{State: s13, Input: "たん"}: s14, // o

	{State: s14, Input: "しか"}: s15, // o
	{State: s15, Input: "のこ"}: s16, // o
	{State: s15, Input: "しか"}: s1,  // x
	{State: s16, Input: "のこ"}: s17, // o
	{State: s16, Input: "しか"}: s1,  // x
	{State: s17, Input: "のこ"}: s18, // o
	{State: s17, Input: "しか"}: s1,  // x
	{State: s18, Input: "こし"}: s19, // o
	{State: s18, Input: "しか"}: s1,  // x
	{State: s19, Input: "たん"}: s20, // o
	{State: s19, Input: "しか"}: s1,  // x
	{State: s20, Input: "たん"}: s21, // o
}

var (
	start    = s0
	final    = automaton.NewStateSet(s21)
	fallback = s0
)

func OKInput(ctx context.Context, ch chan<- string) {
	v := []string{
		"しか", "のこ", "のこ", "のこ", "こし", "たん", "たん",
		"しか", "のこ", "のこ", "のこ", "たん", "たん",
		"しか", "のこ", "のこ", "のこ", "こし", "たん", "たん",
		"しか", "のこ", "のこ", "のこ", "こし", "たん", "たん",
	}
	for _, vv := range v {
		select {
		case <-ctx.Done():
			return
		default:
			ch <- vv
		}
	}
}

func RandomInput(ctx context.Context, ch chan<- string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	v := []string{"しか", "のこ", "こし", "たん"}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ch <- v[r.Intn(len(v))]
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	ch := make(chan string, 1)
	go RandomInput(ctx, ch)
	//go OKInput(ctx, ch)
	a := automaton.New(start, transition, final, automaton.Fallback[string](fallback))
	start := time.Now()
	ok, err := a.Run(ctx, ch, func(_ automaton.State, in string, q automaton.State) {
		fmt.Printf(in)
		if a.Final(q) {
			fmt.Println(" 🦌🦌🦌")
		}
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	d := time.Since(start)
	time.Sleep(1 * time.Second)
	fmt.Println("\n", d)
	if !ok {
		os.Exit(1)
	}
}
