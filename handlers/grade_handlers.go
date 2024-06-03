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

// UploadGrades handles the uploading of grades for a classroom
func UploadGradesToClassroom(w http.ResponseWriter, r *http.Request) {
	var gradesData models.GradesData

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&gradesData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error parsing request payload: %v\n", err)
		return
	}

	log.Printf("Grades Data received: %+v\n", gradesData)

	// Extract the classroom ID from the request URL
	vars := mux.Vars(r)
	classroomID, err := strconv.Atoi(vars["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		log.Printf("Error parsing classroom ID: %v\n", err)
		return
	}

	// Iterate over the grades and insert them into the database
	for _, studentGrades := range gradesData.Grades {
		for _, grade := range studentGrades.Grades {
			err := database.InsertGradesInClassroom(studentGrades.StudentID, studentGrades.SubjectID, grade.Label, grade.Grade, classroomID)
			if err != nil {
				http.Error(w, "Error inserting grades", http.StatusInternalServerError)
				log.Printf("Error inserting grade: %v\n", err)
				return
			}
		}
	}

	// Prepare the JSON response
	response := models.Response{
		Message: "Grades uploaded successfully",
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		log.Printf("Error creating response: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)

	log.Println("Grades uploaded successfully")
}

// Handler function for fetching grades by classroom ID
func GetGradesByClassroomID(w http.ResponseWriter, r *http.Request) {
	// Extract the classroom ID from the URL path
	vars := mux.Vars(r)
	classroomID := vars["classroomID"]

	// Convert classroomID to an integer (assuming it's an integer)
	id, err := strconv.Atoi(classroomID)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Fetch grades for the specified classroom ID
	gradesData, err := database.FetchGradesByClassroomID(id)
	if err != nil {
		http.Error(w, "Error fetching grades", http.StatusInternalServerError)
		log.Printf("Error fetching grades: %v\n", err)
		return
	}

	// Convert gradesData to JSON
	responseData, err := json.Marshal(gradesData)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		log.Printf("Error creating response: %v\n", err)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

	log.Println("Grades retrieved successfully")
}
