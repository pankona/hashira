package hashira

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMergeRankedData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inA  []int
		inB  []int
		want []int
	}{
		{
			inA:  []int{1, 3, 5},
			inB:  []int{2, 4, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			inA:  []int{1, 3, 5, 7, 9, 11, 13},
			inB:  []int{2, 4, 6, 8, 10},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13},
		},
		{
			inA:  []int{1, 3, 5, 7, 9},
			inB:  []int{2, 4, 6, 8, 10, 11, 13},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13},
		},
		{
			inA:  []int{1, 3, 5, 7, 9},
			inB:  nil,
			want: []int{1, 3, 5, 7, 9},
		},
		{
			inA:  nil,
			inB:  []int{1, 3, 5, 7, 9},
			want: []int{1, 3, 5, 7, 9},
		},
		{
			inA:  []int{1, 3, 5, 7, 9},
			inB:  []int{1, 3, 5, 7, 9, 11},
			want: []int{1, 3, 5, 7, 9, 11},
		},
		{
			inA:  []int{1, 3, 5, 7, 9},
			inB:  []int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10, 11},
			want: []int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10, 11},
		},
		{
			inA:  []int{1, 3, 5, 7, 9},
			inB:  []int{9, 7, 5, 3, 1},
			want: []int{1, 3, 5, 7, 9},
		},
		{
			inA:  []int{1, 3, 5, 7, 9, 2},
			inB:  []int{2, 4, 6, 8, 10},
			want: []int{1, 3, 4, 5, 6, 7, 8, 9, 10, 2},
		},
	}

	for i, tt := range tests {
		tt := tt
		if i != 8 {
			continue
		}
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			a := intSliceToRankedDataSlice(tt.inA)
			b := intSliceToRankedDataSlice(tt.inB)

			gotRankedDataSlice := mergeRankedData(a, b)
			got := rankedDataSliceToIntSlice(gotRankedDataSlice)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected result: diff: %s", diff)
			}
		})
	}
}

func intSliceToRankedDataSlice(a []int) []rankedData {
	ret := make([]rankedData, 0, len(a))
	for i := range a {
		ret = append(ret, rankedData{a[i], i})
	}
	return ret
}

func rankedDataSliceToIntSlice(a []rankedData) []int {
	ret := make([]int, 0, len(a))
	for _, v := range a {
		ret = append(ret, v.data.(int))
	}
	return ret
}
