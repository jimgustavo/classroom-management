// database/average.go

package database

import (
	"fmt"
	"log"

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
			var partialAve1, partialAve2 float32
			orderedAverages := make([]models.TermAverageFactor, len(termFactors))

			// Order averages based on the termFactors order
			for i, tf := range termFactors {
				for _, avg := range averages {
					if avg.Term == tf.Term {
						orderedAverages[i] = avg
						break
					}
				}
			}

			// Calculate partial averages
			if len(orderedAverages) > 0 {
				partialAve1 = orderedAverages[0].AveFactor
			}
			if len(orderedAverages) > 1 {
				partialAve1 += orderedAverages[1].AveFactor
			}
			if len(orderedAverages) > 2 {
				partialAve2 = orderedAverages[2].AveFactor
			}
			if len(orderedAverages) > 3 {
				partialAve2 += orderedAverages[3].AveFactor
			}

			for _, avg := range orderedAverages {
				termAve += avg.AveFactor / 2
			}

			averagesData.Averages = append(averagesData.Averages, models.StudentTermAveragesFactor{
				StudentID:   studentID,
				SubjectID:   subjectID,
				Averages:    orderedAverages,
				PartialAve1: partialAve1,
				PartialAve2: partialAve2,
				TermAve:     termAve,
			})
		}
	}

	return averagesData, nil
}

