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

func GetSubjectsByTeacherID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		log.Println("Invalid teacher ID:", err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	log.Println("Fetching subjects for teacher ID:", teacherID)
	subjects, err := database.GetSubjectsByTeacherID(teacherID)
	if err != nil {
		log.Println("Error fetching subjects:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Fetched subjects:", subjects)
	if err := json.NewEncoder(w).Encode(subjects); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateSubject handles the creation of a new subject
func CreateSubject(w http.ResponseWriter, r *http.Request) {
	var subject models.Subject
	err := json.NewDecoder(r.Body).Decode(&subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if subject.Name == "" {
		http.Error(w, "Subject name cannot be empty", http.StatusBadRequest)
		return
	}

	err = database.CreateSubject(&subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subject)
}

// GetAllSubjects retrieves all subjects
func GetAllSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := database.GetAllSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(subjects)
}

// GetSubject retrieves a specific subject by its ID
func GetSubject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	subject, err := database.GetSubjectByID(id)
	if err != nil {
		http.Error(w, "Subject not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(subject)
}

// UpdateSubject updates the details of a specific subject
func UpdateSubject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	var updatedSubject models.Subject
	err = json.NewDecoder(r.Body).Decode(&updatedSubject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedSubject.Name == "" {
		http.Error(w, "Subject name cannot be empty", http.StatusBadRequest)
		return
	}

	err = database.UpdateSubject(id, &updatedSubject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedSubject.ID = id
	json.NewEncoder(w).Encode(updatedSubject)
}

// DeleteSubject deletes a specific subject
func DeleteSubject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteSubject(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSubjectsByStudentIDHandler retrieves all subjects associated with a student by student ID
func GetSubjectsByStudentID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentID, err := strconv.Atoi(params["studentID"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	subjects, err := database.GetSubjectsByStudentID(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(subjects)
}

// RemoveGradeLabelFromSubjectByTerm removes a grade label from a subject for a specific term
func RemoveGradeLabelFromSubjectByTerm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectIDStr := vars["subjectID"]
	gradeLabelIDStr := vars["gradeLabelID"]
	termIDStr := vars["termID"]

	// Convert subjectIDStr, gradeLabelIDStr and termIDStr to int
	subjectID, err := strconv.Atoi(subjectIDStr)
	if err != nil {
		http.Error(w, "Invalid grade label ID", http.StatusBadRequest)
		return
	}

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

	if err := database.RemoveGradeLabelFromSubjectByTerm(subjectID, gradeLabelID, termID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
