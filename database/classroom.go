// database/classroom.go

package database

import (
	"errors"

	"github.com/jimgustavo/classroom-management/models"
)

// CreateClassroom inserts a new classroom record into the database
func CreateClassroom(classroom *models.Classroom) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO classrooms (name) VALUES ($1) RETURNING id"
	err := db.QueryRow(query, classroom.Name).Scan(&classroom.ID)
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

	query := "SELECT id, name FROM classrooms"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classrooms []models.Classroom
	for rows.Next() {
		var classroom models.Classroom
		err := rows.Scan(&classroom.ID, &classroom.Name)
		if err != nil {
			return nil, err
		}
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
	query := "SELECT name FROM classrooms WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&classroom.Name)
	if err != nil {
		return nil, err
	}
	classroom.ID = id

	return &classroom, nil
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

// UpdateClassroom updates the details of a specific classroom in the database
func UpdateClassroom(id int, updatedClassroom *models.Classroom) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "UPDATE classrooms SET name = $1 WHERE id = $2"
	_, err := db.Exec(query, updatedClassroom.Name, id)
	if err != nil {
		return err
	}

	return nil
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
