//  models/student.go

package models

import "github.com/lib/pq"

// Student represents a student entity
type Student struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ClassroomID int    `json:"classroom_id"`
	// Add any other fields related to students here
}

type StudentWithClassroomAndSubjects struct {
	ID               int            `json:"id"`
	Name             string         `json:"name"`
	Classroom        string         `json:"classroom"`
	AssignedSubjects pq.StringArray `json:"assigned_subjects"`
}

// StudentSubject represents the many-to-many relationship between students and subjects
type StudentSubject struct {
	StudentID int `json:"student_id"`
	SubjectID int `json:"subject_id"`
}
