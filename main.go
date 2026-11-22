
package main

import (
	"log"
	"net/http"

	"aplikasi-nilai-mahasiswa/config"
	"aplikasi-nilai-mahasiswa/handlers"
	"aplikasi-nilai-mahasiswa/repositories"
)

func main() {
	// Koneksi PostgreSQL
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Repository
	repo := &repositories.MahasiswaRepository{
		DB: db,
	}

	// Handler
	h, err := handlers.NewHandler(repo)
	if err != nil {
		log.Fatal("Gagal membaca template:", err)
	}

	// Routing
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		h.Dashboard(w, r)
	})

	mux.HandleFunc("/mahasiswa", h.DaftarMahasiswa)
	mux.HandleFunc("/mahasiswa/tambah", h.TambahMahasiswa)
	mux.HandleFunc("/mahasiswa/edit", h.EditMahasiswa)
	mux.HandleFunc("/mahasiswa/hapus", h.HapusMahasiswa)
	mux.HandleFunc("/lulus", h.FilterLulus)
	mux.HandleFunc("/grade", h.GradeMahasiswa)
	mux.HandleFunc("/statistik", h.StatistikKelas)

	// File CSS
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	))

	// Port server
	port := "8081"

	log.Println("Server berjalan di http://localhost:" + port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
