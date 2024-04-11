package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/handlers"
)

func main() {
	// Initialize the database connection
	err := database.InitializeDB("postgres://tavito:mamacita@localhost:5432/classroom_management?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	router.HandleFunc("/students-with-grades", handlers.GetAllStudentsWithGrades).Methods("GET")
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}/grades", handlers.AddGrade).Methods("POST")
	// Route to retrieve grade information by student ID
	router.HandleFunc("/students/{studentID}/grades", handlers.GetGradeByStudentID).Methods("GET")

	// Route to add subject with its grade labels to a classroom and retrive all the subjects with their grade labels in a classroom
	router.HandleFunc("/classrooms/{classroomID}/subject/{subjectID}", handlers.AddSubjectToClassroom).Methods("POST")
	router.HandleFunc("/classrooms/{classroomID}/subjects", handlers.GetSubjectsInClassroom).Methods("GET")

	// Referenced Classroom - Subject - Student routes
	router.HandleFunc("/classrooms/{classroomID}/subjects", handlers.GetSubjectsByClassroomID).Methods("GET")
	router.HandleFunc("/classrooms/{classroom_id}/students", handlers.GetStudentsByClassroom).Methods("GET")
	router.HandleFunc("/classrooms/{classroomID}/subjects/{subjectID}", handlers.DeleteSubjectFromStudents).Methods("DELETE")
	router.HandleFunc("/classrooms/{classroomID}/subjects", handlers.DeleteSubjectsByClassroomID).Methods("DELETE")
	router.HandleFunc("/classrooms/{classroomID}/subjects/{subjectID}", handlers.AddSubjectToStudents).Methods("POST")
	router.HandleFunc("/subjects/{subjectID}/students", handlers.GetStudentsBySubjectID).Methods("GET")
	router.HandleFunc("/students/{studentID}/subjects", handlers.GetSubjectsByStudentID).Methods("GET")
	router.HandleFunc("/students/with-classroom-and-subjects", handlers.GetAllStudentsWithClassroomAndSubjects).Methods("GET")

	// Classroom routes
	router.HandleFunc("/classrooms", handlers.CreateClassroom).Methods("POST")
	router.HandleFunc("/classrooms", handlers.GetAllClassrooms).Methods("GET")
	router.HandleFunc("/classrooms/{id}", handlers.GetClassroom).Methods("GET")
	router.HandleFunc("/classrooms/{id}", handlers.UpdateClassroom).Methods("PUT")
	router.HandleFunc("/classrooms/{id}", handlers.DeleteClassroom).Methods("DELETE")

	// Student routes
	router.HandleFunc("/students", handlers.CreateStudent).Methods("POST")
	router.HandleFunc("/students", handlers.GetAllStudents).Methods("GET")
	router.HandleFunc("/students/{id}", handlers.GetStudent).Methods("GET")
	router.HandleFunc("/students/{id}", handlers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", handlers.DeleteStudent).Methods("DELETE")

	// Subject routes
	router.HandleFunc("/subjects", handlers.CreateSubject).Methods("POST")
	router.HandleFunc("/subjects", handlers.GetAllSubjects).Methods("GET")
	router.HandleFunc("/subjects/{id}", handlers.GetSubject).Methods("GET")
	router.HandleFunc("/subjects/{id}", handlers.UpdateSubject).Methods("PUT")
	router.HandleFunc("/subjects/{id}", handlers.DeleteSubject).Methods("DELETE")

	// Grade Labels routes
	router.HandleFunc("/grade-labels", handlers.CreateGradeLabel).Methods("POST")
	router.HandleFunc("/grade-labels", handlers.GetAllGradeLabels).Methods("GET")
	router.HandleFunc("/grade-labels/{id}", handlers.GetGradeLabel).Methods("GET")
	router.HandleFunc("/grade-labels/{id}", handlers.UpdateGradeLabel).Methods("PUT")
	router.HandleFunc("/grade-labels/{id}", handlers.DeleteGradeLabel).Methods("DELETE")

	// Subject - Grade Label routes
	router.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}", handlers.AssignGradeLabelToSubject).Methods("POST")
	router.HandleFunc("/subjects/{subjectID}/grade-labels", handlers.GetGradeLabelsForSubject).Methods("GET")

	// Start the HTTP server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
///////////////STRUCTURE//////////////////////

classroom-management/
│
├── main.go
│
├── handlers/
│   ├── classroom_handlers.go
│   ├── student_handlers.go
│   ├── subject_handlers.go
│   ├── student_subjects_handlers.go
│	├── grade_label_handlers.go
│	└── grade_handlers.go
│
├── models/
│   ├── classroom.go
│   ├── student.go
│   ├── subject.go
│   ├── student_subjects.go
│	├── grade_label.go
│	└── grade.go
│
└── database/
    └── database.go

///////////Postgres Database//////////
psql

\l

CREATE DATABASE classroom_management;

DROP DATABASE classroom_management;     //for deleting a database

\c classroom_management

pwd

\i /Users/tavito/Documents/go/classroom-management/classroom_management.sql

\dt

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
curl -X POST http://localhost:8080/classrooms/2/subjects/1

Retrieve all students with classroom and subjects assigned:
curl -X GET http://localhost:8080/students/with-classroom-and-subjects

Create a grade label:
curl -X POST -H "Content-Type: application/json" -d '{"subject_id": 1, "label": "Final Term Quiz"}' http://localhost:8080/grade-labels

Update a grade label:
curl -X PUT -H "Content-Type: application/json" -d '{"subject_id": 1, "label": "First Input"}' http://localhost:8080/grade-labels/1

curl -X POST http://localhost:8080/classrooms/1/subject/2

curl -X POST -H "Content-Type: application/json" -d '{"gradeLabelIDs":[11,12,13,14]}' http://localhost:8080/classrooms/4/subject/2
curl -X POST -H "Content-Type: application/json" -d '{"gradeLabelIDs":[]}' http://localhost:8080/classrooms/2/subject/1

curl -X POST -H "Content-Type: application/json" -d '{"label_id":1, "grade":85.5}' http://localhost:8080/students/1/subjects/1/grades
curl http://localhost:8080/students/1/grades


*/
