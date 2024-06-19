package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"

	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

// Handler function to generate the teacher_grades.xlsx report
func GenerateTeacherGradesReport(w http.ResponseWriter, r *http.Request) {
	// Extract the classroom ID and term ID from the URL path
	vars := mux.Vars(r)
	classroomID := vars["classroomID"]
	termID := vars["termID"]

	// Convert classroomID to an integer (assuming it's an integer)
	classroomIDInt, err := strconv.Atoi(classroomID)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Convert termID to an integer (assuming it's an integer)
	termIDInt, err := strconv.Atoi(termID)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	// Fetch grades for the specified classroom ID and term ID
	gradesData, err := database.FetchGradesByClassroomIDAndTermID(classroomIDInt, termIDInt)
	if err != nil {
		http.Error(w, "Error fetching grades", http.StatusInternalServerError)
		log.Printf("Error fetching grades: %v\n", err)
		return
	}

	// Log fetched grades data
	gradesJSON, _ := json.MarshalIndent(gradesData, "", "  ")
	log.Printf("Fetched Grades Data: %s\n", gradesJSON)

	// Fetch students and subjects for the classroom
	students, err := database.GetStudentsByClassroomID(classroomIDInt)
	if err != nil {
		http.Error(w, "Error fetching students", http.StatusInternalServerError)
		log.Printf("Error fetching students: %v\n", err)
		return
	}

	subjects, err := database.GetSubjectsInClassroom(classroomIDInt)
	if err != nil {
		http.Error(w, "Error fetching subjects", http.StatusInternalServerError)
		log.Printf("Error fetching subjects: %v\n", err)
		return
	}

	// Generate the Excel file
	file := excelize.NewFile()
	sheetName := "Teacher Grades"
	file.NewSheet(sheetName)

	// Set the header row
	headers := []string{"Number", "Student Name"}
	for _, subject := range subjects {
		headers = append(headers, subject.Name)
	}
	headers = append(headers, "Average")

	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		file.SetCellValue(sheetName, cell, header)
	}

	// Fill the student grades
	for i, student := range students {
		rowNum := i + 2
		file.SetCellValue(sheetName, "A"+strconv.Itoa(rowNum), i+1)
		file.SetCellValue(sheetName, "B"+strconv.Itoa(rowNum), student.Name)

		totalGrades := 0.0
		gradeCount := 0

		for j, subject := range subjects {
			cell := string(rune('C'+j)) + strconv.Itoa(rowNum)
			grade := getGradeForStudent(gradesData, student.ID, subject.ID, termIDInt)
			if grade != nil {
				file.SetCellValue(sheetName, cell, *grade)
				totalGrades += *grade
				gradeCount++
			} else {
				file.SetCellValue(sheetName, cell, "N/A")
			}
		}

		averageCell := string(rune('C'+len(subjects))) + strconv.Itoa(rowNum)
		if gradeCount > 0 {
			file.SetCellValue(sheetName, averageCell, totalGrades/float64(gradeCount))
		} else {
			file.SetCellValue(sheetName, averageCell, "N/A")
		}
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment;filename=teacher_grades.xlsx")

	// Write the file to the response
	if err := file.Write(w); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		log.Printf("Error writing file: %v\n", err)
		return
	}
}

// Helper function to get the grade for a student in a specific subject and term
func getGradeForStudent(gradesData models.GradesData, studentID, subjectID, termID int) *float64 {
	for _, studentGrades := range gradesData.Grades {
		if studentGrades.StudentID == studentID && studentGrades.SubjectID == subjectID {
			for _, termGrades := range studentGrades.Terms {
				if termGrades.Term == strconv.Itoa(termID) { // Compare with the term string
					if len(termGrades.Grades) > 0 {
						grade := float64(termGrades.Grades[0].Grade) // Convert float32 to float64
						return &grade
					}
				}
			}
		}
	}
	log.Printf("No grade found for student %d, subject %d, term %d\n", studentID, subjectID, termID)
	return nil
}
