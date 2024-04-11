package models

// GradeInput represents a grade input entity
type GradeLabel struct {
	ID        int    `json:"id"`
	SubjectID int    `json:"subject_id"`
	Label     string `json:"label"`
}
