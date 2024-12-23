// handlers/term_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

// GetAllTerms handles the retrieval of all terms
func GetAllTerms(w http.ResponseWriter, r *http.Request) {
	terms, err := database.GetAllTerms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(terms)
}

// GetTerm handles the retrieval of a specific term by ID
func GetTerm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	term, err := database.GetTerm(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(term)
}

// CreateTerm handles the creation of a new term
func CreateTerm(w http.ResponseWriter, r *http.Request) {
	var term models.Term
	err := json.NewDecoder(r.Body).Decode(&term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateTerm(&term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(term)
}

// UpdateTerm handles the updating of an existing term by ID
func UpdateTerm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	var term models.Term
	err = json.NewDecoder(r.Body).Decode(&term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	term.ID = id

	err = database.UpdateTerm(&term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(term)
}

// DeleteTerm handles the deletion of a term by ID
func DeleteTerm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteTerm(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
