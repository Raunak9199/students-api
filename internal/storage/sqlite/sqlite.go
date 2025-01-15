package sqlite

import (
	"database/sql"
	"fmt"

	"githb.com/Raunak9199/students-api/internal/config"
	"githb.com/Raunak9199/students-api/internal/types"
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

func (s *SQlite) GetStudentById(id int64) (types.Student, error) {
	var student types.Student
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}
