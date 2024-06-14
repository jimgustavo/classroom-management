// database/subject.go

package database

import (
	"errors"

	"github.com/jimgustavo/classroom-management/models"
)

// CreateSubject inserts a new subject record into the database
func CreateSubject(subject *models.Subject) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO subjects (name, teacher_id) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, subject.Name, subject.TeacherID).Scan(&subject.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllSubjects retrieves all subjects from the database
func GetAllSubjects() ([]models.Subject, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name, teacher_id FROM subjects"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.ID, &subject.Name, &subject.TeacherID)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

// GetSubjectByID retrieves a specific subject by its ID from the database
func GetSubjectByID(id int) (*models.Subject, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var subject models.Subject
	query := "SELECT name, teacher_id FROM subjects WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&subject.Name, &subject.TeacherID)
	if err != nil {
		return nil, err
	}
	subject.ID = id

	return &subject, nil
}

// UpdateSubject updates the details of a specific subject in the database
func UpdateSubject(id int, updatedSubject *models.Subject) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "UPDATE subjects SET name = $1, teacher_id = $2 WHERE id = $3"
	_, err := db.Exec(query, updatedSubject.Name, updatedSubject.TeacherID, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSubject deletes a specific subject from the database
func DeleteSubject(id int) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "DELETE FROM subjects WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetSubjectsByStudentID retrieves all subjects associated with a student by student ID
func GetSubjectsByStudentID(studentID int) ([]models.Subject, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `
		SELECT subjects.id, subjects.name
		FROM subjects
		INNER JOIN student_subjects ON subjects.id = student_subjects.subject_id
		WHERE student_subjects.student_id = $1
	`
	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.ID, &subject.Name)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

// GetStudentsBySubjectID retrieves all students associated with a subject by subject ID
func GetStudentsBySubjectID(subjectID int) ([]models.Student, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `
		SELECT students.id, students.name
		FROM students
		INNER JOIN student_subjects ON students.id = student_subjects.student_id
		WHERE student_subjects.subject_id = $1
	`
	rows, err := db.Query(query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.Name)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// RemoveGradeLabelFromSubjectByTerm removes a grade label from a subject for a specific term in the database
func RemoveGradeLabelFromSubjectByTerm(subjectID, gradeLabelID, termID int) error {
	_, err := db.Exec("DELETE FROM grade_labels_subjects WHERE subject_id = $1 AND grade_label_id = $2 AND term_id = $3", subjectID, gradeLabelID, termID)
	if err != nil {
		return err
	}
	return nil
}
