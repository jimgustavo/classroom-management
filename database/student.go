// database/student.go

package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jimgustavo/classroom-management/models"
	"github.com/lib/pq"
)

// CreateStudent inserts a new student record into the database
func CreateStudent(student *models.Student) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO students (name, classroom_id) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, student.Name, student.ClassroomID).Scan(&student.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllStudents retrieves all students from the database
func GetAllStudents() ([]models.Student, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name, classroom_id FROM students"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		var classroomID sql.NullInt64
		err := rows.Scan(&student.ID, &student.Name, &classroomID)
		if err != nil {
			return nil, err
		}
		if classroomID.Valid {
			student.ClassroomID = int(classroomID.Int64)
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// GetStudentByID retrieves a specific student by their ID from the database
func GetStudentByID(id int) (*models.Student, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var student models.Student
	var classroomID sql.NullInt64
	query := "SELECT name, classroom_id FROM students WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&student.Name, &classroomID)
	if err != nil {
		return nil, err
	}
	student.ID = id
	if classroomID.Valid {
		student.ClassroomID = int(classroomID.Int64)
	}

	return &student, nil
}

// GetAllStudentsWithClassroomAndSubjects retrieves all students along with their assigned classroom and subjects
func GetAllStudentsWithClassroomAndSubjects() ([]models.StudentWithClassroomAndSubjects, error) {
	var studentsWithClassroomAndSubjects []models.StudentWithClassroomAndSubjects

	query := `
    SELECT 
    s.id, 
    s.name, 
    COALESCE(c.name, 'No classroom assigned') AS classroom, 
    ARRAY_AGG(COALESCE(sub.name, 'No subject assigned')) AS assigned_subjects
	FROM 
		students s
	LEFT JOIN 
		classrooms c ON s.classroom_id = c.id
	LEFT JOIN 
		student_subjects ss ON s.id = ss.student_id
	LEFT JOIN 
		subjects sub ON ss.subject_id = sub.id
	GROUP BY 
		s.id, COALESCE(c.name, 'No classroom assigned')
	`

	// Log the SQL query
	log.Println("Executing SQL query:", query)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error executing SQL query:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var student models.StudentWithClassroomAndSubjects
		var assignedSubjects pq.StringArray // Use pq.StringArray for scanning PostgreSQL arrays

		// Scan data from the row
		if err := rows.Scan(&student.ID, &student.Name, &student.Classroom, &assignedSubjects); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// Log the scanned data
		log.Printf("Scanned student ID: %d, Name: %s, Classroom: %s, Assigned Subjects: %v", student.ID, student.Name, student.Classroom, assignedSubjects)

		// Convert pq.StringArray to []string
		student.AssignedSubjects = []string(assignedSubjects)

		// Append student to the slice
		studentsWithClassroomAndSubjects = append(studentsWithClassroomAndSubjects, student)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return studentsWithClassroomAndSubjects, nil
}

// UpdateStudent updates the details of a specific student in the database
func UpdateStudent(id int, updatedStudent *models.Student) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "UPDATE students SET name = $1 WHERE id = $2"
	_, err := db.Exec(query, updatedStudent.Name, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStudent deletes a specific student from the database
func DeleteStudent(id int) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "DELETE FROM students WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// InsertStudent inserts a new student into the database
func InsertStudent(student *models.Student) error {
	query := `INSERT INTO students (name, classroom_id) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, student.Name, student.ClassroomID).Scan(&student.ID)
	if err != nil {
		return err
	}
	return nil
}
