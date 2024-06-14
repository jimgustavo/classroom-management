// database/classroom.go

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jimgustavo/classroom-management/models"
)

// CreateClassroom inserts a new classroom record into the database
func CreateClassroom(classroom *models.Classroom) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO classrooms (name, teacher_id) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, classroom.Name, classroom.TeacherID).Scan(&classroom.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllClassrooms retrieves all classrooms from the database
func GetAllClassrooms() ([]models.Classroom, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `
        SELECT classrooms.id, classrooms.name, teachers.id, teachers.name
        FROM classrooms
        JOIN teachers ON classrooms.teacher_id = teachers.id
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classrooms []models.Classroom
	for rows.Next() {
		var classroom models.Classroom
		var teacher models.Teacher
		err := rows.Scan(&classroom.ID, &classroom.Name, &teacher.ID, &teacher.Name)
		if err != nil {
			return nil, err
		}
		classroom.Teacher = teacher
		classrooms = append(classrooms, classroom)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classrooms, nil
}

// GetClassroomByID retrieves a specific classroom by its ID from the database
func GetClassroomByID(id int) (*models.Classroom, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var classroom models.Classroom
	var teacher models.Teacher
	query := `
        SELECT classrooms.name, teachers.id, teachers.name
        FROM classrooms
        JOIN teachers ON classrooms.teacher_id = teachers.id
        WHERE classrooms.id = $1
    `
	err := db.QueryRow(query, id).Scan(&classroom.Name, &teacher.ID, &teacher.Name)
	if err != nil {
		return nil, err
	}
	classroom.ID = id
	classroom.Teacher = teacher

	return &classroom, nil
}

// UpdateClassroom updates the details of a specific classroom in the database
func UpdateClassroom(id int, updatedClassroom *models.Classroom) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "UPDATE classrooms SET name = $1, teacher_id = $2 WHERE id = $3"
	_, err := db.Exec(query, updatedClassroom.Name, updatedClassroom.TeacherID, id)
	if err != nil {
		return err
	}

	return nil
}

// GetStudentsByClassroomID retrieves all students belonging to a specific classroom by classroom ID
func GetStudentsByClassroomID(classroomID int) ([]models.Student, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name FROM students WHERE classroom_id = $1"
	rows, err := db.Query(query, classroomID)
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
		student.ClassroomID = classroomID // Set the classroom ID for each student
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// AddSubjectToClassroom adds a subject to a classroom
func AddSubjectToClassroom(classroomID, subjectID int) error {
	log.Printf("Attempting to add subject %d to classroom %d", subjectID, classroomID)

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

	query := `
        SELECT subjects.id, subjects.name, 
               COALESCE(json_agg(json_build_object('id', grade_labels.id, 'label', grade_labels.label, 'term_id', grade_labels_subjects.term_id)) FILTER (WHERE grade_labels.id IS NOT NULL), '[]')
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
		var gradeLabels json.RawMessage
		if err := rows.Scan(&subjectID, &subjectName, &gradeLabels); err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}

		var gradeLabelsParsed []models.GradeLabelTerm
		if err := json.Unmarshal(gradeLabels, &gradeLabelsParsed); err != nil {
			log.Println("Failed to unmarshal grade labels:", err)
			return nil, err
		}

		subjects = append(subjects, models.SubjectWithGradeLabels{
			ID:          subjectID,
			Name:        subjectName,
			GradeLabels: gradeLabelsParsed,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println("Error occurred while iterating through rows:", err)
		return nil, err
	}

	return subjects, nil
}

// DeleteClassroom deletes a specific classroom from the database
func DeleteClassroom(id int) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "DELETE FROM classrooms WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// UnrollStudentFromClassroom removes a student from a classroom
func UnrollStudentFromClassroom(classroomID, studentID string) error {
	query := `UPDATE students SET classroom_id = NULL WHERE id = $1 AND classroom_id = $2`
	result, err := db.Exec(query, studentID, classroomID)
	if err != nil {
		log.Printf("error executing query: %s with classroomID: %s, studentID: %s", query, classroomID, studentID)
		return fmt.Errorf("error unrolling student from classroom: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		return fmt.Errorf("error unrolling student from classroom: %w", err)
	}

	log.Printf("rows affected: %d", rowsAffected)
	return nil
}

// RemoveSubjectFromClassroom removes a subject from a classroom
func RemoveSubjectFromClassroom(classroomID, subjectID string) error {
	query := `DELETE FROM classroom_subjects WHERE classroom_id = $1 AND subject_id = $2`
	_, err := db.Exec(query, classroomID, subjectID)
	if err != nil {
		return fmt.Errorf("error removing subject from classroom: %w", err)
	}
	return nil
}
