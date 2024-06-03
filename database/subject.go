// database/subject.go

package database

import (
	"errors"
	"log"

	"github.com/jimgustavo/classroom-management/models"
	"github.com/lib/pq"
)

// CreateSubject inserts a new subject record into the database
func CreateSubject(subject *models.Subject) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO subjects (name) VALUES ($1) RETURNING id"
	err := db.QueryRow(query, subject.Name).Scan(&subject.ID)
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

	query := "SELECT id, name FROM subjects"
	rows, err := db.Query(query)
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

// GetSubjectByID retrieves a specific subject by its ID from the database
func GetSubjectByID(id int) (*models.Subject, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var subject models.Subject
	query := "SELECT name FROM subjects WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&subject.Name)
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

	query := "UPDATE subjects SET name = $1 WHERE id = $2"
	_, err := db.Exec(query, updatedSubject.Name, id)
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

// DeleteSubjectsByClassroomID removes all subjects associated with a classroom by classroom ID
func DeleteSubjectsByClassroomID(classroomID int) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `
		DELETE FROM student_subjects
		WHERE student_id IN (
			SELECT id
			FROM students
			WHERE classroom_id = $1
		)
	`
	_, err := db.Exec(query, classroomID)
	if err != nil {
		return err
	}

	return nil
}

func AddSubjectToClassroom(classroomID, subjectID int) error {
	log.Printf("Attempting to add subject %d to classroom %d", subjectID, classroomID) // Log before database operation

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if there's an error
			tx.Rollback()
		} else {
			// Commit the transaction if successful
			tx.Commit()
		}
	}()

	// Add subject to classroom
	_, err = tx.Exec("INSERT INTO classroom_subjects (classroom_id, subject_id) VALUES ($1, $2)", classroomID, subjectID)
	if err != nil {
		log.Println("Failed to execute SQL query:", err)
		return err
	}

	return nil
}

func GetSubjectsInClassroom(classroomID int) ([]models.SubjectWithGradeLabels, error) {
	log.Printf("Retrieving subjects with grade labels for classroom %d", classroomID)

	// Query to get subjects and their associated grade labels in classroom
	query := `
        SELECT subjects.id, subjects.name, COALESCE(ARRAY_AGG(COALESCE(grade_labels.label, '')), '{}')
        FROM subjects
        LEFT JOIN classroom_subjects ON subjects.id = classroom_subjects.subject_id
        LEFT JOIN grade_labels_subjects ON subjects.id = grade_labels_subjects.subject_id
        LEFT JOIN grade_labels ON grade_labels.id = grade_labels_subjects.grade_label_id
        WHERE classroom_subjects.classroom_id = $1
        GROUP BY subjects.id, subjects.name
    `

	rows, err := db.Query(query, classroomID)
	if err != nil {
		log.Println("Failed to execute SQL query:", err)
		return nil, err
	}
	defer rows.Close()

	var subjects []models.SubjectWithGradeLabels
	for rows.Next() {
		var subjectID int
		var subjectName string
		var gradeLabels pq.StringArray
		if err := rows.Scan(&subjectID, &subjectName, &gradeLabels); err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}

		subjects = append(subjects, models.SubjectWithGradeLabels{
			ID:          subjectID,
			Name:        subjectName,
			GradeLabels: gradeLabels,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println("Error occurred while iterating through rows:", err)
		return nil, err
	}

	//log.Println("Subjects with grade labels retrieved successfully:", subjects)
	return subjects, nil
}
