package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

// CreateStudent handles the creation of a new student
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateStudent(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// GetAllStudents retrieves all students
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := database.GetAllStudents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(students)
}

// GetStudent retrieves a specific student by their ID
func GetStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := database.GetStudentByID(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(student)
}

// GetAllStudentsWithClassroomAndSubjects retrieves all students along with their assigned classroom and subjects
func GetAllStudentsWithClassroomAndSubjects(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received: GetAllStudentsWithClassroomAndSubjects")

	// Get all students with their assigned classroom and subjects from the database
	students, err := database.GetAllStudentsWithClassroomAndSubjects()
	if err != nil {
		log.Printf("Error retrieving students with classroom and subjects: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Students retrieved successfully")

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode students into JSON and write response
	if err := json.NewEncoder(w).Encode(students); err != nil {
		log.Printf("Error encoding students into JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Response sent successfully")
}

// UpdateStudent updates the details of a specific student
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateStudent(id, &updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedStudent.ID = id
	json.NewEncoder(w).Encode(updatedStudent)
}

// DeleteStudent deletes a specific student
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteStudent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
