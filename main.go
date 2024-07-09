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

	// Initialize router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/logout", handlers.Logout)
	// Grades routes
	router.HandleFunc("/classroom/{classroomID}/grades/get", handlers.GetGradesByClassroomID).Methods("GET")
	router.HandleFunc("/classroom/{classroomID}/terms/{termID}/grades", handlers.GetGradesByClassroomIDAndTermID).Methods("GET")
	// Reinforcement grades routes
	router.HandleFunc("/grade-labels/reinforcement", handlers.GetAllReinforcementGradeLabels).Methods("GET")
	router.HandleFunc("/grade-labels/reinforcement/teacher/{teacherID}", handlers.GetReinforcementGradeLabelsByTeacher).Methods("GET")
	router.HandleFunc("/grade-labels/reinforcement/classroom/{classroomID}/term/{termID}", handlers.GetReinforcementGradeLabelsByClassroomAndTerm).Methods("GET")
	// Averages routes
	router.HandleFunc("/classroom/{classroomID}/averages", handlers.GetAverageGradesByClassroomID).Methods("GET")
	router.HandleFunc("/classroom/{classroomID}/averageswithfactors", handlers.GetAveragesWithFactorsByClassroomID).Methods("GET")
	router.HandleFunc("/classroom/{classroomID}/averageswithreinforcement", handlers.GetAveragesWithReinforcementByClassroomID).Methods("GET")
	router.HandleFunc("/classroom/{classroomID}/averageswithfactors_trimesters", handlers.GetAveragesWithFactorsByClassroomIDForTrimesters).Methods("GET")
	// XLSX REPORTS routes
	router.HandleFunc("/xlsx-report/teachers/{teacherID}/classrooms/{classroomID}/academicPeriod/{academicPeriodID}/terms/{termID}", handlers.GenerateTeacherGradesReport).Methods("GET")
	router.HandleFunc("/xlsx-average/teachers/{teacherID}/classrooms/{classroomID}/academicPeriod/{academicPeriodID}", handlers.GenerateFinalAveragesReport).Methods("GET")
	// TeacherData routes
	router.HandleFunc("/teacherdata", handlers.CreateOrUpdateTeacherDataHandler).Methods("POST")
	router.HandleFunc("/teacherdata/{teacherID}", handlers.GetTeacherDataByTeacherIDHandler).Methods("GET")
	router.HandleFunc("/teacherdata/{id}", handlers.GetTeacherDataByIDHandler).Methods("GET")
	router.HandleFunc("/teacherdata/{id}", handlers.UpdateTeacherDataHandler).Methods("PUT")
	router.HandleFunc("/teacherdata/{id}", handlers.DeleteTeacherDataHandler).Methods("DELETE")
	// Admin routes
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.Authenticate)
	adminRouter.Use(middleware.AdminOnly)
	adminRouter.HandleFunc("/classrooms", handlers.GetAllClassrooms).Methods("GET")
	adminRouter.HandleFunc("/teachers", handlers.GetAllTeachersHandler).Methods("GET")
	adminRouter.HandleFunc("/teacher/{id}", handlers.DeleteTeacherHandler).Methods("DELETE")
	adminRouter.HandleFunc("/teacherdata", handlers.GetAllTeacherDataHandler).Methods("GET")
	adminRouter.HandleFunc("/teacher/{id}/role/{role}", handlers.UpdateTeacherRoleHandler).Methods("PUT")
	// Academic Period routes
	adminRouter.HandleFunc("/academic_periods", handlers.CreateAcademicPeriodHandler).Methods("POST")
	adminRouter.HandleFunc("/academic_periods", handlers.GetAllAcademicPeriodsHandler).Methods("GET")
	adminRouter.HandleFunc("/academic_periods/{id}", handlers.GetAcademicPeriodByIDHandler).Methods("GET")
	adminRouter.HandleFunc("/academic_periods/{id}", handlers.UpdateAcademicPeriodHandler).Methods("PUT")
	adminRouter.HandleFunc("/academic_periods/{id}", handlers.DeleteAcademicPeriodHandler).Methods("DELETE")
	adminRouter.HandleFunc("/academic_periods/{academicPeriodID}/terms/{termID}", handlers.AssignTermToAcademicPeriodHandler).Methods("POST")
	adminRouter.HandleFunc("/academic_periods/{id}/terms", handlers.GetTermsByAcademicPeriod).Methods("GET")
	// Term routes
	adminRouter.HandleFunc("/terms", handlers.GetAllTerms).Methods("GET")
	adminRouter.HandleFunc("/terms", handlers.CreateTerm).Methods("POST")
	adminRouter.HandleFunc("/terms/{id}", handlers.UpdateTerm).Methods("PUT")
	adminRouter.HandleFunc("/terms/{id}", handlers.DeleteTerm).Methods("DELETE")
	//
	adminRouter.HandleFunc("/students", handlers.GetAllStudents).Methods("GET")
	adminRouter.HandleFunc("/subjects", handlers.GetAllSubjects).Methods("GET")
	adminRouter.HandleFunc("/grade-labels", handlers.GetAllGradeLabels).Methods("GET")

	// ProTeacher routes
	proTeacherRouter := router.PathPrefix("/proteacher").Subrouter()
	proTeacherRouter.Use(middleware.Authenticate)
	proTeacherRouter.Use(middleware.ProTeacherOnly)
	// PDF REPORT
	proTeacherRouter.HandleFunc("/pdfminute/teacher/{teacherID}/classroom/{classroomID}/student/{studentID}", handlers.GenerateReportHandler).Methods("GET")

	// Regular teacher routes (protected)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.Authenticate)
	// Classroom routes
	apiRouter.HandleFunc("/classrooms/teacher/{teacherID}", handlers.GetClassroomsByTeacherID).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{id}", handlers.GetClassroom).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/subjects", handlers.GetSubjectsInClassroom).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroom_id}/students", handlers.GetStudentsByClassroom).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/grades/get", handlers.GetGradesByClassroomID).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/terms/{termID}/grades", handlers.GetGradesByClassroomIDAndTermID).Methods("GET")
	apiRouter.HandleFunc("/classrooms/{classroomID}/averages", handlers.GetAverageGradesByClassroomID).Methods("GET")
	apiRouter.HandleFunc("/classrooms", handlers.CreateClassroom).Methods("POST")
	apiRouter.HandleFunc("/classrooms/{classroomID}/subject/{subjectID}", handlers.AddSubjectToClassroom).Methods("POST")
	apiRouter.HandleFunc("/classrooms/{classroomID}/grades", handlers.UploadGradesToClassroom).Methods("POST")
	apiRouter.HandleFunc("/classrooms/{classroomID}/upload-students/{teacherID}", handlers.UploadStudentsFromExcel).Methods("POST") // for uploading students from an Excel file
	apiRouter.HandleFunc("/classrooms/{id}", handlers.UpdateClassroom).Methods("PUT")
	apiRouter.HandleFunc("/classrooms/{id}", handlers.DeleteClassroom).Methods("DELETE")
	apiRouter.HandleFunc("/classrooms/{classroomID}/students/{studentID}", handlers.UnrollStudentFromClassroom).Methods("DELETE")
	apiRouter.HandleFunc("/classrooms/{classroomID}/subjects/{subjectID}", handlers.RemoveSubjectFromClassroom).Methods("DELETE")
	// Student routes
	apiRouter.HandleFunc("/students/teacher/{teacherID}", handlers.GetStudentsByTeacherID).Methods("GET")
	apiRouter.HandleFunc("/students/with-classroom-and-subjects", handlers.GetAllStudentsWithClassroomAndSubjects).Methods("GET")
	apiRouter.HandleFunc("/students/{id}", handlers.GetStudent).Methods("GET")
	apiRouter.HandleFunc("/students/{studentID}/subjects", handlers.GetSubjectsByStudentID).Methods("GET")
	apiRouter.HandleFunc("/students", handlers.CreateStudent).Methods("POST")
	apiRouter.HandleFunc("/students/{id}", handlers.UpdateStudent).Methods("PUT")
	apiRouter.HandleFunc("/students/{id}", handlers.DeleteStudent).Methods("DELETE")
	// Subject routes
	apiRouter.HandleFunc("/subjects/teacher/{teacherID}", handlers.GetSubjectsByTeacherID).Methods("GET")
	apiRouter.HandleFunc("/subjects/{subjectID}/students", handlers.GetStudentsBySubjectID).Methods("GET")
	apiRouter.HandleFunc("/subjects/{subjectID}/terms/{termID}/grade-labels", handlers.GetGradeLabelsForSubject).Methods("GET")
	apiRouter.HandleFunc("/subjects", handlers.CreateSubject).Methods("POST")
	apiRouter.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}/terms/{termID}", handlers.AssignGradeLabelToSubjectByTerm).Methods("POST")
	apiRouter.HandleFunc("/subjects/{id}", handlers.UpdateSubject).Methods("PUT")
	apiRouter.HandleFunc("/subjects/{id}", handlers.DeleteSubject).Methods("DELETE")
	apiRouter.HandleFunc("/subjects/{subjectID}/grade-labels/{gradeLabelID}/terms/{termID}", handlers.RemoveGradeLabelFromSubjectByTerm).Methods("DELETE")
	// Grade Labels routes
	apiRouter.HandleFunc("/grade-labels/teacher/{teacherID}", handlers.GetGradeLabelsByTeacherID).Methods("GET")
	apiRouter.HandleFunc("/grade-labels", handlers.CreateGradeLabel).Methods("POST")
	apiRouter.HandleFunc("/grade-labels/{id}", handlers.UpdateGradeLabel).Methods("PUT")
	apiRouter.HandleFunc("/grade-labels/{id}", handlers.DeleteGradeLabel).Methods("DELETE")
	// Reinforcement routes
	apiRouter.HandleFunc("/grade-labels/reinforcement", handlers.CreateReinforcementGradeLabel).Methods("POST")
	apiRouter.HandleFunc("/grade-labels/reinforcement", handlers.GetAllReinforcementGradeLabels).Methods("GET")
	apiRouter.HandleFunc("/grade-labels/reinforcement/teacher/{teacherID}", handlers.GetReinforcementGradeLabelsByTeacher).Methods("GET")
	apiRouter.HandleFunc("/grade-labels/reinforcement/classroom/{classroomID}/term/{termID}", handlers.GetReinforcementGradeLabelsByClassroomAndTerm).Methods("GET")
	apiRouter.HandleFunc("/grade-labels/reinforcement/{id}", handlers.DeleteReinforcementGradeLabel).Methods("DELETE")
	// Academic Period routes
	apiRouter.HandleFunc("/academic_periods", handlers.GetAllAcademicPeriodsHandler).Methods("GET")
	// Term routes
	apiRouter.HandleFunc("/academic_periods/{id}/terms", handlers.GetTermsByAcademicPeriod).Methods("GET")

	// Serve static files from the "static" directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the HTTP server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", enableCors(router)))
}

