// models/grade_label_classroom.go

package models

// GradeLabelSubject represents a many-to-many relationship between grade labels, subjects, and classrooms
type GradeLabelSubject struct {
	SubjectID    int `json:"subject_id"`
	GradeLabelID int `json:"grade_label_id"`
}
