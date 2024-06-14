package models

// Term represents a term in the system
type Term struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TeacherID int    `json:"teacher_id"`
}
