package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/handlers"
	"github.com/jimgustavo/classroom-management/middleware"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize the database connection
	err := database.InitializeDB("postgres://tavito:mamacita@localhost:5432/classroom_management?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()  // Ensure the database connection is closed when the program exits

	// Initialize router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/logout", handlers.Logout)
	router.HandleFunc("/teachers", handlers.GetAllTeachersHandler).Methods("GET")
	router.HandleFunc("/teachers/{id}", handlers.DeleteTeacherHandler).Methods("DELETE")

	// Protected routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.Authenticate)

	// Classroom routes
	apiRouter.HandleFunc("/classrooms", handlers.GetAllClassrooms).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{id}", handlers.GetClassroom).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/subjects", handlers.GetSubjectsInClassroom).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroom_id}/students", handlers.GetStudentsByClassroom).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/grades/get", handlers.GetGradesByClassroomID).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/terms/{termID}/grades", handlers.GetGradesByClassroomIDAndTermID).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/averages", handlers.GetAverageGradesByClassroomID).Methods("GET")
	apiRouter.HandleFunc("/classrooms", handlers.CreateClassroom).Methods("POST")
	apiRouter.HandleFunc("/classrooms/{classroomID}/subject/{subjectID}", handlers.AddSubjectToClassroom).Methods("POST")
	apiRouter.HandleFunc("/classrooms/{classroomID}/grades", handlers.UploadGradesToClassroom).Methods("POST")
	apiRouter.HandleFunc("/classrooms/{classroomID}/upload-students", handlers.UploadStudentsFromExcel).Methods("POST") // for uploading students from an Excel file
	apiRouter.HandleFunc("/classrooms/{id}", handlers.UpdateClassroom).Methods("PUT")
	apiRouter.HandleFunc("/classrooms/{id}", handlers.DeleteClassroom).Methods("DELETE")
	apiRouter.HandleFunc("/classrooms/{classroomID}/students/{studentID}", handlers.UnrollStudentFromClassroom).Methods("DELETE")
	apiRouter.HandleFunc("/classrooms/{classroomID}/subjects/{subjectID}", handlers.RemoveSubjectFromClassroom).Methods("DELETE")

	// Student routes
	apiRouter.HandleFunc("/students", handlers.GetAllStudents).Methods("GET")
	apiRouter.HandleFunc("/students/with-classroom-and-subjects", handlers.GetAllStudentsWithClassroomAndSubjects).Methods("GET")
	apiRouter.HandleFunc("/students/{id}", handlers.GetStudent).Methods("GET")
	apiRouter.HandleFunc("/students/{studentID}/subjects", handlers.GetSubjectsByStudentID).Methods("GET")
	apiRouter.HandleFunc("/students", handlers.CreateStudent).Methods("POST")
	apiRouter.HandleFunc("/students/{id}", handlers.UpdateStudent).Methods("PUT")
	apiRouter.HandleFunc("/students/{id}", handlers.DeleteStudent).Methods("DELETE")

	// Subject routes
	apiRouter.HandleFunc("/subjects", handlers.GetAllSubjects).Methods("GET")
	apiRouter.HandleFunc("/subjects/{id}", handlers.GetSubject).Methods("GET")
	apiRouter.HandleFunc("/subjects/{subjectID}/students", handlers.GetStudentsBySubjectID).Methods("GET")
	apiRouter.HandleFunc("/subjects/{subjectID}/terms", handlers.GetTermsBySubjectID).Methods("GET")
	apiRouter.HandleFunc("/subjects/{subjectID}/terms/{termID}/grade-labels", handlers.GetGradeLabelsForSubject).Methods("GET")
	apiRouter.HandleFunc("/subjects", handlers.CreateSubject).Methods("POST")
	apiRouter.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}/terms/{termID}", handlers.AssignGradeLabelToSubjectByTerm).Methods("POST")
	apiRouter.HandleFunc("/subjects/{id}", handlers.UpdateSubject).Methods("PUT")
	apiRouter.HandleFunc("/subjects/{id}", handlers.DeleteSubject).Methods("DELETE")
	apiRouter.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}/terms/{termID}", handlers.RemoveGradeLabelFromSubjectByTerm).Methods("DELETE")

	// Grade Labels routes
	apiRouter.HandleFunc("/grade-labels", handlers.CreateGradeLabel).Methods("POST")
	apiRouter.HandleFunc("/grade-labels", handlers.GetAllGradeLabels).Methods("GET")
	apiRouter.HandleFunc("/grade-labels/{id}", handlers.GetGradeLabel).Methods("GET")
	apiRouter.HandleFunc("/grade-labels/{id}", handlers.UpdateGradeLabel).Methods("PUT")
	apiRouter.HandleFunc("/grade-labels/{id}", handlers.DeleteGradeLabel).Methods("DELETE")

	// Term routes
	apiRouter.HandleFunc("/terms", handlers.GetAllTerms).Methods("GET")
	apiRouter.HandleFunc("/terms/{id}", handlers.GetTerm).Methods("GET")
	apiRouter.HandleFunc("/terms", handlers.CreateTerm).Methods("POST")
	apiRouter.HandleFunc("/terms/{id}", handlers.UpdateTerm).Methods("PUT")
	apiRouter.HandleFunc("/terms/{id}", handlers.DeleteTerm).Methods("DELETE")

	// Serve static files from the "static" directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the HTTP server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", enableCors(router)))
	//log.Fatal(http.ListenAndServe(":8080", router))   //In case, we don't use CORS
}

