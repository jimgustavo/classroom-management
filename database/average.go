// database/average.go

package database

import (
	"fmt"

	"github.com/jimgustavo/classroom-management/models"
)

// FetchAverageGradesByClassroomID retrieves the average grades for each term from the database for a specific classroom
func FetchAverageGradesByClassroomID(classroomID int) (models.AveragesData, error) {
	db := GetDB()

	query := `
		SELECT g.student_id, g.subject_id, t.name, AVG(g.grade) as average
        FROM grades g
        JOIN terms t ON g.term_id = t.id
        WHERE g.classroom_id = $1
        GROUP BY g.student_id, g.subject_id, t.name
        ORDER BY g.student_id, g.subject_id, t.name`
	rows, err := db.Query(query, classroomID)
	if err != nil {
		return models.AveragesData{}, fmt.Errorf("error fetching average grades: %w", err)
	}
	defer rows.Close()

	var averagesData models.AveragesData
	averageMap := make(map[int]map[int][]models.TermAverage)

	for rows.Next() {
		var studentID, subjectID int
		var term string
		var average float32

		if err := rows.Scan(&studentID, &subjectID, &term, &average); err != nil {
			return models.AveragesData{}, fmt.Errorf("error scanning row: %w", err)
		}

		if averageMap[studentID] == nil {
			averageMap[studentID] = make(map[int][]models.TermAverage)
		}

		averageMap[studentID][subjectID] = append(averageMap[studentID][subjectID], models.TermAverage{
			Term:    term,
			Average: average,
		})
	}

	for studentID, subjects := range averageMap {
		for subjectID, averages := range subjects {
			averagesData.Averages = append(averagesData.Averages, models.StudentTermAverages{
				StudentID: studentID,
				SubjectID: subjectID,
				Averages:  averages,
			})
		}
	}

	return averagesData, nil
}

func FetchAveragesWithFactorsByClassroomID(classroomID int, termFactors []models.TermFactor) (models.AveragesDataFactor, error) {
	db := GetDB()

	query := `
        SELECT g.student_id, g.subject_id, t.name, AVG(g.grade) as average
        FROM grades g
        JOIN terms t ON g.term_id = t.id
        WHERE g.classroom_id = $1
        GROUP BY g.student_id, g.subject_id, t.name
        ORDER BY g.student_id, g.subject_id, t.name`
	rows, err := db.Query(query, classroomID)
	if err != nil {
		return models.AveragesDataFactor{}, fmt.Errorf("error fetching average grades: %w", err)
	}
	defer rows.Close()

	var averagesData models.AveragesDataFactor
	averageMap := make(map[int]map[int][]models.TermAverageFactor)

	for rows.Next() {
		var studentID, subjectID int
		var term string
		var average float32

		if err := rows.Scan(&studentID, &subjectID, &term, &average); err != nil {
			return models.AveragesDataFactor{}, fmt.Errorf("error scanning row: %w", err)
		}

		factor := float32(1.0)
		for _, tf := range termFactors {
			if tf.Term == term {
				factor = tf.Factor
				break
			}
		}
		fmt.Printf("Checking term: %s, factor: %f\n", term, factor)

		aveFactor := average * factor

		if averageMap[studentID] == nil {
			averageMap[studentID] = make(map[int][]models.TermAverageFactor)
		}

		averageMap[studentID][subjectID] = append(averageMap[studentID][subjectID], models.TermAverageFactor{
			Term:      term,
			Average:   average,
			AveFactor: aveFactor,
		})
	}

	for studentID, subjects := range averageMap {
		for subjectID, averages := range subjects {
			var termAve float32
			for _, avg := range averages {
				termAve += avg.AveFactor
			}
			averagesData.Averages = append(averagesData.Averages, models.StudentTermAveragesFactor{
				StudentID: studentID,
				SubjectID: subjectID,
				Averages:  averages,
				TermAve:   termAve,
			})
		}
	}

	return averagesData, nil
}
