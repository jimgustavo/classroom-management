package models

// Term represents a term in the system
type Term struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	AcademicPeriodID int    `json:"academic_period_id"`
}
