package database

import (
	"database/sql"

	"github.com/jimgustavo/classroom-management/models"
)

// GetAllTerms retrieves all terms from the database
func GetAllTerms() ([]models.Term, error) {
	rows, err := db.Query("SELECT id, name FROM terms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []models.Term
	for rows.Next() {
		var term models.Term
		err := rows.Scan(&term.ID, &term.Name)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}
	return terms, nil
}

// GetTerm retrieves a term by its ID
func GetTerm(id int) (*models.Term, error) {
	row := db.QueryRow("SELECT id, name FROM terms WHERE id = $1", id)
	var term models.Term
	err := row.Scan(&term.ID, &term.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &term, nil
}

// CreateTerm inserts a new term into the database
func CreateTerm(term *models.Term) error {
	err := db.QueryRow("INSERT INTO terms (name) VALUES ($1) RETURNING id", term.Name).Scan(&term.ID)
	return err
}

// UpdateTerm updates an existing term in the database
func UpdateTerm(term *models.Term) error {
	_, err := db.Exec("UPDATE terms SET name = $1 WHERE id = $2", term.Name, term.ID)
	return err
}

// DeleteTerm deletes a term from the database by ID
func DeleteTerm(id int) error {
	_, err := db.Exec("DELETE FROM terms WHERE id = $1", id)
	return err
}

func GetTermsBySubjectID(subjectID string) ([]models.Term, error) {
	rows, err := db.Query(`
        SELECT t.id, t.name
        FROM terms t
        JOIN grade_labels_subjects gls ON t.id = gls.term_id
        WHERE gls.subject_id = $1
    `, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []models.Term
	for rows.Next() {
		var term models.Term
		if err := rows.Scan(&term.ID, &term.Name); err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}

	return terms, nil
}
