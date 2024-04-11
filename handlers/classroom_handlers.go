package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

// CreateClassroom handles the creation of a new classroom
func CreateClassroom(w http.ResponseWriter, r *http.Request) {
	var classroom models.Classroom
	err := json.NewDecoder(r.Body).Decode(&classroom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert classroom into the database
	err = database.CreateClassroom(&classroom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(classroom)
}

// GetAllClassrooms retrieves all classrooms
func GetAllClassrooms(w http.ResponseWriter, r *http.Request) {
	classrooms, err := database.GetAllClassrooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(classrooms)
}

// GetClassroom retrieves a specific classroom by its ID
func GetClassroom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	classroom, err := database.GetClassroomByID(id)
	if err != nil {
		http.Error(w, "Classroom not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(classroom)
}

/*
// GetSubjectsByClassroomID retrieves subjects assigned to a classroom
func GetSubjectsByClassroomID(w http.ResponseWriter, r *http.Request) {
	// Get the classroom ID from the request parameters
	params := mux.Vars(r)
	classroomID := params["classroomID"]

	// Get subjects assigned to the specified classroom from the database
	subjects, err := database.GetSubjectsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode subjects into JSON and write response
	if err := json.NewEncoder(w).Encode(subjects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
*/
// GetSubjectsByClassroomID retrieves subjects assigned to a classroom along with students assigned to each subject by classroom ID
func GetSubjectsByClassroomID(w http.ResponseWriter, r *http.Request) {
	// Get the classroom ID from the request parameters
	params := mux.Vars(r)
	classroomID := params["classroomID"]

	// Get subjects assigned to the specified classroom along with students assigned to each subject
	subjectsWithStudents, err := database.GetSubjectsAndStudentsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode subjectsWithStudents into JSON and write response
	if err := json.NewEncoder(w).Encode(subjectsWithStudents); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler function to get students by classroom ID
func GetStudentsByClassroom(w http.ResponseWriter, r *http.Request) {
	// Get the classroom ID from the request URL
	params := mux.Vars(r)
	classroomID, err := strconv.Atoi(params["classroom_id"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Retrieve the list of students associated with the specified classroom ID from the database
	students, err := database.GetStudentsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of students in the response
	json.NewEncoder(w).Encode(students)
}

// UpdateClassroom updates the details of a specific classroom
func UpdateClassroom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	var updatedClassroom models.Classroom
	err = json.NewDecoder(r.Body).Decode(&updatedClassroom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateClassroom(id, &updatedClassroom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedClassroom.ID = id
	json.NewEncoder(w).Encode(updatedClassroom)
}

// DeleteClassroom deletes a specific classroom
func DeleteClassroom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteClassroom(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
