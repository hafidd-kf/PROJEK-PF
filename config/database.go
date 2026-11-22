package config

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=hafidd dbname=db_nilai_mahasiswa sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("koneksi database gagal: %w", err)
	}

	fmt.Println("Berhasil terhubung ke PostgreSQL!")

	return db, nil
}
