// models/grade_label.go

package models

// GradeLabel represents a grade label entity
type GradeLabel struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Date      string `json:"date"`
	Skill     string `json:"skill"`
	TeacherID int    `json:"teacher_id"`
}

type GradeLabelTerm struct {
	ID     int    `json:"id"`
	Label  string `json:"label"`
	TermID int    `json:"term_id"`
}

type ReinforcementGradeLabel struct {
	ID          int     `json:"id"`
	StudentID   int     `json:"student_id"`
	ClassroomID int     `json:"classroom_id"`
	SubjectID   int     `json:"subject_id"`
	TermID      int     `json:"term_id"`
	Label       string  `json:"label"`
	Date        string  `json:"date"`
	Skill       string  `json:"skill"`
	TeacherID   int     `json:"teacher_id"`
	Grade       float64 `json:"grade"`
}
