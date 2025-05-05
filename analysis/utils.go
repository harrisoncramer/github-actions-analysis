package analysis

func percentile(sorted []int, p float64) int {
	if len(sorted) == 0 {
		return 0
	}
	k := int(float64(len(sorted)-1) * p)
	return sorted[k]
}
