package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jimgustavo/classroom-management/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// CreateTeacher inserts a new teacher record into the database with a hashed password
func CreateTeacher(name, email, password string) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `INSERT INTO teachers (name, email, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, name, email, password)
	if err != nil {
		log.Println("Error storing teacher in the database:", err)
	}
	return err
}

// AuthenticateTeacher verifies the teacher's email and password, and returns the teacher's ID
func AuthenticateTeacher(email, password string) (int, error) {
	if db == nil {
		return 0, errors.New("database connection is not initialized")
	}

	var (
		id             int
		hashedPassword string
	)

	query := `SELECT id, password FROM teachers WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Teacher not found")
			return 0, errors.New("teacher not found")
		}
		log.Println("Error retrieving teacher from database:", err)
		return 0, err
	}

	log.Println("Provided password:", password)
	log.Println("Retrieved hashed password:", hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Error comparing password hash:", err)
		return 0, errors.New("invalid credentials")
	}

	log.Println("Authentication successful for teacher with ID:", id)

	return id, nil
}

// GetTeacherByID retrieves a teacher's details by their ID
func GetTeacherByID(teacherID int) (*models.Teacher, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var teacher models.Teacher
	query := `SELECT id, name, email FROM teachers WHERE id = $1`
	err := db.QueryRow(query, teacherID).Scan(&teacher.ID, &teacher.Name, &teacher.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("teacher not found")
		}
		return nil, err
	}

	return &teacher, nil
}

// GetAllTeachers retrieves all teachers from the database without requiring authorization
func GetAllTeachers() ([]models.Teacher, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name, email FROM teachers"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Email)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teachers, nil
}

// DeleteTeacher deletes a specific teacher from the database without requiring authorization
func DeleteTeacher(id string) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "DELETE FROM teachers WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

/*
func AuthenticateTeacher(email, password string) (int, error) {
	if db == nil {
		return 0, errors.New("database connection is not initialized")
	}

	var (
		id             int
		hashedPassword string
	)

	// Retrieve the hashed password and ID from the database
	query := `SELECT id, password FROM teachers WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Teacher not found")
			return 0, errors.New("teacher not found")
		}
		log.Println("Error retrieving teacher from database:", err)
		return 0, err
	}

	log.Println("Email:", email)
	log.Println("Generated hashed password:", string(hashedPassword))
	log.Println("Retrieved hashed password:", hashedPassword)
	log.Println("Generated hashed password length:", len(hashedPassword))
	log.Println("Retrieved hashed password length:", len(hashedPassword))

	// Trim both hashed passwords
	trimmedGeneratedHashedPassword := strings.TrimSpace(hashedPassword)
	trimmedProvidedPassword := strings.TrimSpace(password)

	log.Println("Trimmed generated hashed password:", trimmedGeneratedHashedPassword)
	log.Println("Trimmed provided password:", trimmedProvidedPassword)

	// Compare the trimmed hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(trimmedGeneratedHashedPassword), []byte(trimmedProvidedPassword))
	if err != nil {
		log.Println("Error comparing password hash:", err)
		return 0, errors.New("invalid credentials")
	}

	log.Println("Authentication successful for teacher with ID:", id)

	return id, nil
}
*/
