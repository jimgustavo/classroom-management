// database.go

package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jimgustavo/classroom-management/models"
	"github.com/lib/pq"
)

var db *sql.DB

// InitializeDB initializes the database connection
func InitializeDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

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

/*
// GetSubjectsByClassroomID retrieves subjects assigned to a classroom by classroom ID
func GetSubjectsByClassroomID(classroomID string) ([]models.Subject, error) {
	// Initialize a slice to hold the subjects
	var subjects []models.Subject

	// Query to retrieve subjects assigned to the specified classroom
	query := `
		SELECT s.id, s.name
		FROM subjects s
		INNER JOIN student_subjects ss ON s.id = ss.subject_id
		INNER JOIN students st ON ss.student_id = st.id
		WHERE st.classroom_id = $1
	`

	// Execute the query
	rows, err := db.Query(query, classroomID)
	if err != nil {
		// Handle query error
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var subject models.Subject
		// Scan subject data from the row
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			// Handle scanning error
			return nil, err
		}
		// Append subject to the slice
		subjects = append(subjects, subject)
	}
	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		// Handle iteration error
		return nil, err
	}

	return subjects, nil
}
*/
// GetSubjectsAndStudentsByClassroomID retrieves subjects assigned to a classroom along with students assigned to each subject by classroom ID
func GetSubjectsAndStudentsByClassroomID(classroomID string) ([]models.SubjectWithStudents, error) {
	// Initialize a slice to hold the subjects along with students assigned to each subject
	var subjectsWithStudents []models.SubjectWithStudents

	// Query to retrieve subjects assigned to the specified classroom along with students assigned to each subject
	query := `
        SELECT s.id AS subject_id, s.name AS subject_name,
               st.id AS student_id, st.name AS student_name
        FROM subjects s
        INNER JOIN student_subjects ss ON s.id = ss.subject_id
        INNER JOIN students st ON ss.student_id = st.id
        WHERE st.classroom_id = $1
        ORDER BY s.id, st.id
    `

	// Execute the query
	rows, err := db.Query(query, classroomID)
	if err != nil {
		// Handle query error
		return nil, err
	}
	defer rows.Close()

	var currentSubjectID string
	var subjectWithStudents models.SubjectWithStudents

	// Iterate over the rows
	for rows.Next() {
		var subjectID string
		var subjectName string
		var studentID sql.NullInt64
		var studentName sql.NullString

		// Scan data from the row
		if err := rows.Scan(&subjectID, &subjectName, &studentID, &studentName); err != nil {
			// Handle scanning error
			return nil, err
		}

		// Check if this row represents a new subject
		if subjectID != currentSubjectID {
			// If it's not the first subject, append the previous subjectWithStudents to the slice
			if currentSubjectID != "" {
				subjectsWithStudents = append(subjectsWithStudents, subjectWithStudents)
			}
			// Initialize a new subjectWithStudents
			subjectWithStudents = models.SubjectWithStudents{
				ID:   subjectID,
				Name: subjectName,
			}
			// Update currentSubjectID
			currentSubjectID = subjectID
		}

		// If studentID and studentName are valid, append student to the current subjectWithStudents
		if studentID.Valid && studentName.Valid {
			subjectWithStudents.Students = append(subjectWithStudents.Students, models.Student{
				ID:   int(studentID.Int64),
				Name: studentName.String,
			})
		}
	}

	// Append the last subjectWithStudents to the slice
	if currentSubjectID != "" {
		subjectsWithStudents = append(subjectsWithStudents, subjectWithStudents)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		// Handle iteration error
		return nil, err
	}

	return subjectsWithStudents, nil
}

/*
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
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
*/

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

/*
// GetAllStudentsWithClassroomAndSubjects retrieves all students along with their assigned classroom and subjects
func GetAllStudentsWithClassroomAndSubjects() ([]models.StudentWithClassroomAndSubjects, error) {
	var studentsWithClassroomAndSubjects []models.StudentWithClassroomAndSubjects

	// Query to retrieve all students along with their assigned classroom and subjects
	query := `
        SELECT s.id, s.name, c.name AS classroom, array_agg(sub.id) AS assigned_subjects
        FROM students s
        INNER JOIN classrooms c ON s.classroom_id = c.id
        LEFT JOIN student_subjects ss ON s.id = ss.student_id
        LEFT JOIN subjects sub ON ss.subject_id = sub.id
        GROUP BY s.id, c.name
    `

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var student models.StudentWithClassroomAndSubjects
		var subjectIDs []int

		// Scan data from the row
		if err := rows.Scan(&student.ID, &student.Name, &student.Classroom, pq.Array(&subjectIDs)); err != nil {
			return nil, err
		}

		// Convert subject IDs to Subject structs
		for _, id := range subjectIDs {
			student.AssignedSubjects = append(student.AssignedSubjects, models.Subject{ID: id})
		}

		// Append student to the slice
		studentsWithClassroomAndSubjects = append(studentsWithClassroomAndSubjects, student)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return studentsWithClassroomAndSubjects, nil
}
*/

