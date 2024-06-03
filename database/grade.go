// database/grade.go

package database

import (
	"fmt"
	"log"

	"github.com/jimgustavo/classroom-management/models"
)

// InsertGrade inserts a grade into the database
func InsertGradesInClassroom(studentID, subjectID int, label, grade string, classroomID int) error {
	db := GetDB()

	query := `
        INSERT INTO grades (student_id, subject_id, label, grade, classroom_id)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (student_id, subject_id, label) DO UPDATE
        SET grade = $4`

	log.Printf("Executing query: %s with values studentID=%d, subjectID=%d, label=%s, grade=%s, classroomID=%d",
		query, studentID, subjectID, label, grade, classroomID)

	_, err := db.Exec(query, studentID, subjectID, label, grade, classroomID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// FetchGradesFromDB retrieves all the grades from the database
func FetchGradesByClassroomID(classroomID int) (models.GradesData, error) {
	db := GetDB()

	query := `
        SELECT student_id, subject_id, label, grade
        FROM grades
        WHERE classroom_id = $1`

	rows, err := db.Query(query, classroomID)
	if err != nil {
		return models.GradesData{}, fmt.Errorf("error fetching grades: %w", err)
	}
	defer rows.Close()

	var gradesData models.GradesData
	gradeMap := make(map[int]map[int][]models.Grade)

	for rows.Next() {
		var studentID, subjectID int
		var label, grade string

		if err := rows.Scan(&studentID, &subjectID, &label, &grade); err != nil {
			return models.GradesData{}, fmt.Errorf("error scanning row: %w", err)
		}

		if gradeMap[studentID] == nil {
			gradeMap[studentID] = make(map[int][]models.Grade)
		}

		gradeMap[studentID][subjectID] = append(gradeMap[studentID][subjectID], models.Grade{Label: label, Grade: grade})
	}

	for studentID, subjects := range gradeMap {
		for subjectID, grades := range subjects {
			gradesData.Grades = append(gradesData.Grades, models.StudentGrade{
				StudentID: studentID,
				SubjectID: subjectID,
				Grades:    grades,
			})
		}
	}

	return gradesData, nil
}
