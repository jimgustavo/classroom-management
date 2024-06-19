// handlers/classroom_handlers.go

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

func GetClassroomsByTeacherID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	classrooms, err := database.GetClassroomsByTeacherID(teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(classrooms)
}

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

// AddSubjectToClassroom adds a subject to a classroom
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

	log.Printf("Adding subject %d to classroom %d", subjectID, classroomID)

	// Add subject to classroom
	if err := database.AddSubjectToClassroom(classroomID, subjectID); err != nil {
		log.Println("Failed to add subject to classroom:", err)
		http.Error(w, "Failed to add subject to classroom", http.StatusInternalServerError)
		return
	}

	log.Println("Subject added to classroom successfully")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Subject added to classroom successfully"})
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

func GetSubjectsInClassroom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classroomID, err := strconv.Atoi(vars["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	subjects, err := database.GetSubjectsInClassroom(classroomID)
	if err != nil {
		http.Error(w, "Failed to get subjects in classroom", http.StatusInternalServerError)
		return
	}

	jsonSubjects, err := json.Marshal(subjects)
	if err != nil {
		http.Error(w, "Failed to marshal subjects to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonSubjects)
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

// UnrollStudentFromClassroom removes a student from a classroom
func UnrollStudentFromClassroom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classroomID := vars["classroomID"]
	studentID := vars["studentID"]

	log.Printf("Attempting to unroll student %s from classroom %s", studentID, classroomID)

	err := database.UnrollStudentFromClassroom(classroomID, studentID)
	if err != nil {
		log.Printf("Error unrolling student from classroom: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully unrolled student %s from classroom %s", studentID, classroomID)
	w.WriteHeader(http.StatusNoContent)
}

// RemoveSubjectFromClassroom removes a subject from a classroom
func RemoveSubjectFromClassroom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classroomID := vars["classroomID"]
	subjectID := vars["subjectID"]

	err := database.RemoveSubjectFromClassroom(classroomID, subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
