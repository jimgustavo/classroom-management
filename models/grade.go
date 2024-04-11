package models

// Grade represents a grade entry for a student.
type Grade struct {
	StudentID   int     `json:"student_id"`
	ClassroomID int     `json:"classroom_id"`
	SubjectID   int     `json:"subject_id"`
	LabelID     int     `json:"label_id"`
	Grade       float64 `json:"grade"`
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
