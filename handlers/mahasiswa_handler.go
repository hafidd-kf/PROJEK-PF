package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"aplikasi-nilai-mahasiswa/functional"
	"aplikasi-nilai-mahasiswa/models"
	"aplikasi-nilai-mahasiswa/repositories"
)

type MahasiswaView struct {
	models.Mahasiswa
	NilaiAkhir float64
	Grade      string
	Status     string
}

type Statistik struct {
	Total      int
	Lulus      int
	TidakLulus int
	RataRata   float64
	Tertinggi  float64
	Terendah   float64
	NamaTinggi string
	NamaRendah string
}

type Handler struct {
	Repo *repositories.MahasiswaRepository
	Tmpl *template.Template
}

func NewHandler(repo *repositories.MahasiswaRepository) (*Handler, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}).ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		Repo: repo,
		Tmpl: tmpl,
	}, nil
}

func (h *Handler) render(w http.ResponseWriter, page string, data any) {
	err := h.Tmpl.ExecuteTemplate(w, page, data)
	if err != nil {
		http.Error(w, "Gagal menampilkan halaman", http.StatusInternalServerError)
	}
}

func nilaiAkhir(m models.Mahasiswa) float64 {
	return functional.HitungNilaiAkhir(m.Tugas, m.UTS, m.UAS)
}

func buatView(data []models.Mahasiswa) []MahasiswaView {
	return functional.Map(data, func(m models.Mahasiswa) MahasiswaView {
		nilai := nilaiAkhir(m)
		status := "TIDAK LULUS"

		if nilai >= 70 {
			status = "LULUS"
		}

		return MahasiswaView{
			Mahasiswa:  m,
			NilaiAkhir: nilai,
			Grade:      functional.KonversiGrade(nilai),
			Status:     status,
		}
	})
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data", 500)
		return
	}

	view := buatView(data)

	var rataRata float64
	if len(view) > 0 {
		total := functional.Reduce(view, float64(0),
			func(total float64, m MahasiswaView) float64 {
				return total + m.NilaiAkhir
			})
		rataRata = total / float64(len(view))
	}

	stat := Statistik{
		Total:    len(view),
		RataRata: rataRata,
	}

	for _, m := range view {
		if m.NilaiAkhir >= 70 {
			stat.Lulus++
		} else {
			stat.TidakLulus++
		}
	}

	h.render(w, "dashboard.html", map[string]any{
		"Statistik": stat,
	})
}

func (h *Handler) DaftarMahasiswa(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data", 500)
		return
	}

	h.render(w, "mahasiswa.html", map[string]any{
		"Mahasiswa": buatView(data),
	})
}

func (h *Handler) TambahMahasiswa(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.render(w, "tambah.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", 405)
		return
	}

	m, err := bacaForm(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := h.Repo.Create(m); err != nil {
		http.Error(w, "Gagal menyimpan data. Periksa NIM apakah sudah digunakan.", 400)
		return
	}

	http.Redirect(w, r, "/mahasiswa", http.StatusSeeOther)
}

func bacaForm(r *http.Request) (models.Mahasiswa, error) {
	var m models.Mahasiswa

	m.NIM = strings.TrimSpace(r.FormValue("nim"))
	m.Nama = strings.TrimSpace(r.FormValue("nama"))

	if m.NIM == "" || m.Nama == "" {
		return m, fmt.Errorf("NIM dan nama wajib diisi")
	}

	var err error

	m.Tugas, err = strconv.ParseFloat(r.FormValue("tugas"), 64)
	if err != nil {
		return m, fmt.Errorf("nilai tugas harus berupa angka")
	}

	m.UTS, err = strconv.ParseFloat(r.FormValue("uts"), 64)
	if err != nil {
		return m, fmt.Errorf("nilai UTS harus berupa angka")
	}

	m.UAS, err = strconv.ParseFloat(r.FormValue("uas"), 64)
	if err != nil {
		return m, fmt.Errorf("nilai UAS harus berupa angka")
	}

	if m.Tugas < 0 || m.Tugas > 100 ||
		m.UTS < 0 || m.UTS > 100 ||
		m.UAS < 0 || m.UAS > 100 {
		return m, fmt.Errorf("nilai harus berada antara 0 sampai 100")
	}

	return m, nil
}

func (h *Handler) EditMahasiswa(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", 400)
		return
	}

	if r.Method == http.MethodGet {
		m, err := h.Repo.GetByID(id)
		if err != nil {
			http.Error(w, "Mahasiswa tidak ditemukan", 404)
			return
		}

		h.render(w, "edit.html", m)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", 405)
		return
	}

	m, err := bacaForm(r)
	if err != nil {
		http.Error(w, "Input tidak valid", 400)
		return
	}

	m.ID = id

	if err := h.Repo.Update(m); err != nil {
		http.Error(w, "Gagal mengubah data", 400)
		return
	}

	http.Redirect(w, r, "/mahasiswa", 303)
}

func (h *Handler) HapusMahasiswa(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", 405)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", 400)
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, "Gagal menghapus data", 400)
		return
	}

	http.Redirect(w, r, "/mahasiswa", 303)
}

func (h *Handler) FilterLulus(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data", 500)
		return
	}

	view := buatView(data)

	lulus := functional.Filter(view, func(m MahasiswaView) bool {
		return m.NilaiAkhir >= 70
	})

	h.render(w, "lulus.html", map[string]any{
		"Mahasiswa": lulus,
	})
}

func (h *Handler) GradeMahasiswa(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data", 500)
		return
	}

	h.render(w, "grade.html", map[string]any{
		"Mahasiswa": buatView(data),
	})
}

func (h *Handler) StatistikKelas(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data", 500)
		return
	}

	view := buatView(data)
	stat := Statistik{Total: len(view)}

	if len(view) > 0 {
		stat.Tertinggi = view[0].NilaiAkhir
		stat.Terendah = view[0].NilaiAkhir
		stat.NamaTinggi = view[0].Nama
		stat.NamaRendah = view[0].Nama

		total := functional.Reduce(view, float64(0),
			func(total float64, m MahasiswaView) float64 {
				return total + m.NilaiAkhir
			})

		stat.RataRata = total / float64(len(view))

		for _, m := range view {
			if m.NilaiAkhir > stat.Tertinggi {
				stat.Tertinggi = m.NilaiAkhir
				stat.NamaTinggi = m.Nama
			}

			if m.NilaiAkhir < stat.Terendah {
				stat.Terendah = m.NilaiAkhir
				stat.NamaRendah = m.Nama
			}

			if m.NilaiAkhir >= 70 {
				stat.Lulus++
			} else {
				stat.TidakLulus++
			}
		}
	}

	h.render(w, "statistik.html", stat)
}
