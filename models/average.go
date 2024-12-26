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

type TermFactor struct {
	Term   string
	Factor float32
}

type TermAverageFactor struct {
	Term      string  `json:"term"`
	Average   float32 `json:"average"`
	AveFactor float32 `json:"ave_factor"`
	Label     string  `json:"label"`
}

type StudentTermAveragesFactor struct {
	StudentID   int                 `json:"student_id"`
	SubjectID   int                 `json:"subject_id"`
	Averages    []TermAverageFactor `json:"averages"`
	PartialAve1 float32             `json:"partial_ave_1"`
	PartialAve2 float32             `json:"partial_ave_2"`
	TermAve     float32             `json:"term_ave"`
}

type StudentTermAveragesTrimester struct {
	StudentID   int                 `json:"student_id"`
	SubjectID   int                 `json:"subject_id"`
	Averages    []TermAverageFactor `json:"averages"`
	PartialAve1 float32             `json:"partial_ave_1"`
	PartialAve2 float32             `json:"partial_ave_2"`
	PartialAve3 float32             `json:"partial_ave_3"`
	TermAve     float32             `json:"term_ave"`
}

type AveragesDataFactor struct {
	Averages []StudentTermAveragesFactor `json:"averages"`
}

type AveragesDataTrimester struct {
	Averages []StudentTermAveragesTrimester `json:"averages"`
}
