package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
	"github.com/xuri/excelize/v2"
)

// CreateStudent handles the creation of a new student
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateStudent(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// GetAllStudents retrieves all students
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := database.GetAllStudents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(students)
}

// GetStudent retrieves a specific student by their ID
func GetStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := database.GetStudentByID(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(student)
}

// UpdateStudent updates the details of a specific student
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateStudent(id, &updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedStudent.ID = id
	json.NewEncoder(w).Encode(updatedStudent)
}

// DeleteStudent deletes a specific student
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteStudent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllStudentsWithClassroomAndSubjects retrieves all students along with their assigned classroom and subjects
func GetAllStudentsWithClassroomAndSubjects(w http.ResponseWriter, r *http.Request) {
	// Get all students with their assigned classroom and subjects from the database
	students, err := database.GetAllStudentsWithClassroomAndSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Students retrieved successfully")

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode students into JSON and write response
	if err := json.NewEncoder(w).Encode(students); err != nil {
		log.Printf("Error encoding students into JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Response sent successfully")
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

// UploadStudentsFromExcel handles uploading students from an Excel file
func UploadStudentsFromExcel(w http.ResponseWriter, r *http.Request) {
	classroomID := mux.Vars(r)["classroomID"]
	startCell := r.URL.Query().Get("startCell")
	endCell := r.URL.Query().Get("endCell")
	sheetName := r.URL.Query().Get("sheetName")

	// Parse classroomID to int
	classroomIDInt, err := strconv.Atoi(classroomID)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Parse the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Open the Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Error opening Excel file", http.StatusBadRequest)
		return
	}

	// Extract students from the specified range
	students, err := extractStudentsFromExcel(f, sheetName, startCell, endCell)
	if err != nil {
		http.Error(w, "Error extracting students from Excel file", http.StatusBadRequest)
		return
	}

	// Insert students into the database
	for _, student := range students {
		student.ClassroomID = classroomIDInt
		err := database.InsertStudent(&student)
		if err != nil {
			http.Error(w, "Error inserting student into database", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Students uploaded successfully")
}

func extractStudentsFromExcel(f *excelize.File, sheetName, startCell, endCell string) ([]models.Student, error) {
	var students []models.Student

	startCol, startRow, err := excelize.CellNameToCoordinates(startCell)
	if err != nil {
		return nil, err
	}
	endCol, endRow, err := excelize.CellNameToCoordinates(endCell)
	if err != nil {
		return nil, err
	}

	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			cell, err := excelize.CoordinatesToCellName(col, row)
			if err != nil {
				return nil, err
			}
			name, err := f.GetCellValue(sheetName, cell)
			if err != nil {
				return nil, err
			}
			if name != "" {
				students = append(students, models.Student{Name: name})
			}
		}
	}

	return students, nil
}
