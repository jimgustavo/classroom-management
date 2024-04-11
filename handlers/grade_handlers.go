// handlers/grade_handlers.go

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

// AddGrade handles the request to add a grade to a specific student in a specific subject and label.
func AddGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["studentID"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	subjectID, err := strconv.Atoi(vars["subjectID"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	var grade models.Grade
	// Decode the request body into the grade model
	err = json.NewDecoder(r.Body).Decode(&grade)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Assign student ID and subject ID to the grade
	grade.StudentID = studentID
	grade.SubjectID = subjectID

	// Call the database function to add the grade
	err = database.AddGrade(grade)
	if err != nil {
		http.Error(w, "Failed to add grade to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllStudentsWithGrades retrieves all students with their grades, labels, subjects, and classrooms.
func GetAllStudentsWithGrades(w http.ResponseWriter, r *http.Request) {
	// Log that the function has started
	log.Println("Retrieving all students with grades")

	// Call the database function to get all students with their grades, labels, subjects, and classrooms
	studentGradeInfo, err := database.GetAllStudentsWithGrades()
	if err != nil {
		// Log the error if retrieval fails
		log.Println("Failed to retrieve students with grades:", err)
		http.Error(w, "Failed to retrieve students with grades", http.StatusInternalServerError)
		return
	}

	// Log that retrieval was successful
	log.Println("Students with grades retrieved successfully")

	// Encode the response as JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(studentGradeInfo)
	if err != nil {
		// Log the encoding error
		log.Println("Failed to encode response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetGradeByStudentID handles the request to retrieve a JSON with the grade, label, subject, and classroom by providing the student ID.
func GetGradeByStudentID(w http.ResponseWriter, r *http.Request) {
	// Parse studentID from URL parameter
	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["studentID"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Call the database function to retrieve the grade information
	gradeInfo, err := database.GetGradeByStudentID(studentID)
	if err != nil {
		http.Error(w, "Failed to retrieve grade information", http.StatusInternalServerError)
		return
	}

	// Marshal the grade information into JSON
	gradeJSON, err := json.Marshal(gradeInfo)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(gradeJSON)
}
