package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"githb.com/Raunak9199/students-api/internal/storage"
	"githb.com/Raunak9199/students-api/internal/types"
	"githb.com/Raunak9199/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			slog.Info("No student data provided")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body can't be empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError((err)))
			return
		}

		// Request Validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("Student created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			slog.Error("Failed to create student", slog.String("Er", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// w.Write([]byte("Welcome to students api"))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("GetById", slog.String("id", id))

		intId, er := strconv.ParseInt(id, 10, 64)
		if er != nil {
			slog.Error("Failed to parse id", slog.String("Er", er.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(er))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("Failed to get student", slog.String("Er", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetStudents()

		if err != nil {
			slog.Error("Failed to get students", slog.String("Er", err.Error()))

			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := r.PathValue("id")
		if idstr == "" {
			slog.Error("Missing student ID")
			response.WriteJson(w, http.StatusBadRequest, "student ID is required")
			return
		}

		id, err := strconv.ParseInt(idstr, 10, 64)

		if err != nil {
			slog.Error("Invalid student ID", slog.String("id", idstr), slog.String("Er", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, "invalid student ID")
			return
		}

		slog.Info("Deleting student", slog.Int64("id", id))

		deletedStudent, err := storage.DeleteStudent(id)
		if err != nil {
			slog.Error("Failed to delete student", slog.String("Er", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with the details of the deleted student
		response.WriteJson(w, http.StatusOK, deletedStudent)
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			response.WriteJson(w, http.StatusBadRequest, "student ID is required")
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, "invalid student ID")
			return
		}

		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}

		updatedStudent, err := storage.UpdateStudent(id, student)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJson(w, http.StatusOK, updatedStudent)
	}
}
