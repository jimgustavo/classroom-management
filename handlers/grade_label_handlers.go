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

func GetGradeLabelsByTeacherID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		log.Println("Invalid teacher ID:", err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	log.Println("Fetching grade labels for teacher ID:", teacherID)
	gradeLabels, err := database.GetGradeLabelsByTeacherID(teacherID)
	if err != nil {
		log.Println("Error fetching grade labels:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Fetched grade labels:", gradeLabels)
	if err := json.NewEncoder(w).Encode(gradeLabels); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateGradeLabel handles the creation of a new grade label
func CreateGradeLabel(w http.ResponseWriter, r *http.Request) {
	var gradeLabel models.GradeLabel
	err := json.NewDecoder(r.Body).Decode(&gradeLabel)
	if err != nil {
		http.Error(w, createJSONError("Invalid request payload", err.Error()), http.StatusBadRequest)
		return
	}

	err = database.CreateGradeLabel(&gradeLabel)
	if err != nil {
		http.Error(w, createJSONError("Failed to create grade label", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gradeLabel)
}

// createJSONError creates a JSON error response
func createJSONError(message, details string) string {
	errorResponse := map[string]string{
		"message": message,
		"details": details,
	}
	errorJSON, _ := json.Marshal(errorResponse)
	return string(errorJSON)
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
	log.Printf("Deleting grade label with id: %s", gradeLabelID)

	err := database.DeleteGradeLabel(gradeLabelID)
	if err != nil {
		log.Printf("Error deleting grade label with id %s: %v", gradeLabelID, err)
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

///////////////////////////////ACADEMIC REINFORCEMENT/////////////////////////////

func CreateReinforcementGradeLabel(w http.ResponseWriter, r *http.Request) {
	var gradeLabel models.ReinforcementGradeLabel
	if err := json.NewDecoder(r.Body).Decode(&gradeLabel); err != nil {
		log.Println("Error decoding request body:", err)
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.AddReinforcementGradeLabel(gradeLabel); err != nil {
		log.Println("Error adding reinforcement grade label to database:", err)
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetAllReinforcementGradeLabels(w http.ResponseWriter, r *http.Request) {
	gradeLabels, err := database.GetAllReinforcementGradeLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gradeLabels)
}

func GetReinforcementGradeLabelsByTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gradeLabels, err := database.GetReinforcementGradeLabelsByTeacher(teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gradeLabels)
}

func GetReinforcementGradeLabelsByClassroomAndTerm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classroomID, err := strconv.Atoi(vars["classroomID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	termID, err := strconv.Atoi(vars["termID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gradeLabels, err := database.GetReinforcementGradeLabelsByClassroomAndTerm(classroomID, termID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if gradeLabels == nil {
		gradeLabels = []models.ReinforcementGradeLabel{} // Return an empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gradeLabels)
}

func DeleteReinforcementGradeLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DeleteReinforcementGradeLabel(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
