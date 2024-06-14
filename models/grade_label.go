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
