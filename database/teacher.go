package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jimgustavo/classroom-management/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// CreateTeacher inserts a new teacher record into the database with a hashed password and role
func CreateTeacher(name, email, password, role string) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `INSERT INTO teachers (name, email, password, role) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, name, email, password, role)
	if err != nil {
		log.Println("Error storing teacher in the database:", err)
	}
	return err
}

// AuthenticateTeacher verifies the teacher's email and password, and returns the teacher's ID and role
func AuthenticateTeacher(email, password string) (int, string, error) {
	if db == nil {
		return 0, "", errors.New("database connection is not initialized")
	}

	var (
		id             int
		hashedPassword string
		role           string
	)

	query := `SELECT id, password, role FROM teachers WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&id, &hashedPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Teacher not found")
			return 0, "", errors.New("teacher not found")
		}
		log.Println("Error retrieving teacher from database:", err)
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Error comparing password hash:", err)
		return 0, "", errors.New("invalid credentials")
	}

	return id, role, nil
}

// GetAllTeachers retrieves all teachers from the database without requiring authorization
func GetAllTeachers() ([]models.Teacher, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, name, email, role FROM teachers"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Role)
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

// GetTeacherByID retrieves a teacher's details by their ID
func GetTeacherByID(teacherID int) (*models.Teacher, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var teacher models.Teacher
	query := `SELECT id, name, email, role FROM teachers WHERE id = $1`
	err := db.QueryRow(query, teacherID).Scan(&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("teacher not found")
		}
		return nil, err
	}

	return &teacher, nil
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

// UpdateTeacherRole updates the role of a teacher
func UpdateTeacherRole(teacherID int, role string) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `UPDATE teachers SET role = $1 WHERE id = $2`
	_, err := db.Exec(query, role, teacherID)
	if err != nil {
		log.Println("Error updating teacher role in the database:", err)
	}
	return err
}

// ////////////////////TEACHER DATA////////////////////
func CreateTeacherData(db *sql.DB, teacherData models.TeacherData) (int, error) {
	query := `
        INSERT INTO teacher_data (school, school_year, school_hours, country, city, teacher_id, teacher_full_name, id_number, labor_dependency_relationship, principal, vice_principal, dece, inspector, institutional_email, phone, teacher_birthday)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        ON CONFLICT (teacher_id) 
        DO UPDATE SET 
            school = EXCLUDED.school,
            school_year = EXCLUDED.school_year,
            school_hours = EXCLUDED.school_hours,
            country = EXCLUDED.country,
            city = EXCLUDED.city,
            teacher_full_name = EXCLUDED.teacher_full_name,
            id_number = EXCLUDED.id_number,
            labor_dependency_relationship = EXCLUDED.labor_dependency_relationship,
            principal = EXCLUDED.principal,
            vice_principal = EXCLUDED.vice_principal,
            dece = EXCLUDED.dece,
            inspector = EXCLUDED.inspector,
            institutional_email = EXCLUDED.institutional_email,
            phone = EXCLUDED.phone,
            teacher_birthday = EXCLUDED.teacher_birthday
        RETURNING id;
    `
	var id int
	err := db.QueryRow(query, teacherData.School, teacherData.SchoolYear, teacherData.SchoolHours, teacherData.Country, teacherData.City, teacherData.TeacherID, teacherData.TeacherFullName, teacherData.TeacherIDNumber, teacherData.LaborDependencyRelationship, teacherData.Principal, teacherData.VicePrincipal, teacherData.Dece, teacherData.Inspector, teacherData.InstitutionalEmail, teacherData.Phone, teacherData.TeacherBirthday).Scan(&id)
	return id, err
}

func GetAllTeacherData(db *sql.DB) ([]models.TeacherData, error) {
	query := `SELECT id, school, school_year, school_hours, country, city, teacher_id, teacher_full_name, id_number, labor_dependency_relationship, principal, vice_principal, dece, inspector, institutional_email, phone, teacher_birthday FROM teacher_data`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teacherDataList []models.TeacherData
	for rows.Next() {
		var teacherData models.TeacherData
		err := rows.Scan(&teacherData.ID, &teacherData.School, &teacherData.SchoolYear, &teacherData.SchoolHours, &teacherData.Country, &teacherData.City, &teacherData.TeacherID, &teacherData.TeacherFullName, &teacherData.TeacherIDNumber, &teacherData.LaborDependencyRelationship, &teacherData.Principal, &teacherData.VicePrincipal, &teacherData.Dece, &teacherData.Inspector, &teacherData.InstitutionalEmail, &teacherData.Phone, &teacherData.TeacherBirthday)
		if err != nil {
			return nil, err
		}
		teacherDataList = append(teacherDataList, teacherData)
	}

	return teacherDataList, nil
}

