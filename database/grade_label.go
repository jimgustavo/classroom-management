// database/grade_label.go

package database

import (
	"database/sql"
	"errors"

	"github.com/jimgustavo/classroom-management/models"
)

// CreateGradeLabel inserts a new grade label into the database
func CreateGradeLabel(gradeLabel *models.GradeLabel) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := "INSERT INTO grade_labels (label) VALUES ($1) RETURNING id"
	err := db.QueryRow(query, gradeLabel.Label).Scan(&gradeLabel.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllGradeLabels retrieves all grade labels from the database
func GetAllGradeLabels() ([]models.GradeLabel, error) {
	query := "SELECT id, label FROM grade_labels"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gradeLabels []models.GradeLabel
	for rows.Next() {
		var gradeLabel models.GradeLabel
		if err := rows.Scan(&gradeLabel.ID, &gradeLabel.Label); err != nil {
			return nil, err
		}
		gradeLabels = append(gradeLabels, gradeLabel)
	}

	return gradeLabels, nil
}

// GetGradeLabelByID retrieves a specific grade label by its ID from the database
func GetGradeLabelByID(id string) (*models.GradeLabel, error) {
	var gradeLabel models.GradeLabel
	query := "SELECT id, label FROM grade_labels WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&gradeLabel.ID, &gradeLabel.Label)
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
	query := "UPDATE grade_labels SET label = $2 WHERE id = $3"
	_, err := db.Exec(query, gradeLabel.Label, gradeLabel.ID)
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

// AssignGradeLabelToSubject assigns a grade label to a subject in the database
func AssignGradeLabelToSubject(subjectID string, gradeLabelID int) error {
	_, err := db.Exec("INSERT INTO grade_labels_subjects (subject_id, grade_label_id) VALUES ($1, $2)", subjectID, gradeLabelID)
	if err != nil {
		return err
	}
	return nil
}

// GetGradeLabelsForSubject retrieves all grade labels assigned to a subject from the database
func GetGradeLabelsForSubject(subjectID string) ([]models.GradeLabel, error) {
	var gradeLabels []models.GradeLabel

	rows, err := db.Query("SELECT id, label FROM grade_labels WHERE id IN (SELECT grade_label_id FROM grade_labels_subjects WHERE subject_id = $1)", subjectID)
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
