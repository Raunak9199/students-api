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

func (s *SQlite) GetStudents() ([]types.Student, error) {
	var students []types.Student
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			// slog.Error("query error: %w", err)
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
func (s *SQlite) DeleteStudent(id int64) (types.Student, error) {
	var student types.Student

	// Fetch the student details before deletion
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to prepare SELECT query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("failed to fetch student details: %w", err)
	}

	// Delete the student
	_, err = s.Db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to delete student with id %d: %w", id, err)
	}

	// Return the deleted student details
	return student, nil
}

func (s *SQlite) UpdateStudent(id int64, student types.Student) (types.Student, error) {

	var existingStudent types.Student
	err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).
		Scan(&existingStudent.Id, &existingStudent.Name, &existingStudent.Email, &existingStudent.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with id %d not found", id)
		}
		return types.Student{}, fmt.Errorf("failed to fetch student details: %w", err)
	}

	_, err = s.Db.Exec(
		"UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?",
		student.Name, student.Email, student.Age, id,
	)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to update student: %w", err)
	}

	return types.Student{
		Id:    id,
		Name:  student.Name,
		Email: student.Email,
		Age:   student.Age,
	}, nil
}
