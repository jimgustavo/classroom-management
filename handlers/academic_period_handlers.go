// handlers/academic_period_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

func CreateAcademicPeriodHandler(w http.ResponseWriter, r *http.Request) {
	var period models.AcademicPeriod
	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := database.CreateAcademicPeriod(period.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	period.ID = id
	json.NewEncoder(w).Encode(period)
}

func GetAllAcademicPeriodsHandler(w http.ResponseWriter, r *http.Request) {
	periods, err := database.GetAllAcademicPeriods()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(periods)
}

func GetAcademicPeriodByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	period, err := database.GetAcademicPeriodByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(period)
}

func UpdateAcademicPeriodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var period models.AcademicPeriod
	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateAcademicPeriod(id, period.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAcademicPeriodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteAcademicPeriod(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AssignTermToAcademicPeriodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	academicPeriodID, err := strconv.Atoi(vars["academicPeriodID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	termID, err := strconv.Atoi(vars["termID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.AssignTermToAcademicPeriod(academicPeriodID, termID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetTermsByAcademicPeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	academicPeriodID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid academic period ID", http.StatusBadRequest)
		return
	}

	terms, err := database.FetchTermsByAcademicPeriodFromDB(academicPeriodID)
	if err != nil {
		http.Error(w, "Failed to fetch terms", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(terms)
}
