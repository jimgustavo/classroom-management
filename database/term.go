package database

import (
	"database/sql"

	"github.com/jimgustavo/classroom-management/models"
)

// GetAllTerms retrieves all terms from the database
func GetAllTerms() ([]models.Term, error) {
	rows, err := db.Query("SELECT id, name, academic_period_id FROM terms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []models.Term
	for rows.Next() {
		var term models.Term
		err := rows.Scan(&term.ID, &term.Name, &term.AcademicPeriodID)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}
	return terms, nil
}

// GetTerm retrieves a term by its ID
func GetTerm(id int) (*models.Term, error) {
	row := db.QueryRow("SELECT id, name, academic_period_id FROM terms WHERE id = $1", id)
	var term models.Term
	err := row.Scan(&term.ID, &term.Name, &term.AcademicPeriodID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &term, nil
}

// CreateTerm inserts a new term into the database
func CreateTerm(term *models.Term) error {
	err := db.QueryRow("INSERT INTO terms (name, academic_period_id) VALUES ($1, $2) RETURNING id", term.Name, term.AcademicPeriodID).Scan(&term.ID)
	return err
}

// UpdateTerm updates an existing term in the database
func UpdateTerm(term *models.Term) error {
	_, err := db.Exec("UPDATE terms SET name = $1, academic_period_id = $2 WHERE id = $3", term.Name, term.AcademicPeriodID, term.ID)
	return err
}

// DeleteTerm deletes a term from the database by ID
func DeleteTerm(id int) error {
	_, err := db.Exec("DELETE FROM terms WHERE id = $1", id)
	return err
}
