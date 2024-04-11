package models

import "github.com/lib/pq"

// Student represents a student entity
type Student struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ClassroomID int    `json:"classroom_id"`
	// Add any other fields related to students here
}

// StudentWithClassroomAndSubjects represents a student along with their classroom and assigned subjects
/*
type StudentWithClassroomAndSubjects struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Classroom        string    `json:"classroom"`
	AssignedSubjects []Subject `json:"assigned_subjects"`
}
*/
type StudentWithClassroomAndSubjects struct {
	ID               int            `json:"id"`
	Name             string         `json:"name"`
	Classroom        string         `json:"classroom"`
	AssignedSubjects pq.StringArray `json:"assigned_subjects"`
}
