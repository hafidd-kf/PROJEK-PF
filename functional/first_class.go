package functional

type OperasiNilai func(float64) string

func GunakanOperasi(
	nilai float64,
	operasi OperasiNilai,
) string {
	return operasi(nilai)
}