/*
///////////////PROJECT STRUCTURE//////////////////////

classroom-management/
│
├── main.go
│
├── handlers/
│   ├── classroom_handlers.go
│   ├── grade_handlers.go
│   ├── grade_label_handlers.go
│   ├── student_handlers.go
│   ├── subject_handlers.go
│   ├── teacher_handlers.go
│   └── term_handlers.go
│
├── models/
│   ├── classroom.go
│   ├── grade_label.go
│   ├── grade.go
│   ├── student.go
│   ├── subject.go
│   ├── teacher.go
│   └── term.go
│
├── database/
│   ├── database.go
│   ├── classroom.go
│   ├── grade_label.go
│   ├── grade.go
│   ├── student.go
│   ├── subject.go
│   ├── teacher.go
│   └── term.go
│
├── middleware/
│   └── auth.go
│
├── classroom_management.sql
├── go.mod
├── go.sum
└── static/
    ├── index.html
    ├── script.js
    ├── styles.css
	├── main.html
    ├── main.js
    ├── main.css
	├── signup.html
    ├── classroom-grades-display.html
    ├── classroom-grades-display.js
    ├── classroom-grades-display.css
    ├── classroom-grades-upload.html
    ├── classroom-grades-upload.js
    └── classroom-grades-upload.css

///////////Postgres Database//////////
psql

\l

CREATE DATABASE classroom_management;

DROP DATABASE classroom_management;     //for deleting a database

\c classroom_management

pwd

\i /Users/tavito/Documents/go/classroom-management/classroom_management.sql

\dt

To show all the elements in terms:

SELECT * FROM terms;

To delete tables:

DROP TABLE grade_labels;

To create table:

-- Table to store grade labels for each classroom and subject
CREATE TABLE grade_labels (
    id SERIAL PRIMARY KEY,
    label VARCHAR(255), -- Label for the grade (e.g., "1st input", "2nd input", "lesson", "quiz", etc.)
    date DATE,          -- New field for date
    skill VARCHAR(255), -- New field for skill
    teacher_id INT REFERENCES teachers(id)
);

////////////////CURL COMMANDS///////////////////
Sign Up:

curl -X POST http://localhost:8080/signup \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Gustavo",
        "email": "jimgustavo1987@gmail.com",
        "password": "hola"
    }'

curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{
        "email": "jimgustavo1987@gmail.com",
        "password": "hola"
    }'

	returned token: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjoxLCJleHAiOjE3MTgzODE5NjN9.J7v5VPJgaRgaVCgZqOL4KG9aUHe8RvGqG5JLM7dHSCc"}



Create a Classroom:
curl -X POST http://localhost:8080/api/classrooms \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjoxLCJleHAiOjE3MTgzODE5NjN9.J7v5VPJgaRgaVCgZqOL4KG9aUHe8RvGqG5JLM7dHSCc" \
    -d '{
        "name": "Room A",
        "teacher_id": 1
    }'

Get all Classrooms:

curl -X GET http://localhost:8080/api/classrooms \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjoxLCJleHAiOjE3MTgzODE5NjN9.J7v5VPJgaRgaVCgZqOL4KG9aUHe8RvGqG5JLM7dHSCc"


Create a Student:

curl -X POST http://localhost:8080/api/students \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjoxLCJleHAiOjE3MTg0MDUyOTZ9.l0dx49smDnHw9hIVHJpBfygumPUp9hJLKI2fRDjNagU" \
    -d '{
        "name": "Isabelita",
        "classroom_id": 1,
        "teacher_id": 1
    }'
Get all students:

curl -X GET http://localhost:8080/api/students \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjoxLCJleHAiOjE3MTg0MDUyOTZ9.l0dx49smDnHw9hIVHJpBfygumPUp9hJLKI2fRDjNagU"


Update a Student:
curl -X PUT -H "Content-Type: application/json" -d '{
    "name": "Enrique Ruiz",
    "classroom_id":1
}' http://localhost:8080/students/16
Delete a Student:
curl -X DELETE http://localhost:8080/students/11

Create a Subject:
curl -X POST -H "Content-Type: application/json" -d '{"name":"Mathematics"}' http://localhost:8080/subjects

Assign a Subject to a Classroom:
curl -X POST http://localhost:8080/classrooms/1/subject/7

Upload students list from and xlsx file:
curl -X POST "http://localhost:8080/classrooms/4/upload-students?startCell=B7&endCell=B21&sheetName=Matematicas" \
     -F "file=@./consolidado.xlsx"


curl -X DELETE http://localhost:8080/teachers/1

curl -X POST http://localhost:8080/classrooms/1/subject/2
*/
