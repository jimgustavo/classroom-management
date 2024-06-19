package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

func GetTermsByTeacherID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		log.Println("Invalid teacher ID:", err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	log.Println("Fetching terms for teacher ID:", teacherID)
	terms, err := database.GetTermsByTeacherID(teacherID)
	if err != nil {
		log.Println("Error fetching terms:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Fetched terms:", terms)
	if err := json.NewEncoder(w).Encode(terms); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

func GetTermsBySubjectID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectID := vars["subjectID"]

	terms, err := database.GetTermsBySubjectID(subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(terms)
}
