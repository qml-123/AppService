package utils

func SliceToMap(s []int64) map[int64]bool {
	m := make(map[int64]bool)
	for _, a := range s {
		m[a] = true
	}
	return m
}

func MapToSlice(m map[int64]bool) []int64 {
	s := make([]int64, 0)
	for a := range m {
		s = append(s, a)
	}
	return s
}

func StrSliceToMap(s []string) map[string]bool {
	m := make(map[string]bool)
	for _, a := range s {
		m[a] = true
	}
	return m
}

func StrMapToSlice(m map[string]bool) []string {
	s := make([]string, 0)
	for a := range m {
		s = append(s, a)
	}
	return s
}
