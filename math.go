package rrh

func MaxInt(values ...int) int {
	if len(values) == 0 {
		panic("MaxInt requires at least one arguments")
	}
	max := values[0]
	for _, v := range values {
		if max < v {
			max = v
		}
	}
	return max
}
