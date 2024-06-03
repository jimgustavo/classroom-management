package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

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

// DeleteSubjectsByClassroomIDHandler removes all subjects associated with a classroom by classroom ID
func DeleteSubjectsByClassroomID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classroomID, err := strconv.Atoi(params["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteSubjectsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