// GetAllStudentsWithClassroomAndSubjects retrieves all students along with their assigned classroom and subjects
func GetAllStudentsWithClassroomAndSubjects() ([]models.StudentWithClassroomAndSubjects, error) {
	var studentsWithClassroomAndSubjects []models.StudentWithClassroomAndSubjects

	// Query to retrieve all students along with their assigned classroom and subjects

	query := `
    SELECT 
        s.id, 
        s.name, 
        c.name AS classroom, 
        array_agg(sub.name) AS assigned_subjects
    FROM 
        students s
    INNER JOIN 
        classrooms c ON s.classroom_id = c.id
    LEFT JOIN 
        student_subjects ss ON s.id = ss.student_id
    LEFT JOIN 
        subjects sub ON ss.subject_id = sub.id
    GROUP BY 
        s.id, c.name
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

// AddSubjectToStudents adds a subject to all students in a classroom
func AddSubjectToStudents(classroomID, subjectID int) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `
		INSERT INTO student_subjects (student_id, subject_id)
		SELECT students.id, $1
		FROM students
		WHERE students.classroom_id = $2
	`
	_, err := db.Exec(query, subjectID, classroomID)
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

// DeleteSubjectFromStudents removes a subject from all students in a classroom
func DeleteSubjectFromStudents(classroomID, subjectID int) error {
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
		AND subject_id = $2
	`
	_, err := db.Exec(query, classroomID, subjectID)
	if err != nil {
		return err
	}

	return nil
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

// CreateGradeLabel inserts a new grade label into the database
func CreateGradeLabel(gradeLabel *models.GradeLabel) error {
	query := "INSERT INTO grade_labels (subject_id, label) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, gradeLabel.SubjectID, gradeLabel.Label).Scan(&gradeLabel.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllGradeLabels retrieves all grade labels from the database
func GetAllGradeLabels() ([]models.GradeLabel, error) {
	query := "SELECT id, subject_id, label FROM grade_labels"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gradeLabels []models.GradeLabel
	for rows.Next() {
		var gradeLabel models.GradeLabel
		if err := rows.Scan(&gradeLabel.ID, &gradeLabel.SubjectID, &gradeLabel.Label); err != nil {
			return nil, err
		}
		gradeLabels = append(gradeLabels, gradeLabel)
	}

	return gradeLabels, nil
}

// GetGradeLabelByID retrieves a specific grade label by its ID from the database
func GetGradeLabelByID(id string) (*models.GradeLabel, error) {
	var gradeLabel models.GradeLabel
	query := "SELECT id, subject_id, label FROM grade_labels WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&gradeLabel.ID, &gradeLabel.SubjectID, &gradeLabel.Label)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("grade label not found")
		}
		return nil, err
	}
	return &gradeLabel, nil
}

