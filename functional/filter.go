package functional

func Filter[T any](
	data []T,
	kondisi func(T) bool,
) []T {
	hasil := make([]T, 0)

	for _, item := range data {
		if kondisi(item) {
			hasil = append(hasil, item)
		}
	}

	return hasil
}
