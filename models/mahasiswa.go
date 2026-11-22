package models

import "time"

type Mahasiswa struct {
	ID        int
	NIM       string
	Nama      string
	Tugas     float64
	UTS       float64
	UAS       float64
	CreatedAt time.Time
}
