package functional

func HitungNilaiAkhir(tugas, uts, uas float64) float64 {
	return (tugas + uts + uas) / 3
}

func KonversiGrade(nilai float64) string {
	switch {
	case nilai >= 80:
		return "A"
	case nilai >= 70:
		return "B"
	case nilai >= 60:
		return "C"
	case nilai >= 50:
		return "D"
	default:
		return "E"
	}
}
