package models

// Subject represents a subject entity
type Subject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SubjectWithGradeLabels represents a subject along with its associated grade labels
type SubjectWithGradeLabels struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	GradeLabels []string `json:"grade_labels"`
}

type AddSubjectRequest struct {
	GradeLabelIDs []int `json:"gradeLabelIDs"`
}
