// database/grade_label.go

package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jimgustavo/classroom-management/models"
)

// CreateGradeLabel inserts a new grade label into the database
func CreateGradeLabel(gradeLabel *models.GradeLabel) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO grade_labels (label, date, skill, teacher_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := db.QueryRow(query, gradeLabel.Label, gradeLabel.Date, gradeLabel.Skill, gradeLabel.TeacherID).Scan(&gradeLabel.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllGradeLabels retrieves all grade labels from the database
func GetAllGradeLabels() ([]models.GradeLabel, error) {
	query := "SELECT id, label, date, skill FROM grade_labels"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gradeLabels []models.GradeLabel
	for rows.Next() {
		var gradeLabel models.GradeLabel
		if err := rows.Scan(&gradeLabel.ID, &gradeLabel.Label, &gradeLabel.Date, &gradeLabel.Skill); err != nil {
			return nil, err
		}
		gradeLabels = append(gradeLabels, gradeLabel)
	}

	return gradeLabels, nil
}

// GetGradeLabelByID retrieves a specific grade label by its ID from the database
func GetGradeLabelByID(id string) (*models.GradeLabel, error) {
	var gradeLabel models.GradeLabel
	query := "SELECT id, label, date, skill FROM grade_labels WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&gradeLabel.ID, &gradeLabel.Label, &gradeLabel.Date, &gradeLabel.Skill)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("grade label not found")
		}
		return nil, err
	}
	return &gradeLabel, nil
}

// UpdateGradeLabel updates an existing grade label in the database
func UpdateGradeLabel(gradeLabel *models.GradeLabel) error {
	query := "UPDATE grade_labels SET label = $1, date = $2, skill = $3 WHERE id = $4"
	_, err := db.Exec(query, gradeLabel.Label, gradeLabel.Date, gradeLabel.Skill, gradeLabel.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGradeLabel deletes a grade label from the database by its ID
func DeleteGradeLabel(id string) error {
	query := "DELETE FROM grade_labels WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// AssignGradeLabelToSubject assigns a grade label to a subject for a specific term in the database
func AssignGradeLabelToSubjectByTerm(subjectID string, gradeLabelID, termID int) error {
	_, err := db.Exec("INSERT INTO grade_labels_subjects (subject_id, grade_label_id, term_id) VALUES ($1, $2, $3)", subjectID, gradeLabelID, termID)
	if err != nil {
		return err
	}
	return nil
}

// GetGradeLabelsForSubject retrieves all grade labels assigned to a subject for a specific term from the database
func GetGradeLabelsForSubject(subjectID int, termID int) ([]models.GradeLabel, error) {
	log.Printf("Attempting to add subject %d to classroom %d", subjectID, termID)
	var gradeLabels []models.GradeLabel

	rows, err := db.Query("SELECT id, label FROM grade_labels WHERE id IN (SELECT grade_label_id FROM grade_labels_subjects WHERE subject_id = $1 AND term_id = $2)", subjectID, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gradeLabel models.GradeLabel
		if err := rows.Scan(&gradeLabel.ID, &gradeLabel.Label); err != nil {
			return nil, err
		}
		gradeLabels = append(gradeLabels, gradeLabel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gradeLabels, nil
}
