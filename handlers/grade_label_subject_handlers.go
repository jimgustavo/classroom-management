// handlers/grade_label_subject_handlers.go

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
)

// AssignGradeLabelToSubject assigns a grade label to a subject
func AssignGradeLabelToSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectID := vars["subjectID"]
	gradeLabelIDStr := vars["gradeLabelID"] // Extract grade label ID from URL path

	// Convert gradeLabelIDStr to int
	gradeLabelID, err := strconv.Atoi(gradeLabelIDStr)
	if err != nil {
		http.Error(w, "Invalid grade label ID", http.StatusBadRequest)
		return
	}

	if err := database.AssignGradeLabelToSubject(subjectID, gradeLabelID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetGradeLabelsForSubject retrieves all grade labels assigned to a subject
func GetGradeLabelsForSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectID := vars["subjectID"]

	gradeLabels, err := database.GetGradeLabelsForSubject(subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gradeLabels)
}
