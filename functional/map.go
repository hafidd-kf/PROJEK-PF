package functional

func Map[T any, R any](
	data []T,
	ubah func(T) R,
) []R {
	hasil := make([]R, 0, len(data))

	for _, item := range data {
		hasil = append(hasil, ubah(item))
	}

	return hasil
}
