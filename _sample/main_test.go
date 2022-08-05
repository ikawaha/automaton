package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/ikawaha/automaton"
)

func Test_RunWithoutOptions(t *testing.T) {
	testdata := []struct {
		name    string
		input   []string
		want    bool
		wantErr error
		output  string
	}{
		{
			name:    "exact match",
			input:   []string{"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ"},
			want:    true,
			wantErr: nil,
			output:  "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name: "postfix match",
			input: []string{
				"スコ", "スコ", "ドド", "スコ", "スコ", "ドド",
				"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ"},
			want:    false,
			wantErr: errors.New("undefined transition, state: 3, input: ドド"),
			output:  "スコスコドドスコスコ",
		},
		{
			name: "prefix match",
			input: []string{
				"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ",
				"スコ", "スコ", "ドド",
			},
			want:    true,
			wantErr: nil,
			output:  "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name:    "not accept",
			input:   []string{"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ"},
			want:    false,
			wantErr: nil,
			output:  "ドドスコスコスコドドスコスコスコドドスコスコ",
		},
	}
	a := automaton.New(start, transition, final)
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			ch := make(chan string, 1)
			go func(ctx context.Context, input []string) {
				select {
				default:
					for _, v := range input {
						ch <- v
					}
					close(ch)
				case <-ctx.Done():
					close(ch)
				}
			}(ctx, tt.input)
			var b bytes.Buffer
			got, err := a.Run(context.TODO(), ch, func(p automaton.State, in string, q automaton.State) {
				b.WriteString(in)
				if a.Final(q) {
					b.WriteString("ラブ注入♡")
				}
			})
			if err == nil && tt.wantErr != nil || err != nil && tt.wantErr == nil {
				t.Errorf("want error %v, but got %v", tt.wantErr, err)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("want error %v, but got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("want %t, but got %t", tt.want, got)
			}
			if got, want := b.String(), tt.output; got != want {
				t.Errorf("want %s, got %s", want, got)
			}
		})
	}
}

func Test_RunWithFallback(t *testing.T) {
	testdata := []struct {
		name   string
		input  []string
		want   bool
		output string
	}{
		{
			name:   "exact match",
			input:  []string{"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ"},
			want:   false,
			output: "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name: "postfix match",
			input: []string{
				"スコ", "スコ", "ドド", "スコ", "スコ", "ドド",
				"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ"},
			want:   false,
			output: "スコスコドドスコスコドド" + "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name: "prefix match",
			input: []string{
				"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ",
				"スコ", "スコ", "ドド",
			},
			want:   true,
			output: "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name:   "not accept",
			input:  []string{"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "ドド"},
			want:   false,
			output: "ドドスコスコスコドドスコスコスコドドスコスコドド",
		},
	}
	a := automaton.New(start, transition, final, automaton.Fallback[string](fallback))
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			ch := make(chan string, 1)
			go func(ctx context.Context, input []string) {
				select {
				default:
					for _, v := range input {
						ch <- v
					}
					close(ch)
				case <-ctx.Done():
					close(ch)
				}
			}(ctx, tt.input)
			var b bytes.Buffer
			got, err := a.Run(context.TODO(), ch, func(p automaton.State, in string, q automaton.State) {
				b.WriteString(in)
				if a.Final(q) {
					b.WriteString("ラブ注入♡")
				}
			})
			if err != nil {
				t.Fatalf("unexpected error, %v", err)
			}
			if got != tt.want {
				t.Errorf("want %t, but got %t", tt.want, got)
			}
			if got, want := b.String(), tt.output; got != want {
				t.Errorf("want %s, got %s", want, got)
			}
		})
	}
}

func Test_RunWithFallbackLoop(t *testing.T) {
	testdata := []struct {
		name   string
		input  []string
		want   bool
		output string
	}{
		{
			name:   "exact match",
			input:  []string{"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ"},
			want:   false,
			output: "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name: "postfix match",
			input: []string{
				"スコ", "スコ", "ドド", "スコ", "スコ", "ドド",
				"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ"},
			want:   false,
			output: "スコスコドドスコスコドド" + "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡",
		},
		{
			name: "prefix match",
			input: []string{
				"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ",
				"スコ", "スコ", "ドド",
			},
			want:   false,
			output: "ドドスコスコスコドドスコスコスコドドスコスコスコラブ注入♡" + "スコスコドド",
		},
		{
			name:   "not accept",
			input:  []string{"ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "スコ", "ドド", "スコ", "スコ", "ドド"},
			want:   false,
			output: "ドドスコスコスコドドスコスコスコドドスコスコドド",
		},
	}
	a := automaton.New(start, transition, final, automaton.Fallback[string](fallback), automaton.Loop[string]())
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			ch := make(chan string, 1)
			go func(ctx context.Context, input []string) {
				select {
				default:
					for _, v := range input {
						ch <- v
					}
					close(ch)
				case <-ctx.Done():
					close(ch)
				}
			}(ctx, tt.input)
			var b bytes.Buffer
			got, err := a.Run(context.TODO(), ch, func(p automaton.State, in string, q automaton.State) {
				b.WriteString(in)
				if a.Final(q) {
					b.WriteString("ラブ注入♡")
				}
			})
			if err != nil {
				t.Fatalf("unexpected error, %v", err)
			}
			if got != tt.want {
				t.Errorf("want %t, but got %t", tt.want, got)
			}
			if got, want := b.String(), tt.output; got != want {
				t.Errorf("want %s, got %s", want, got)
			}
		})
	}
}