func FetchAveragesWithReinforcementByClassroomID(classroomID int, termFactors []models.TermFactor) (models.AveragesDataFactor, error) {
	db := GetDB()

	regularQuery := `
        SELECT g.student_id, g.subject_id, t.name AS term, g.grade
        FROM grades g
        JOIN terms t ON g.term_id = t.id
        WHERE g.classroom_id = $1`

	reinforcementQuery := `
        SELECT r.student_id, r.subject_id, t.name AS term, r.grade
        FROM reinforcement_grade_labels r
        JOIN terms t ON r.term_id = t.id
        WHERE r.classroom_id = $1`

	log.Printf("Executing query to fetch regular grades for classroom ID %d", classroomID)
	regularRows, err := db.Query(regularQuery, classroomID)
	if err != nil {
		log.Printf("Error executing regular grades query: %v", err)
		return models.AveragesDataFactor{}, fmt.Errorf("error fetching regular grades: %w", err)
	}
	defer regularRows.Close()

	log.Printf("Executing query to fetch reinforcement grades for classroom ID %d", classroomID)
	reinforcementRows, err := db.Query(reinforcementQuery, classroomID)
	if err != nil {
		log.Printf("Error executing reinforcement grades query: %v", err)
		return models.AveragesDataFactor{}, fmt.Errorf("error fetching reinforcement grades: %w", err)
	}
	defer reinforcementRows.Close()

	type Grade struct {
		StudentID int
		SubjectID int
		Term      string
		Grade     float32
		GradeType string
	}

	var regularGrades []Grade
	var reinforcementGrades []Grade

	log.Println("Processing regular grades query results...")
	for regularRows.Next() {
		var g Grade
		if err := regularRows.Scan(&g.StudentID, &g.SubjectID, &g.Term, &g.Grade); err != nil {
			log.Printf("Error scanning regular grade row: %v", err)
			return models.AveragesDataFactor{}, fmt.Errorf("error scanning regular grade row: %w", err)
		}
		g.GradeType = "regular"
		regularGrades = append(regularGrades, g)
	}

	log.Println("Processing reinforcement grades query results...")
	for reinforcementRows.Next() {
		var g Grade
		if err := reinforcementRows.Scan(&g.StudentID, &g.SubjectID, &g.Term, &g.Grade); err != nil {
			log.Printf("Error scanning reinforcement grade row: %v", err)
			return models.AveragesDataFactor{}, fmt.Errorf("error scanning reinforcement grade row: %w", err)
		}
		g.GradeType = "reinforcement"
		reinforcementGrades = append(reinforcementGrades, g)
	}

	log.Println("Combining regular and reinforcement grades...")
	gradeMap := make(map[int]map[int]map[string][]Grade)
	for _, g := range regularGrades {
		if gradeMap[g.StudentID] == nil {
			gradeMap[g.StudentID] = make(map[int]map[string][]Grade)
		}
		if gradeMap[g.StudentID][g.SubjectID] == nil {
			gradeMap[g.StudentID][g.SubjectID] = make(map[string][]Grade)
		}
		gradeMap[g.StudentID][g.SubjectID][g.GradeType] = append(gradeMap[g.StudentID][g.SubjectID][g.GradeType], g)
	}

	for _, g := range reinforcementGrades {
		if gradeMap[g.StudentID] == nil {
			gradeMap[g.StudentID] = make(map[int]map[string][]Grade)
		}
		if gradeMap[g.StudentID][g.SubjectID] == nil {
			gradeMap[g.StudentID][g.SubjectID] = make(map[string][]Grade)
		}
		gradeMap[g.StudentID][g.SubjectID][g.GradeType] = append(gradeMap[g.StudentID][g.SubjectID][g.GradeType], g)
	}

	log.Println("Calculating averages with factors...")
	var averagesData models.AveragesDataFactor
	for studentID, subjects := range gradeMap {
		for subjectID, gradeTypes := range subjects {
			var termAve float32
			var partialAve1, partialAve2 float32
			orderedAverages := make([]models.TermAverageFactor, len(termFactors))

			log.Printf("Processing grades for studentID=%d, subjectID=%d", studentID, subjectID)

			// Order averages based on the termFactors order
			for i, tf := range termFactors {
				var totalGrade, totalAveFactor, count float32
				var label string

				if regularAverages, ok := gradeTypes["regular"]; ok {
					for _, g := range regularAverages {
						if g.Term == tf.Term {
							totalGrade += g.Grade
							totalAveFactor += g.Grade * tf.Factor
							count++
							label = "regular"
						}
					}
				}
				if reinforcementAverages, ok := gradeTypes["reinforcement"]; ok {
					for _, g := range reinforcementAverages {
						if g.Term == tf.Term {
							totalGrade += g.Grade
							totalAveFactor += g.Grade * tf.Factor
							count++
							label = "includes_reinforcement"
						}
					}
				}

				if count > 0 {
					orderedAverages[i] = models.TermAverageFactor{
						Term:      tf.Term,
						Average:   totalGrade / count,
						AveFactor: totalAveFactor / count,
						Label:     label,
					}
				} else {
					orderedAverages[i] = models.TermAverageFactor{
						Term:      tf.Term,
						Average:   0,
						AveFactor: 0,
						Label:     "regular",
					}
				}
			}

			// Calculate partial averages
			if len(orderedAverages) > 0 {
				partialAve1 = orderedAverages[0].AveFactor
			}
			if len(orderedAverages) > 1 {
				partialAve1 += orderedAverages[1].AveFactor
			}
			if len(orderedAverages) > 2 {
				partialAve2 = orderedAverages[2].AveFactor
			}
			if len(orderedAverages) > 3 {
				partialAve2 += orderedAverages[3].AveFactor
			}

			// Calculate term average considering both partial averages
			for _, avg := range orderedAverages {
				termAve += avg.AveFactor / 2
			}

			log.Printf("Calculated partial and term averages: partialAve1=%f, partialAve2=%f, termAve=%f", partialAve1, partialAve2, termAve)

			averagesData.Averages = append(averagesData.Averages, models.StudentTermAveragesFactor{
				StudentID:   studentID,
				SubjectID:   subjectID,
				Averages:    orderedAverages,
				PartialAve1: partialAve1,
				PartialAve2: partialAve2,
				TermAve:     termAve,
			})
		}
	}

	log.Println("Completed fetching and processing averages with factors.")
	return averagesData, nil
}