// UpdateGradeLabel updates an existing grade label in the database
func UpdateGradeLabel(gradeLabel *models.GradeLabel) error {
	query := "UPDATE grade_labels SET subject_id = $1, label = $2 WHERE id = $3"
	_, err := db.Exec(query, gradeLabel.SubjectID, gradeLabel.Label, gradeLabel.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGradeLabel deletes a grade label from the database by its ID
func DeleteGradeLabel(id string) error {
	query := "DELETE FROM grade_labels WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// AssignGradeLabelToSubject assigns a grade label to a subject in the database
func AssignGradeLabelToSubject(subjectID string, gradeLabelID int) error {
	_, err := db.Exec("INSERT INTO grade_labels_subjects (subject_id, grade_label_id) VALUES ($1, $2)", subjectID, gradeLabelID)
	if err != nil {
		return err
	}
	return nil
}

// GetGradeLabelsForSubject retrieves all grade labels assigned to a subject from the database
func GetGradeLabelsForSubject(subjectID string) ([]models.GradeLabel, error) {
	var gradeLabels []models.GradeLabel

	rows, err := db.Query("SELECT id, subject_id, label FROM grade_labels WHERE id IN (SELECT grade_label_id FROM grade_labels_subjects WHERE subject_id = $1)", subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gradeLabel models.GradeLabel
		if err := rows.Scan(&gradeLabel.ID, &gradeLabel.SubjectID, &gradeLabel.Label); err != nil {
			return nil, err
		}
		gradeLabels = append(gradeLabels, gradeLabel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gradeLabels, nil
}

func AddSubjectToClassroom(classroomID, subjectID int, req models.AddSubjectRequest) error {
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

	// Add grade labels to the subject
	for _, gradeLabelID := range req.GradeLabelIDs {
		_, err = tx.Exec("INSERT INTO grade_labels_subjects (subject_id, grade_label_id) VALUES ($1, $2)", subjectID, gradeLabelID)
		if err != nil {
			log.Println("Failed to execute SQL query:", err)
			return err
		}
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

	log.Println("Subjects with grade labels retrieved successfully:", subjects)
	return subjects, nil
}

// AddGrade adds a grade to the database for a specific student in a specific subject and label.
func AddGrade(grade models.Grade) error {
	log.Printf("Adding grade %+v", grade)

	// Check if the student ID exists in the students table
	var studentIDExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", grade.StudentID).Scan(&studentIDExists)
	if err != nil {
		log.Println("Failed to check if student ID exists:", err)
		return err
	}
	if !studentIDExists {
		return errors.New("student ID does not exist")
	}

	// Perform database operations to add the grade
	_, err = db.Exec("INSERT INTO grades (student_id, subject_id, label_id, grade) VALUES ($1, $2, $3, $4)",
		grade.StudentID, grade.SubjectID, grade.LabelID, grade.Grade)
	if err != nil {
		log.Println("Failed to execute SQL query:", err)
		return err
	}

	log.Println("Grade added successfully")
	return nil
}

// GetAllStudentsWithGrades retrieves all students with their grades, labels, subjects, and classrooms from the database.
func GetAllStudentsWithGrades() ([]models.StudentGradeInfo, error) {
	// Execute the SQL query to fetch all students with their grades, labels, subjects, and classrooms
	rows, err := db.Query(`
	SELECT
    students.id AS student_id,
    students.name AS student_name,
    COALESCE(grades.grade, 0.0) AS grade,
    COALESCE(grade_labels.label, 'no label added yet') AS label,
	COALESCE(subjects.name, 'no subject added yet') AS subject,
    classrooms.name AS classroom
FROM
    students
LEFT JOIN
    grades ON students.id = grades.student_id
LEFT JOIN
    grade_labels ON grades.label_id = grade_labels.id
LEFT JOIN
    subjects ON grades.subject_id = subjects.id
LEFT JOIN
    classrooms ON students.classroom_id = classrooms.id

	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store student information with multiple grades and subjects
	studentMap := make(map[int]*models.StudentGradeInfo)

	// Iterate through the rows and populate StudentGradeInfo structs
	for rows.Next() {
		var studentID int
		var studentName, label, subject, classroom string
		var grade float64
		err := rows.Scan(&studentID, &studentName, &grade, &label, &subject, &classroom)
		if err != nil {
			return nil, err
		}

		// If the student already exists in the map, append the grade info
		if student, ok := studentMap[studentID]; ok {
			student.Grades = append(student.Grades, models.GradeInfo{
				Grade:     grade,
				Label:     label,
				Subject:   subject,
				Classroom: classroom,
			})
		} else {
			// Otherwise, create a new entry in the map
			studentMap[studentID] = &models.StudentGradeInfo{
				StudentID:   studentID,
				StudentName: studentName,
				Grades: []models.GradeInfo{
					{
						Grade:     grade,
						Label:     label,
						Subject:   subject,
						Classroom: classroom,
					},
				},
			}
		}
	}

	// Convert the map to a slice of StudentGradeInfo structs
	var studentGradeInfo []models.StudentGradeInfo
	for _, student := range studentMap {
		studentGradeInfo = append(studentGradeInfo, *student)
	}

	return studentGradeInfo, nil
}

// GetGradeByStudentID retrieves the grade, label, subject, and classroom by providing the student ID.
func GetGradeByStudentID(studentID int) (*models.GradeInfo, error) {
	// Prepare the SQL query to retrieve grade information
	query := `
		SELECT grades.grade, grade_labels.label, subjects.name, classrooms.name
		FROM grades
		INNER JOIN grade_labels ON grades.label_id = grade_labels.id
		INNER JOIN subjects ON grades.subject_id = subjects.id
		INNER JOIN students ON grades.student_id = students.id
		INNER JOIN classrooms ON students.classroom_id = classrooms.id
		WHERE students.id = $1;
	`

	// Execute the query and retrieve the row
	row := db.QueryRow(query, studentID)

	// Initialize variables to store the retrieved data
	var grade float64
	var label, subject, classroom string

	// Scan the row to extract the values
	err := row.Scan(&grade, &label, &subject, &classroom)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No grade information found for student:", studentID)
			return nil, err
		}
		log.Println("Failed to scan row:", err)
		return nil, err
	}

	// Create a GradeInfo struct to hold the retrieved data
	gradeInfo := &models.GradeInfo{
		Grade:     grade,
		Label:     label,
		Subject:   subject,
		Classroom: classroom,
	}

	return gradeInfo, nil
}
