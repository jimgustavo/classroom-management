package models

type Teacher struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TeacherData struct {
	ID                          int    `json:"id"`
	School                      string `json:"school"`
	SchoolYear                  string `json:"school_year"`
	SchoolHours                 string `json:"school_hours"`
	Country                     string `json:"country"`
	City                        string `json:"city"`
	TeacherID                   int    `json:"teacher_id"`
	TeacherFullName             string `json:"teacher_full_name"`
	TeacherBirthday             string `json:"teacher_birthday"`
	TeacherIDNumber             string `json:"id_number"`
	LaborDependencyRelationship string `json:"labor_dependency_relationship"`
	InstitutionalEmail          string `json:"institutional_email"`
	Phone                       string `json:"phone"`
	Principal                   string `json:"principal"`
	VicePrincipal               string `json:"vice_principal"`
	Dece                        string `json:"dece"`
	Inspector                   string `json:"inspector"`
}
