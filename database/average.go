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
