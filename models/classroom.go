// models/classroom.go
package models

type Classroom struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	TeacherID int     `json:"teacher_id"`
	Teacher   Teacher `json:"teacher"`
	// You can include other fields as needed
}
