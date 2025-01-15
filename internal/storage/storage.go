package storage

import "githb.com/Raunak9199/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	DeleteStudent(id int64) (types.Student, error)
	UpdateStudent(id int64, student types.Student) (types.Student, error)
}
