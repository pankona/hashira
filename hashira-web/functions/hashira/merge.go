package hashira

import "sort"

func mergeStringSlice(a, b []string) []string {
	ar := stringSliceToRankedDataSlice(a)
	br := stringSliceToRankedDataSlice(b)

	ret := mergeRankedData(ar, br)

	return rankedDataSliceToStringSlice(ret)
}

func stringSliceToRankedDataSlice(a []string) []rankedData {
	ret := make([]rankedData, 0, len(a))
	for i := range a {
		ret = append(ret, rankedData{a[i], i})
	}
	return ret
}

func rankedDataSliceToStringSlice(a []rankedData) []string {
	ret := make([]string, 0, len(a))
	for _, v := range a {
		ret = append(ret, v.data.(string))
	}
	return ret
}

type rankedData struct {
	data interface{}
	rank int
}

func mergeRankedData(ar, br []rankedData) []rankedData {
	br = removeDuplicates(br, ar)

	ret := append(ar, br...)

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].rank < ret[j].rank
	})

	return ret
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// return a - b
func removeDuplicates(a, b []rankedData) []rankedData {
	ret := make([]rankedData, 0, max(len(a), len(b)))
	for i, v := range a {
		if includes(b, v) {
			continue
		}
		ret = append(ret, a[i])
	}
	return ret
}

func includes(a []rankedData, i rankedData) bool {
	for _, v := range a {
		if v.data == i.data {
			return true
		}
	}
	return false
}
