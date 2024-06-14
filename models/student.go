package models

type Student struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ClassroomID int    `json:"classroom_id"`
	TeacherID   int    `json:"teacher_id"`
}

type StudentWithClassroomAndSubjects struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	ClassroomID   int      `json:"classroom_id"`
	TeacherID     int      `json:"teacher_id"`
	ClassroomName string   `json:"classroom_name"`
	Subjects      []string `json:"subjects"`
}

// StudentSubject represents the many-to-many relationship between students and subjects
type StudentSubject struct {
	StudentID int `json:"student_id"`
	SubjectID int `json:"subject_id"`
}
