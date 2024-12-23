// models/grade.go
package models

// Grade represents a grade with associated label ID, grade value, and teacher ID.
type Grade struct {
	LabelID   int     `json:"label_id"`
	Grade     float32 `json:"grade"`
	TeacherID int     `json:"teacher_id"`
}

// GradeSkill represents a grade with associated skill and grade value.
type GradeSkill struct {
	Skill string  `json:"skill"`
	Grade float32 `json:"grade"`
	Date  string  `json:"date"`
}

// StudentTermGradeSkills represents the grades with associated skills for a student in a particular subject.
type StudentTermGradeSkills struct {
	StudentID   int               `json:"student_id"`
	SubjectID   int               `json:"subject_id"`
	SubjectName string            `json:"subject_name"` // New field
	Terms       []TermGradeSkills `json:"terms"`
}

// TermGradeSkills represents the grades for a particular term with associated skills.
type TermGradeSkills struct {
	Term   string       `json:"term"`
	Grades []GradeSkill `json:"grades"`
}

// GradesDataSkills represents the overall data structure for grades including skills.
type GradesDataSkills struct {
	Grades []StudentTermGradeSkills `json:"grades"`
}

// TermGrades represents the grades for a particular term.
type TermGrades struct {
	Term   string  `json:"term"`
	Grades []Grade `json:"grades"`
}

// StudentTermGrades represents the grades for a student in a particular subject.
type StudentTermGrades struct {
	StudentID int          `json:"student_id"`
	SubjectID int          `json:"subject_id"`
	Terms     []TermGrades `json:"terms"`
}

// GradesData represents the overall data structure for grades.
type GradesData struct {
	Grades []StudentTermGrades `json:"grades"`
}

// Response represents a generic response message.
type Response struct {
	Message string `json:"message"`
}
