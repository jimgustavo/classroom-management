package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
)

// AddSubjectToStudentsHandler adds a subject to all students in a classroom
func AddSubjectToStudents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classroomID, err := strconv.Atoi(params["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	subjectID, err := strconv.Atoi(params["subjectID"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	err = database.AddSubjectToStudents(classroomID, subjectID)
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

// GetStudentsBySubjectIDHandler retrieves all students associated with a subject by subject ID
func GetStudentsBySubjectID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	subjectID, err := strconv.Atoi(params["subjectID"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	students, err := database.GetStudentsBySubjectID(subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(students)
}

// DeleteSubjectFromStudentsHandler removes a subject from all students in a classroom
func DeleteSubjectFromStudents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classroomID, err := strconv.Atoi(params["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	subjectID, err := strconv.Atoi(params["subjectID"])
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteSubjectFromStudents(classroomID, subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