func FetchAveragesWithFactorsByClassroomIDForTrimesters(classroomID int, termFactors []models.TermFactor) (models.AveragesDataTrimester, error) {
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
		return models.AveragesDataTrimester{}, fmt.Errorf("error fetching average grades: %w", err)
	}
	defer rows.Close()

	// Map to hold the averages data by student and subject
	averageMap := make(map[int]map[int][]models.TermAverageFactor)

	for rows.Next() {
		var studentID, subjectID int
		var term string
		var average float32

		if err := rows.Scan(&studentID, &subjectID, &term, &average); err != nil {
			return models.AveragesDataTrimester{}, fmt.Errorf("error scanning row: %w", err)
		}

		// Calculate the factor based on the term
		factor := float32(1.0)
		for _, tf := range termFactors {
			if tf.Term == term {
				factor = tf.Factor
				break
			}
		}

		aveFactor := average * factor

		// Initialize the map if it doesn't exist for this student and subject
		if averageMap[studentID] == nil {
			averageMap[studentID] = make(map[int][]models.TermAverageFactor)
		}

		averageMap[studentID][subjectID] = append(averageMap[studentID][subjectID], models.TermAverageFactor{
			Term:      term,
			Average:   average,
			AveFactor: aveFactor,
			Label:     "regular", // Initialize as "regular"
		})
	}

	// Process the reinforcement grades
	reinforcementQuery := `
        SELECT r.student_id, r.subject_id, t.name, r.grade
        FROM reinforcement_grade_labels r
        JOIN terms t ON r.term_id = t.id
        WHERE r.classroom_id = $1`
	reinforcementRows, err := db.Query(reinforcementQuery, classroomID)
	if err != nil {
		return models.AveragesDataTrimester{}, fmt.Errorf("error fetching reinforcement grades: %w", err)
	}
	defer reinforcementRows.Close()

	for reinforcementRows.Next() {
		var studentID, subjectID int
		var term string
		var grade float32

		if err := reinforcementRows.Scan(&studentID, &subjectID, &term, &grade); err != nil {
			return models.AveragesDataTrimester{}, fmt.Errorf("error scanning reinforcement grade row: %w", err)
		}

		// Calculate the factor based on the term
		factor := float32(1.0)
		for _, tf := range termFactors {
			if tf.Term == term {
				factor = tf.Factor
				break
			}
		}

		// If reinforcement grade exists, combine it with the regular grade
		if averages, exists := averageMap[studentID][subjectID]; exists {
			for i, avg := range averages {
				if avg.Term == term {
					// Update the average to include reinforcement
					totalGrade := avg.Average + grade
					totalAveFactor := avg.AveFactor + (grade * factor)
					count := 2 // Since we are averaging two grades now (regular and reinforcement)
					averageMap[studentID][subjectID][i].Average = totalGrade / float32(count)
					averageMap[studentID][subjectID][i].AveFactor = totalAveFactor / float32(count)
					averageMap[studentID][subjectID][i].Label = "includes_reinforcement"
				}
			}
		} else {
			// If no regular grade exists, just add the reinforcement grade
			averageMap[studentID][subjectID] = append(averageMap[studentID][subjectID], models.TermAverageFactor{
				Term:      term,
				Average:   grade,
				AveFactor: grade * factor,
				Label:     "includes_reinforcement",
			})
		}
	}

	// Final processing to calculate partial and term averages
	var averagesData models.AveragesDataTrimester
	for studentID, subjects := range averageMap {
		for subjectID, averages := range subjects {
			var termAve float32
			var partialAve1, partialAve2, partialAve3 float32
			orderedAverages := make([]models.TermAverageFactor, len(termFactors))

			// Order averages based on the termFactors order
			for i, tf := range termFactors {
				for _, avg := range averages {
					if avg.Term == tf.Term {
						orderedAverages[i] = avg
						break
					}
				}
			}

			// Calculate partial averages
			if len(orderedAverages) > 0 {
				partialAve1 = orderedAverages[0].AveFactor
			}
			if len(orderedAverages) > 1 {
				partialAve1 += orderedAverages[1].AveFactor
			}
			if len(orderedAverages) > 2 {
				partialAve2 = orderedAverages[2].AveFactor
			}
			if len(orderedAverages) > 3 {
				partialAve2 += orderedAverages[3].AveFactor
			}
			if len(orderedAverages) > 4 {
				partialAve3 = orderedAverages[4].AveFactor
			}
			if len(orderedAverages) > 5 {
				partialAve3 += orderedAverages[5].AveFactor
			}

			// Calculate term average considering both partial averages
			for _, avg := range orderedAverages {
				termAve += avg.AveFactor / 3
			}

			averagesData.Averages = append(averagesData.Averages, models.StudentTermAveragesTrimester{
				StudentID:   studentID,
				SubjectID:   subjectID,
				Averages:    orderedAverages,
				PartialAve1: partialAve1,
				PartialAve2: partialAve2,
				PartialAve3: partialAve3,
				TermAve:     termAve,
			})
		}
	}

	return averagesData, nil
}

