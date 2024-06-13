// models/average.go

package models

type TermAverage struct {
	Term    string  `json:"term"`
	Average float32 `json:"average"`
}

type StudentTermAverages struct {
	StudentID int           `json:"student_id"`
	SubjectID int           `json:"subject_id"`
	Averages  []TermAverage `json:"averages"`
}

type AveragesData struct {
	Averages []StudentTermAverages `json:"averages"`
}
