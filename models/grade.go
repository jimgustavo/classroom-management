package models

type Grade struct {
	LabelID   int     `json:"label_id"`
	Grade     float32 `json:"grade"`
	TeacherID int     `json:"teacher_id"`
}

type TermGrades struct {
	Term   string  `json:"term"`
	Grades []Grade `json:"grades"`
}

type StudentTermGrades struct {
	StudentID int          `json:"student_id"`
	SubjectID int          `json:"subject_id"`
	Terms     []TermGrades `json:"terms"`
}

type GradesData struct {
	Grades []StudentTermGrades `json:"grades"`
}

type Response struct {
	Message string `json:"message"`
}
