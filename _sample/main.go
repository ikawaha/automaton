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
)

var transition = automaton.Transition[string]{
	{State: s0, Input: "ドド"}: s1, // o
	{State: s0, Input: "スコ"}: s0, // x
	//{State: s1, Input: "ドド"}: s1, // x
	{State: s1, Input: "スコ"}: s2, // o
	//{State: s2, Input: "ドド"}: s1, // x
	{State: s2, Input: "スコ"}: s3, // o
	//{State: s3, Input: "ドド"}: s1, // x
	{State: s3, Input: "スコ"}: s4, // o

	{State: s4, Input: "ドド"}: s5, // o
	{State: s4, Input: "スコ"}: s0, // x
	//{State: s5, Input: "ドド"}: s1, // x
	{State: s5, Input: "スコ"}: s6, // o
	//{State: s6, Input: "ドド"}: s1, // x
	{State: s6, Input: "スコ"}: s7, // o
	//{State: s7, Input: "ドド"}: s1, // x
	{State: s7, Input: "スコ"}: s8, // o

	{State: s8, Input: "ドド"}: s9, // o
	{State: s8, Input: "スコ"}: s0, // x
	//{State: s9, Input: "ドド"}:  s1,  // x
	{State: s9, Input: "スコ"}: s10, // o
	//{State: s10, Input: "ドド"}: s1,  // x
	{State: s10, Input: "スコ"}: s11, // o
	//{State: s11, Input: "ドド"}: s1,  // x
	{State: s11, Input: "スコ"}: s12, // o
}

var (
	start    = s0
	final    = automaton.NewStateSet(s12)
	fallback = s1
)

func RandomInput(ctx context.Context, ch chan<- string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	v := []string{"ドド", "スコ"}
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ch := make(chan string, 1)
	go RandomInput(ctx, ch)

	a := automaton.New(start, transition, final, automaton.Fallback[string](fallback))
	ok, err := a.Run(ctx, ch, func(_ automaton.State, in string, q automaton.State) {
		fmt.Printf(in)
		if a.Final(q) {
			fmt.Printf("ラブ注入♡")
		}
	})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		os.Exit(1)
	}
}