/*
func FetchAveragesWithFactorsByClassroomIDForTrimesters(classroomID int, termFactors []models.TermFactor) (models.AveragesDataTrimester, error) {
	db := GetDB()

	// Query for regular grades
	regularQuery := `
        SELECT g.student_id, g.subject_id, t.name AS term, AVG(g.grade) as average
        FROM grades g
        JOIN terms t ON g.term_id = t.id
        WHERE g.classroom_id = $1
        GROUP BY g.student_id, g.subject_id, t.name
        ORDER BY g.student_id, g.subject_id, t.name`

	// Query for reinforcement grades
	reinforcementQuery := `
        SELECT r.student_id, r.subject_id, t.name AS term, AVG(r.grade) as average
        FROM reinforcement_grade_labels r
        JOIN terms t ON r.term_id = t.id
        WHERE r.classroom_id = $1
        GROUP BY r.student_id, r.subject_id, t.name
        ORDER BY r.student_id, r.subject_id, t.name`

	// Execute regular grades query
	regularRows, err := db.Query(regularQuery, classroomID)
	if err != nil {
		return models.AveragesDataTrimester{}, fmt.Errorf("error fetching regular grades: %w", err)
	}
	defer regularRows.Close()

	// Execute reinforcement grades query
	reinforcementRows, err := db.Query(reinforcementQuery, classroomID)
	if err != nil {
		return models.AveragesDataTrimester{}, fmt.Errorf("error fetching reinforcement grades: %w", err)
	}
	defer reinforcementRows.Close()

	// Struct to hold grade data
	type Grade struct {
		StudentID int
		SubjectID int
		Term      string
		Average   float32
		GradeType string
	}

	// Process query results into maps
	gradeMap := make(map[int]map[int]map[string]Grade)

	// Process regular grades
	for regularRows.Next() {
		var g Grade
		if err := regularRows.Scan(&g.StudentID, &g.SubjectID, &g.Term, &g.Average); err != nil {
			return models.AveragesDataTrimester{}, fmt.Errorf("error scanning regular grade row: %w", err)
		}
		g.GradeType = "regular"
		if gradeMap[g.StudentID] == nil {
			gradeMap[g.StudentID] = make(map[int]map[string]Grade)
		}
		if gradeMap[g.StudentID][g.SubjectID] == nil {
			gradeMap[g.StudentID][g.SubjectID] = make(map[string]Grade)
		}
		gradeMap[g.StudentID][g.SubjectID][g.Term] = g
	}

	// Process reinforcement grades
	for reinforcementRows.Next() {
		var g Grade
		if err := reinforcementRows.Scan(&g.StudentID, &g.SubjectID, &g.Term, &g.Average); err != nil {
			return models.AveragesDataTrimester{}, fmt.Errorf("error scanning reinforcement grade row: %w", err)
		}
		g.GradeType = "includes_reinforcement"
		if gradeMap[g.StudentID] == nil {
			gradeMap[g.StudentID] = make(map[int]map[string]Grade)
		}
		if gradeMap[g.StudentID][g.SubjectID] == nil {
			gradeMap[g.StudentID][g.SubjectID] = make(map[string]Grade)
		}
		gradeMap[g.StudentID][g.SubjectID][g.Term] = g
	}

	// Process and calculate averages
	var averagesData models.AveragesDataTrimester

	for studentID, subjects := range gradeMap {
		for subjectID, terms := range subjects {
			var termAve float32
			var partialAve1, partialAve2, partialAve3 float32
			orderedAverages := make([]models.TermAverageFactor, len(termFactors))

			// Order and process averages based on term factors
			for i, tf := range termFactors {
				if g, ok := terms[tf.Term]; ok {
					factor := tf.Factor
					orderedAverages[i] = models.TermAverageFactor{
						Term:      g.Term,
						Average:   g.Average,
						AveFactor: g.Average * factor,
						Label:     g.GradeType,
					}
				}
			}

			// Calculate partial and term averages
			if len(orderedAverages) > 0 {
				partialAve1 = orderedAverages[0].AveFactor
			}
			if len(orderedAverages) > 1 {
				partialAve1 += orderedAverages[1].AveFactor
			}
			if len(orderedAverages) > 2 {
				partialAve2 = orderedAverages[2].AveFactor
			}
			if len(orderedAverages) > 3 {
				partialAve2 += orderedAverages[3].AveFactor
			}
			if len(orderedAverages) > 4 {
				partialAve3 = orderedAverages[4].AveFactor
			}
			if len(orderedAverages) > 5 {
				partialAve3 += orderedAverages[5].AveFactor
			}

			for _, avg := range orderedAverages {
				termAve += avg.AveFactor / 3
			}

			// Append the results to the output structure
			averagesData.Averages = append(averagesData.Averages, models.StudentTermAveragesTrimester{
				StudentID:   studentID,
				SubjectID:   subjectID,
				Averages:    orderedAverages,
				PartialAve1: partialAve1,
				PartialAve2: partialAve2,
				PartialAve3: partialAve3,
				TermAve:     termAve,
			})
		}
	}

	return averagesData, nil
}
*/

/*
http://localhost:8080/classroom/3/averageswithfactors_trimesters?trimestre_1=0.7&sumativa_t1=0.3&trimestre_2=0.7&sumativa_t2=0.3&trimestre_3=0.7&sumativa_t3=0.3
*/
