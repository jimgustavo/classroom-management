// models/grade_label.go

package models

// GradeInput represents a grade input entity
type GradeLabel struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// GradeLabelWithID represents a grade label with its ID and name
type GradeLabelWithID struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type GradeLabelSubject struct {
	SubjectID    int `json:"subject_id"`
	GradeLabelID int `json:"grade_label_id"`
	TermID       int `json:"term_id"`
}

type GradeLabelTerm struct {
	ID     int    `json:"id"`
	Label  string `json:"label"`
	TermID int    `json:"term_id"`
}

/*
// GradeLabelTermPair represents a pair of grade label ID and term ID
type GradeLabelTermPair struct {
	GradeLabelID int `json:"grade_label_id"`
	TermID       int `json:"term_id"`
}
*/
