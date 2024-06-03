// models/grade_label.go

package models

// GradeInput represents a grade input entity
type GradeLabel struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// GradeLabelSubject represents a many-to-many relationship between grade labels, subjects, and classrooms
type GradeLabelSubject struct {
	SubjectID    int `json:"subject_id"`
	GradeLabelID int `json:"grade_label_id"`
}
