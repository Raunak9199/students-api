package sqlite

import (
	"database/sql"

	"githb.com/Raunak9199/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	age INTEGER NOT NULL
	)`)

	if err != nil {
		return nil, err
	}

	return &SQlite{
		Db: db,
	}, nil

}

func (s *SQlite) CreateStudent(name string, email string, age int) (int64, error) {

	stat, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stat.Close()

	result, err := stat.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil

}
