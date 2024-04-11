package models

// StudentSubject represents the many-to-many relationship between students and subjects
type StudentSubject struct {
	StudentID int `json:"student_id"`
	SubjectID int `json:"subject_id"`
}