/*

	ROUTES NOT USED FOR NOW:
	apiRouter.HandleFunc("/subjects/{id}", handlers.GetSubject).Methods("GET")
	apiRouter.HandleFunc("/grade-labels/{id}", handlers.GetGradeLabel).Methods("GET")
	apiRouter.HandleFunc("/terms/{id}", handlers.GetTerm).Methods("GET")


///////////Postgres Database//////////
psql

\l

CREATE DATABASE classroom_management;

DROP DATABASE classroom_management;     //for deleting a database

\c classroom_management;

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

ENABLE pgcrypto:
-- Enable the pgcrypto extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Verify the extension is enabled
\dx



SIGNUP ADMIN USERS:
-- Connect to your PostgreSQL database
\c classroom_management;

-- Insert a new admin user
INSERT INTO teachers (name, email, password, role)
VALUES ('tavito', 'jimgustavo@icloud.com', crypt('mamacita', gen_salt('bf')), 'admin');

DELETE ADMIN USERS:
-- Connect to your PostgreSQL database
\c classroom_management;

-- Delete an admin user by email
DELETE FROM teachers WHERE email = 'jimgustavo@icloud.com' AND role = 'admin';


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
        "email": "jimgustavo@icloud.com",
        "password": "mamacita"
    }'

	returned token: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo0LCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MTk5NjkzMzN9.HRAF33FNjNju6gK6V_mkuNR79kfA9SbMwP3FThFz5yA"}


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
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo2LCJleHAiOjE3MTg1MDY4MTd9.aYaUjQN3j_-5xEgzaLntBxDOK1ZkEr_xpmj8HTI-Kxw"

curl -X GET http://localhost:8080/api/classrooms/teacher/4 \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo2LCJleHAiOjE3MTg1MDY4MTd9.aYaUjQN3j_-5xEgzaLntBxDOK1ZkEr_xpmj8HTI-Kxw"


Get all teachers:

curl -X GET http://localhost:8080/admin/teachers \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo1LCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MTk5NzY4MTV9.yKqCkntDL_M55EJX5kgWfF-UJY0rQIPnNtnABAE_Vk4"

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

curl -X GET http://localhost:8080/api/students/ \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo0LCJleHAiOjE3MTg1MTE3OTh9.RLvS-QfxBuhMHmNaEaUwMPh-Qu5b76x-tK4vGBXqvLM"

curl -X GET http://localhost:8080/api/students/teacher/4 \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo0LCJleHAiOjE3MTg1MTE3OTh9.RLvS-QfxBuhMHmNaEaUwMPh-Qu5b76x-tK4vGBXqvLM"


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

Teacher Data:

curl -X POST http://localhost:8080/teacherdata \
-H "Content-Type: application/json" \
-d '{
    "school": "Sample School",
    "school_year": "2023-2024",
    "school_hours": "8:00 AM - 3:00 PM",
    "country": "Ecuador",
    "city": "Zamora",
    "teacher_id": 1,
    "teacher_full_name": "John Doe",
	"teacher_birthday": "1987-10-17",
    "id_number": "123456789",
    "labor_dependency_relationship": "Sample Relationship",
	"institutional_email": "jimmy.ruiz@educacion.gob.ec",
	"phone": "0961763152",
    "principal": "Principal Name",
    "vice_principal": "Vice Principal Name",
    "dece": "DECE Name",
    "inspector": "Inspector Name"
}'

curl -X PUT http://localhost:8080/admin/teacher/1/role/proteacher \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFjaGVyX2lkIjo1LCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MjAwMjgyMzV9.1BJSr31K5Oqyq3_0gXbAyTA6YCu7HYXD38_-6sIf4Nw"

*/
