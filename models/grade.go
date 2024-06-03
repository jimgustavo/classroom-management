package models

type Grade struct {
	Label string `json:"label"`
	Grade string `json:"grade"`
}

type StudentGrade struct {
	StudentID int     `json:"studentID"`
	SubjectID int     `json:"subjectID"`
	Grades    []Grade `json:"grades"`
}

// StudentGrades structure
type StudentGrades struct {
	StudentID int     `json:"studentID"`
	SubjectID int     `json:"subjectID"`
	Grades    []Grade `json:"grades"`
}

type GradesData struct {
	Grades []StudentGrade `json:"grades"`
}

// Response structure
type Response struct {
	Message string `json:"message"`
}

// GradeInfo represents information about a grade, including the grade value, label, subject, and classroom.
type GradeInfo struct {
	Grade     float64 `json:"grade"`
	Label     string  `json:"label"`
	Subject   string  `json:"subject"`
	Classroom string  `json:"classroom"`
}

// StudentGradeInfo represents information about a student's grade, label, subject, and classroom.
type StudentGradeInfo struct {
	StudentID   int         `json:"student_id"`
	StudentName string      `json:"student_name"`
	Grades      []GradeInfo `json:"grades"`
}

type ClassroomGrades struct {
	ClassroomID int                `json:"classroom_id"`
	Grades      []StudentGradeInfo `json:"grades"`
}
