// models/subject.go
package models

// Subject represents a subject with its ID, name, and associated teacher ID.
type Subject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TeacherID int    `json:"teacher_id"`
}

type SubjectWithGradeLabels struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	GradeLabels []GradeLabelTerm `json:"grade_labels"`
}

type AddSubjectRequest struct {
	GradeLabelIDs []int `json:"gradeLabelIDs"`
}
