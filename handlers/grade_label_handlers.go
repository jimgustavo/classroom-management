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

// CreateGradeLabel handles the creation of a new grade label
func CreateGradeLabel(w http.ResponseWriter, r *http.Request) {
	var gradeLabel models.GradeLabel
	err := json.NewDecoder(r.Body).Decode(&gradeLabel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateGradeLabel(&gradeLabel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gradeLabel)
}

// GetAllGradeLabels retrieves all grade labels
func GetAllGradeLabels(w http.ResponseWriter, r *http.Request) {
	gradeLabels, err := database.GetAllGradeLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gradeLabels)
}

// GetGradeLabel retrieves a specific grade label by its ID
func GetGradeLabel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gradeLabelID := params["id"]

	gradeLabel, err := database.GetGradeLabelByID(gradeLabelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gradeLabel)
}

// UpdateGradeLabel updates an existing grade label
func UpdateGradeLabel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gradeLabelIDStr := params["id"]

	// Convert gradeLabelIDStr to int
	gradeLabelID, err := strconv.Atoi(gradeLabelIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedGradeLabel models.GradeLabel
	err = json.NewDecoder(r.Body).Decode(&updatedGradeLabel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign the converted ID to updatedGradeLabel
	updatedGradeLabel.ID = gradeLabelID
	err = database.UpdateGradeLabel(&updatedGradeLabel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGradeLabel)
}

// DeleteGradeLabel deletes a grade label by its ID
func DeleteGradeLabel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gradeLabelID := params["id"]

	err := database.DeleteGradeLabel(gradeLabelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AssignGradeLabelToSubject assigns a grade label to a subject for a specific term
func AssignGradeLabelToSubjectByTerm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectID := vars["subjectID"]
	gradeLabelIDStr := vars["gradeLabelID"] // Extract grade label ID from URL path
	termIDStr := vars["termID"]             // Extract term ID from URL path

	// Convert gradeLabelIDStr and termIDStr to int
	gradeLabelID, err := strconv.Atoi(gradeLabelIDStr)
	if err != nil {
		http.Error(w, "Invalid grade label ID", http.StatusBadRequest)
		return
	}
	termID, err := strconv.Atoi(termIDStr)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	if err := database.AssignGradeLabelToSubjectByTerm(subjectID, gradeLabelID, termID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetGradeLabelsForSubject retrieves all grade labels assigned to a subject for a specific term
func GetGradeLabelsForSubject(w http.ResponseWriter, r *http.Request) {
	log.Printf("hit the endpoint: GetGradeLabelsForSubject")
	vars := mux.Vars(r)
	subjectIDStr := vars["subjectID"]
	termIDStr := vars["termID"] // Extract term ID from URL path

	// Convert termIDStr to int
	subjectID, err := strconv.Atoi(subjectIDStr)
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	// Convert termIDStr to int
	termID, err := strconv.Atoi(termIDStr)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	gradeLabels, err := database.GetGradeLabelsForSubject(subjectID, termID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gradeLabels)
}
