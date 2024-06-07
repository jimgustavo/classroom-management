package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/handlers"
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
	//defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Classroom routes
	router.HandleFunc("/classrooms", handlers.GetAllClassrooms).Methods("GET")
	router.HandleFunc("/classrooms/{id}", handlers.GetClassroom).Methods("GET")
	router.HandleFunc("/classrooms/{classroomID}/subjects", handlers.GetSubjectsInClassroom).Methods("GET")
	router.HandleFunc("/classrooms/{classroom_id}/students", handlers.GetStudentsByClassroom).Methods("GET")
	router.HandleFunc("/classrooms/{classroomID}/grades/get", handlers.GetGradesByClassroomID).Methods("GET")
	router.HandleFunc("/classrooms/{classroomID}/terms/{termID}/grades", handlers.GetGradesByClassroomIDAndTermID).Methods("GET")
	router.HandleFunc("/classrooms", handlers.CreateClassroom).Methods("POST")
	router.HandleFunc("/classrooms/{classroomID}/subject/{subjectID}", handlers.AddSubjectToClassroom).Methods("POST")
	router.HandleFunc("/classrooms/{classroomID}/grades", handlers.UploadGradesToClassroom).Methods("POST")
	router.HandleFunc("/classrooms/{classroomID}/upload-students", handlers.UploadStudentsFromExcel).Methods("POST") // for uploading students from an Excel file
	router.HandleFunc("/classrooms/{id}", handlers.UpdateClassroom).Methods("PUT")
	router.HandleFunc("/classrooms/{id}", handlers.DeleteClassroom).Methods("DELETE")
	router.HandleFunc("/classrooms/{classroomID}/students/{studentID}", handlers.UnrollStudentFromClassroom).Methods("DELETE")
	router.HandleFunc("/classrooms/{classroomID}/subjects/{subjectID}", handlers.RemoveSubjectFromClassroom).Methods("DELETE")

	// Student routes
	router.HandleFunc("/students", handlers.GetAllStudents).Methods("GET")
	router.HandleFunc("/students/with-classroom-and-subjects", handlers.GetAllStudentsWithClassroomAndSubjects).Methods("GET")
	router.HandleFunc("/students/{id}", handlers.GetStudent).Methods("GET")
	router.HandleFunc("/students/{studentID}/subjects", handlers.GetSubjectsByStudentID).Methods("GET")
	router.HandleFunc("/students", handlers.CreateStudent).Methods("POST")
	router.HandleFunc("/students/{id}", handlers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", handlers.DeleteStudent).Methods("DELETE")

	// Subject routes
	router.HandleFunc("/subjects", handlers.GetAllSubjects).Methods("GET")
	router.HandleFunc("/subjects/{id}", handlers.GetSubject).Methods("GET")
	router.HandleFunc("/subjects/{subjectID}/students", handlers.GetStudentsBySubjectID).Methods("GET")
	router.HandleFunc("/subjects/{subjectID}/terms", handlers.GetTermsBySubjectID).Methods("GET")
	router.HandleFunc("/subjects/{subjectID}/terms/{termID}/grade-labels", handlers.GetGradeLabelsForSubject).Methods("GET")
	router.HandleFunc("/subjects", handlers.CreateSubject).Methods("POST")
	router.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}/terms/{termID}", handlers.AssignGradeLabelToSubjectByTerm).Methods("POST")
	router.HandleFunc("/subjects/{id}", handlers.UpdateSubject).Methods("PUT")
	router.HandleFunc("/subjects/{id}", handlers.DeleteSubject).Methods("DELETE")
	router.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}/terms/{termID}", handlers.RemoveGradeLabelFromSubjectByTerm).Methods("DELETE")

	// Grade Labels routes
	router.HandleFunc("/grade-labels", handlers.CreateGradeLabel).Methods("POST")
	router.HandleFunc("/grade-labels", handlers.GetAllGradeLabels).Methods("GET")
	router.HandleFunc("/grade-labels/{id}", handlers.GetGradeLabel).Methods("GET")
	router.HandleFunc("/grade-labels/{id}", handlers.UpdateGradeLabel).Methods("PUT")
	router.HandleFunc("/grade-labels/{id}", handlers.DeleteGradeLabel).Methods("DELETE")

	// Term routes
	router.HandleFunc("/terms", handlers.GetAllTerms).Methods("GET")
	router.HandleFunc("/terms/{id}", handlers.GetTerm).Methods("GET")
	router.HandleFunc("/terms", handlers.CreateTerm).Methods("POST")
	router.HandleFunc("/terms/{id}", handlers.UpdateTerm).Methods("PUT")
	router.HandleFunc("/terms/{id}", handlers.DeleteTerm).Methods("DELETE")

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
│ 	├── student_handlers.go
│ 	├── subject_handlers.go
│	  └── term_handlers.go
│
├── models/
│   ├── classroom.go
│   ├── grade_label.go
│   ├── grade.go
│	  ├── student.go
│	  ├── subject.go
│ 	└── term.go
│
└── database/
│   ├── database.go
│   ├── classroom.go
│   ├── grade_label.go
│   ├── grade.go
│	  ├── student.go
│	  ├── subject.go
│  	└── term.go
│
│── classroom_management.sql
│── go.mod
│── go.sum
└── static/
    ├── index.html
    ├── script.js
    ├── styles.css
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

DROP TABLE grades;

To create table:

CREATE TABLE grades (
    student_id INT,
    subject_id INT,
    term VARCHAR(50),
    label VARCHAR(50),
    grade FLOAT,
    classroom_id INT,
    PRIMARY KEY (student_id, subject_id, term, label),
    CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES students(id),
    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(id),
    CONSTRAINT fk_classroom FOREIGN KEY (classroom_id) REFERENCES classrooms(id)
);

////////////////CURL COMMANDS///////////////////
Create a Classroom:
curl -X POST -H "Content-Type: application/json" -d '{"name":"Room A"}' http://localhost:8080/classrooms

Create a Student:
curl -X POST -H "Content-Type: application/json" -d '{"name":"Jimmy Ruiz", "classroom_id":1}' http://localhost:8080/students

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


curl -X DELETE http://localhost:8080/subjects/1/grade-labels/5/terms/1

curl -X POST http://localhost:8080/classrooms/1/subject/2
*/
