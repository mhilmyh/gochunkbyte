package gochunkbyte

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-test/deep"
)

type testCase struct {
	name    string
	execute func() interface{}
	want    func() interface{}
}

func runTestCase(tc testCase) func(t *testing.T) {
	return func(t *testing.T) {
		got := tc.execute()
		want := tc.want()
		if diff := deep.Equal(got, want); diff != nil {
			t.Error(diff)
		}
	}
}

func TestSplit(t *testing.T) {
	tests := []testCase{
		{
			name: "split into one chunk",
			execute: func() interface{} {
				content := []byte("abcdefghijkl")
				chunks, _ := Split(content, 1)
				return chunks
			},
			want: func() interface{} {
				return []*chunk{
					{
						Bytes: []byte("abcdefghijkl"),
						Index: 0,
						Start: 0,
					},
				}
			},
		},
		{
			name: "split into two chunk",
			execute: func() interface{} {
				content := []byte("abcdefghijkl")
				chunks, _ := Split(content, 2)
				return chunks
			},
			want: func() interface{} {
				return []*chunk{
					{
						Bytes: []byte("abcdef"),
						Index: 0,
						Start: 0,
					},
					{
						Bytes: []byte("ghijkl"),
						Index: 1,
						Start: 6,
					},
				}
			},
		},
		{
			name: "split into four chunk",
			execute: func() interface{} {
				content := []byte("abcdefghijkl")
				chunks, _ := Split(content, 4)
				return chunks
			},
			want: func() interface{} {
				return []*chunk{
					{
						Bytes: []byte("abc"),
						Index: 0,
						Start: 0,
					},
					{
						Bytes: []byte("def"),
						Index: 1,
						Start: 3,
					},
					{
						Bytes: []byte("ghi"),
						Index: 2,
						Start: 6,
					},
					{
						Bytes: []byte("jkl"),
						Index: 3,
						Start: 9,
					},
				}
			},
		},
		{
			name: "split into five chunk",
			execute: func() interface{} {
				content := []byte("abcdefghijkl")
				chunks, _ := Split(content, 5)
				return chunks
			},
			want: func() interface{} {
				return []*chunk{
					{
						Bytes: []byte("abc"),
						Index: 0,
						Start: 0,
					},
					{
						Bytes: []byte("def"),
						Index: 1,
						Start: 3,
					},
					{
						Bytes: []byte("gh"),
						Index: 2,
						Start: 6,
					},
					{
						Bytes: []byte("ij"),
						Index: 3,
						Start: 8,
					},
					{
						Bytes: []byte("kl"),
						Index: 4,
						Start: 10,
					},
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, runTestCase(tc))
	}
}

func BenchmarkSplit(b *testing.B) {
	for input := 1; input <= 100000; input *= 10 {
		b.Run(fmt.Sprintf("input_%d", input), func(b *testing.B) {
			rand.Seed(time.Now().Unix())

			charset := "abcdefghijklmnopqrstuvwxyz1234567890"
			content := ""

			for i := 0; i < input; i++ {
				content += string(charset[rand.Int()%len(charset)])
			}

			payload := []byte(content)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				Split(payload, rand.Int())
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []testCase{
		{
			name: "merge one byte",
			execute: func() interface{} {
				chunks := []*chunk{
					{
						Index: 0,
						Start: 0,
						Bytes: []byte{1},
					},
				}
				merged := Merge(chunks, 1)
				return merged
			},
			want: func() interface{} {
				return []byte{1}
			},
		},
		{
			name: "merge multiple chunk",
			execute: func() interface{} {
				chunks := []*chunk{
					{
						Index: 0,
						Start: 3,
						Bytes: []byte{1},
					},
					{
						Index: 1,
						Start: 0,
						Bytes: []byte{2, 3, 4},
					},
				}
				merged := Merge(chunks, 4)
				return merged
			},
			want: func() interface{} {
				return []byte{2, 3, 4, 1}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, runTestCase(tc))
	}
}
