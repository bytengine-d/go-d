package lang

func Index[T any](arr []T, item T, comparison func(v1, v2 T) bool) int64 {
	var index int64 = -1

	for i, v := range arr {
		if comparison(item, v) {
			index = int64(i)
			break
		}
	}

	return index
}