func GetTeacherDataByID(db *sql.DB, id int) (models.TeacherData, error) {
	query := `SELECT id, school, school_year, school_hours, country, city, teacher_id, teacher_full_name, id_number, labor_dependency_relationship, principal, vice_principal, dece, inspector, institutional_email, phone, teacher_birthday FROM teacher_data WHERE id = $1`
	var teacherData models.TeacherData
	err := db.QueryRow(query, id).Scan(&teacherData.ID, &teacherData.School, &teacherData.SchoolYear, &teacherData.SchoolHours, &teacherData.Country, &teacherData.City, &teacherData.TeacherID, &teacherData.TeacherFullName, &teacherData.TeacherIDNumber, &teacherData.LaborDependencyRelationship, &teacherData.Principal, &teacherData.VicePrincipal, &teacherData.Dece, &teacherData.Inspector, &teacherData.InstitutionalEmail, &teacherData.Phone, &teacherData.TeacherBirthday)
	if err != nil {
		return teacherData, err
	}
	return teacherData, nil
}

// GetTeacherDataByTeacherID retrieves teacher data by teacher ID
func GetTeacherDataByTeacherID(teacherID int) (*models.TeacherData, error) {
	var teacherData models.TeacherData

	query := `SELECT id, school, school_year, school_hours, country, city, teacher_id, teacher_full_name, id_number, labor_dependency_relationship, principal, vice_principal, dece, inspector, institutional_email, phone, teacher_birthday FROM teacher_data WHERE teacher_id = $1`

	row := db.QueryRow(query, teacherID)
	err := row.Scan(&teacherData.ID, &teacherData.School, &teacherData.SchoolYear, &teacherData.SchoolHours, &teacherData.Country, &teacherData.City, &teacherData.TeacherID, &teacherData.TeacherFullName, &teacherData.TeacherIDNumber, &teacherData.LaborDependencyRelationship, &teacherData.Principal, &teacherData.VicePrincipal, &teacherData.Dece, &teacherData.Inspector, &teacherData.InstitutionalEmail, &teacherData.Phone, &teacherData.TeacherBirthday)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no teacher data found for teacher ID %d", teacherID)
		}
		return nil, err
	}

	return &teacherData, nil
}

func UpdateTeacherData(db *sql.DB, teacherData models.TeacherData) error {
	query := `UPDATE teacher_data SET school = $1, school_year = $2, school_hours = $3, country = $4, city = $5, teacher_id = $6, teacher_full_name = $7, id_number = $8, labor_dependency_relationship = $9, principal = $10, vice_principal = $11, dece = $12, inspector = $13, institutional_email = $14, phone = $15, teacher_birthday = $16 WHERE id = $17`
	_, err := db.Exec(query, teacherData.School, teacherData.SchoolYear, teacherData.SchoolHours, teacherData.Country, teacherData.City, teacherData.TeacherID, teacherData.TeacherFullName, teacherData.TeacherIDNumber, teacherData.LaborDependencyRelationship, teacherData.Principal, teacherData.VicePrincipal, teacherData.Dece, teacherData.Inspector, teacherData.InstitutionalEmail, teacherData.Phone, teacherData.TeacherBirthday, teacherData.ID)
	return err
}

func DeleteTeacherData(db *sql.DB, id int) error {
	query := `DELETE FROM teacher_data WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func SaveLogo(teacherID int, fileBytes []byte) error {
    query := `
        INSERT INTO logos (teacher_id, logo)
        VALUES ($1, $2)
        ON CONFLICT (teacher_id)
        DO UPDATE SET logo = EXCLUDED.logo`
    _, err := db.Exec(query, teacherID, fileBytes)
    if err != nil {
        fmt.Println("Database error:", err)
    }
    return err
}


func GetLogo(teacherID int) ([]byte, error) {
    var logo []byte
    query := `SELECT logo FROM logos WHERE teacher_id = $1`
    err := db.QueryRow(query, teacherID).Scan(&logo)
    return logo, err
}