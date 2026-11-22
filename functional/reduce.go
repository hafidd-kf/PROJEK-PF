package functional

func Reduce[T any, R any](
	data []T,
	awal R,
	operasi func(R, T) R,
) R {
	hasil := awal

	for _, item := range data {
		hasil = operasi(hasil, item)
	}

	return hasil
}
