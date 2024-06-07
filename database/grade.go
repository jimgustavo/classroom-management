// database/grade.go

package database

import (
	"fmt"

	"github.com/jimgustavo/classroom-management/models"
)

// InsertGradesInClassroom inserts a grade into the database
func InsertGradesInClassroom(studentID, subjectID int, term string, labelID int, grade float32, classroomID int) error {
	db := GetDB()

	// Check if term exists, if not insert it
	var termID int
	err := db.QueryRow(`SELECT id FROM terms WHERE name = $1`, term).Scan(&termID)
	if err != nil {
		err = db.QueryRow(`INSERT INTO terms (name) VALUES ($1) RETURNING id`, term).Scan(&termID)
		if err != nil {
			return fmt.Errorf("failed to insert term: %w", err)
		}
	}

	query := `
        INSERT INTO grades (student_id, subject_id, term_id, label_id, grade, classroom_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (student_id, subject_id, term_id, label_id) DO UPDATE
        SET grade = $5`

	_, err = db.Exec(query, studentID, subjectID, termID, labelID, grade, classroomID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// FetchGradesByClassroomID retrieves all the grades from the database for a specific classroom
func FetchGradesByClassroomID(classroomID int) (models.GradesData, error) {
	db := GetDB()

	query := `
        SELECT g.student_id, g.subject_id, t.name, g.label_id, g.grade
        FROM grades g
        JOIN terms t ON g.term_id = t.id
        WHERE g.classroom_id = $1`

	rows, err := db.Query(query, classroomID)
	if err != nil {
		return models.GradesData{}, fmt.Errorf("error fetching grades: %w", err)
	}
	defer rows.Close()

	var gradesData models.GradesData
	gradeMap := make(map[int]map[int]map[string][]models.Grade)

	for rows.Next() {
		var studentID, subjectID, labelID int
		var term string
		var grade float32

		if err := rows.Scan(&studentID, &subjectID, &term, &labelID, &grade); err != nil {
			return models.GradesData{}, fmt.Errorf("error scanning row: %w", err)
		}

		if gradeMap[studentID] == nil {
			gradeMap[studentID] = make(map[int]map[string][]models.Grade)
		}
		if gradeMap[studentID][subjectID] == nil {
			gradeMap[studentID][subjectID] = make(map[string][]models.Grade)
		}

		gradeMap[studentID][subjectID][term] = append(gradeMap[studentID][subjectID][term], models.Grade{LabelID: labelID, Grade: grade})
	}

	for studentID, subjects := range gradeMap {
		for subjectID, terms := range subjects {
			var termGrades []models.TermGrades
			for term, grades := range terms {
				termGrades = append(termGrades, models.TermGrades{
					Term:   term,
					Grades: grades,
				})
			}
			gradesData.Grades = append(gradesData.Grades, models.StudentTermGrades{
				StudentID: studentID,
				SubjectID: subjectID,
				Terms:     termGrades,
			})
		}
	}

	return gradesData, nil
}

// FetchGradesByClassroomIDAndTermID retrieves all the grades from the database for a specific classroom and term
func FetchGradesByClassroomIDAndTermID(classroomID, termID int) (models.GradesData, error) {
	db := GetDB()

	query := `
        SELECT g.student_id, g.subject_id, t.name, g.label_id, g.grade
        FROM grades g
        JOIN terms t ON g.term_id = t.id
        WHERE g.classroom_id = $1 AND g.term_id = $2`

	rows, err := db.Query(query, classroomID, termID)
	if err != nil {
		return models.GradesData{}, fmt.Errorf("error fetching grades: %w", err)
	}
	defer rows.Close()

	var gradesData models.GradesData
	gradeMap := make(map[int]map[int]map[string][]models.Grade)

	for rows.Next() {
		var studentID, subjectID, labelID int
		var term string
		var grade float32

		if err := rows.Scan(&studentID, &subjectID, &term, &labelID, &grade); err != nil {
			return models.GradesData{}, fmt.Errorf("error scanning row: %w", err)
		}

		if gradeMap[studentID] == nil {
			gradeMap[studentID] = make(map[int]map[string][]models.Grade)
		}
		if gradeMap[studentID][subjectID] == nil {
			gradeMap[studentID][subjectID] = make(map[string][]models.Grade)
		}

		gradeMap[studentID][subjectID][term] = append(gradeMap[studentID][subjectID][term], models.Grade{LabelID: labelID, Grade: grade})
	}

	for studentID, subjects := range gradeMap {
		for subjectID, terms := range subjects {
			var termGrades []models.TermGrades
			for term, grades := range terms {
				termGrades = append(termGrades, models.TermGrades{
					Term:   term,
					Grades: grades,
				})
			}
			gradesData.Grades = append(gradesData.Grades, models.StudentTermGrades{
				StudentID: studentID,
				SubjectID: subjectID,
				Terms:     termGrades,
			})
		}
	}

	return gradesData, nil
}
