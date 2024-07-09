// database/academic_period.go

package database

import (
	"database/sql"

	"github.com/jimgustavo/classroom-management/models"
)

func CreateAcademicPeriod(name string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO academic_periods (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAllAcademicPeriods() ([]models.AcademicPeriod, error) {
	rows, err := db.Query("SELECT id, name FROM academic_periods")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []models.AcademicPeriod
	for rows.Next() {
		var period models.AcademicPeriod
		if err := rows.Scan(&period.ID, &period.Name); err != nil {
			return nil, err
		}
		periods = append(periods, period)
	}
	return periods, nil
}

func GetAcademicPeriodByID(id int) (models.AcademicPeriod, error) {
	var period models.AcademicPeriod
	err := db.QueryRow("SELECT id, name FROM academic_periods WHERE id = $1", id).Scan(&period.ID, &period.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return period, nil
		}
		return period, err
	}
	return period, nil
}

func UpdateAcademicPeriod(id int, name string) error {
	_, err := db.Exec("UPDATE academic_periods SET name = $1 WHERE id = $2", name, id)
	return err
}

func DeleteAcademicPeriod(id int) error {
	_, err := db.Exec("DELETE FROM academic_periods WHERE id = $1", id)
	return err
}

func AssignTermToAcademicPeriod(academicPeriodID, termID int) error {
	_, err := db.Exec("INSERT INTO academic_period_terms (academic_period_id, term_id) VALUES ($1, $2)", academicPeriodID, termID)
	return err
}

func FetchTermsByAcademicPeriodFromDB(academicPeriodID int) ([]models.Term, error) {

	rows, err := db.Query(`
        SELECT t.id, t.name
        FROM terms t
        JOIN academic_period_terms apt ON t.id = apt.term_id
        WHERE apt.academic_period_id = $1
    `, academicPeriodID)
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
