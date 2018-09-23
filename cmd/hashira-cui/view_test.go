package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	src := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	}
	in := "z"

	tcs := []struct {
		inSrc []string
		index int
		want  []string
	}{
		{
			inSrc: src,
			index: 0,
			want:  []string{in, "a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			inSrc: src,
			index: 1,
			want:  []string{"a", in, "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			inSrc: src,
			index: 100,
			want:  []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", in},
		},
		{
			inSrc: src,
			index: -1,
			want:  nil,
		},
	}

	for i, tc := range tcs {
		t.Logf("case %d testing", i)
		ret := insert(tc.inSrc, in, tc.index)
		require.Equal(t, ret, tc.want, "no.%d failed: [want] %v [got] %v", i, tc.want, ret)
	}
}

func TestRemove(t *testing.T) {
	src := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	}

	tcs := []struct {
		in   string
		want []string
	}{
		{
			in:   "a",
			want: []string{"b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			in:   "b",
			want: []string{"a", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			in:   "z",
			want: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			in:   "1",
			want: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
	}

	for i, tc := range tcs {
		t.Logf("case %d testing", i)
		ret := remove(src, tc.in)
		require.Equal(t, ret, tc.want, "no.%d failed: [want] %v [got] %v", i, tc.want, ret)
	}
}
