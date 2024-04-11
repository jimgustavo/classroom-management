package models

// Classroom represents a classroom entity
type Classroom struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// You can include other fields as needed
}

// SubjectWithStudents represents a subject along with students assigned to it
type SubjectWithStudents struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Students []Student `json:"students"`
}
