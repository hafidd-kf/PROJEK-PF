package repositories

import (
	"database/sql"
	"errors"

	"aplikasi-nilai-mahasiswa/models"
)

type MahasiswaRepository struct {
	DB *sql.DB
}

// READ: mengambil semua mahasiswa
func (r *MahasiswaRepository) GetAll() ([]models.Mahasiswa, error) {
	rows, err := r.DB.Query(`
		SELECT id, nim, nama, tugas, uts, uas, created_at
		FROM mahasiswa ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.Mahasiswa{}

	for rows.Next() {
		var m models.Mahasiswa

		err := rows.Scan(
			&m.ID, &m.NIM, &m.Nama,
			&m.Tugas, &m.UTS, &m.UAS,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		data = append(data, m)
	}

	return data, rows.Err()
}

// READ: mengambil mahasiswa berdasarkan ID
func (r *MahasiswaRepository) GetByID(id int) (models.Mahasiswa, error) {
	var m models.Mahasiswa

	err := r.DB.QueryRow(`
		SELECT id, nim, nama, tugas, uts, uas, created_at
		FROM mahasiswa WHERE id = $1
	`, id).Scan(
		&m.ID, &m.NIM, &m.Nama,
		&m.Tugas, &m.UTS, &m.UAS,
		&m.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return m, errors.New("mahasiswa tidak ditemukan")
	}

	return m, err
}

// CREATE: menambahkan mahasiswa
func (r *MahasiswaRepository) Create(m models.Mahasiswa) error {
	_, err := r.DB.Exec(`
		INSERT INTO mahasiswa (nim, nama, tugas, uts, uas)
		VALUES ($1, $2, $3, $4, $5)
	`, m.NIM, m.Nama, m.Tugas, m.UTS, m.UAS)

	return err
}

// UPDATE: mengubah mahasiswa
func (r *MahasiswaRepository) Update(m models.Mahasiswa) error {
	result, err := r.DB.Exec(`
		UPDATE mahasiswa
		SET nim = $1, nama = $2, tugas = $3, uts = $4, uas = $5
		WHERE id = $6
	`, m.NIM, m.Nama, m.Tugas, m.UTS, m.UAS, m.ID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("mahasiswa tidak ditemukan")
	}

	return nil
}

// DELETE: menghapus mahasiswa
func (r *MahasiswaRepository) Delete(id int) error {
	result, err := r.DB.Exec(
		"DELETE FROM mahasiswa WHERE id = $1", id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("mahasiswa tidak ditemukan")
	}

	return nil
}
