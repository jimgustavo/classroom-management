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

	query := "INSERT INTO students (name, classroom_id, teacher_id) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(query, student.Name, student.ClassroomID, student.TeacherID).Scan(&student.ID)
	if err != nil {
		return err
	}

	return nil
}

// CreateStudentWithID inserts a new student record with a specific ID into the database
func CreateStudentWithID(student *models.Student) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	// Insert student with the provided ID
	query := "INSERT INTO students (id, name, classroom_id, teacher_id) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(query, student.ID, student.Name, student.ClassroomID, student.TeacherID)

	if err != nil {
		log.Printf("Error inserting student with ID: %v", err)
		return err
	}

	return nil
}

func GetStudentsByTeacherID(teacherID int) ([]models.Student, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name, COALESCE(classroom_id, 0) as classroom_id, teacher_id FROM students WHERE teacher_id = $1"
	rows, err := db.Query(query, teacherID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.Name, &student.ClassroomID, &student.TeacherID)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows iteration:", err)
		return nil, err
	}

	log.Println("Students found:", students)
	return students, nil
}

// GetAllStudents retrieves all students from the database
func GetAllStudents() ([]models.Student, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name, classroom_id, teacher_id FROM students"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		var classroomID, teacherID sql.NullInt64
		err := rows.Scan(&student.ID, &student.Name, &classroomID, &teacherID)
		if err != nil {
			return nil, err
		}
		if classroomID.Valid {
			student.ClassroomID = int(classroomID.Int64)
		}
		if teacherID.Valid {
			student.TeacherID = int(teacherID.Int64)
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
	var classroomID, teacherID sql.NullInt64
	query := "SELECT name, classroom_id, teacher_id FROM students WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&student.Name, &classroomID, &teacherID)
	if err != nil {
		return nil, err
	}
	student.ID = id
	if classroomID.Valid {
		student.ClassroomID = int(classroomID.Int64)
	}
	if teacherID.Valid {
		student.TeacherID = int(teacherID.Int64)
	}

	return &student, nil
}

// UpdateStudent updates the details of a specific student in the database
func UpdateStudent(id int, updatedStudent *models.Student) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "UPDATE students SET name = $1, classroom_id = $2, teacher_id = $3 WHERE id = $4"
	_, err := db.Exec(query, updatedStudent.Name, updatedStudent.ClassroomID, updatedStudent.TeacherID, id)
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

// GetAllStudentsWithClassroomAndSubjects retrieves all students with their assigned classrooms and subjects
func GetAllStudentsWithClassroomAndSubjects() ([]models.StudentWithClassroomAndSubjects, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `
		SELECT s.id, s.name, s.classroom_id, s.teacher_id, c.name, array_agg(sub.name)
		FROM students s
		LEFT JOIN classrooms c ON s.classroom_id = c.id
		LEFT JOIN student_subjects ss ON s.id = ss.student_id
		LEFT JOIN subjects sub ON ss.subject_id = sub.id
		GROUP BY s.id, c.name
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.StudentWithClassroomAndSubjects
	for rows.Next() {
		var student models.StudentWithClassroomAndSubjects
		var classroomName sql.NullString
		var subjectNames pq.StringArray
		err := rows.Scan(&student.ID, &student.Name, &student.ClassroomID, &student.TeacherID, &classroomName, &subjectNames)
		if err != nil {
			return nil, err
		}
		if classroomName.Valid {
			student.ClassroomName = classroomName.String
		}
		student.Subjects = subjectNames
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// InsertStudent inserts a new student into the database
func InsertStudent(student *models.Student, teacherID int) error {
	query := `INSERT INTO students (name, classroom_id, teacher_id) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, student.Name, student.ClassroomID, teacherID).Scan(&student.ID)
	if err != nil {
		return err
	}
	return nil
}
