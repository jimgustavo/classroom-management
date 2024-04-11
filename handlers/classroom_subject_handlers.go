// handlers/classroom_subjects_handlers.go

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

func AddSubjectToClassroom(w http.ResponseWriter, r *http.Request) {
	// Parse classroomID from URL parameter
	vars := mux.Vars(r)
	classroomID, err := strconv.Atoi(vars["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Parse subjectID from URL parameter
	subjectID, err := strconv.Atoi(vars["subjectID"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var requestBody models.AddSubjectRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	log.Printf("Adding subject %d to classroom %d", subjectID, classroomID) // Log before calling database function

	// Add subject to classroom
	if err := database.AddSubjectToClassroom(classroomID, subjectID, requestBody); err != nil {
		log.Println("Failed to add subject to classroom:", err) // Log the error
		http.Error(w, "Failed to add subject to classroom", http.StatusInternalServerError)
		return
	}

	log.Println("Subject added to classroom successfully")
	w.WriteHeader(http.StatusCreated)
}

func GetSubjectsInClassroom(w http.ResponseWriter, r *http.Request) {
	// Parse classroomID from URL parameter
	vars := mux.Vars(r)
	classroomID, err := strconv.Atoi(vars["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Get subjects in classroom
	subjects, err := database.GetSubjectsInClassroom(classroomID)
	if err != nil {
		http.Error(w, "Failed to get subjects in classroom", http.StatusInternalServerError)
		return
	}

	// Marshal subjects to JSON
	jsonSubjects, err := json.Marshal(subjects)
	if err != nil {
		http.Error(w, "Failed to marshal subjects to JSON", http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonSubjects)
}
