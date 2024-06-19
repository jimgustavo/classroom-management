// handlers/average_handlers.go

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
)

// Handler function for fetching average grades by classroom ID
func GetAverageGradesByClassroomID(w http.ResponseWriter, r *http.Request) {
	// Extract the classroom ID from the URL path
	vars := mux.Vars(r)
	classroomID := vars["classroomID"]

	// Convert classroomID to an integer (assuming it's an integer)
	id, err := strconv.Atoi(classroomID)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Fetch average grades for the specified classroom ID
	averagesData, err := database.FetchAverageGradesByClassroomID(id)
	if err != nil {
		http.Error(w, "Error fetching average grades", http.StatusInternalServerError)
		log.Printf("Error fetching average grades: %v\n", err)
		return
	}

	// Convert averagesData to JSON
	responseData, err := json.Marshal(averagesData)
	if err != nil {
		http.Error(w, "Error encoding response data", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

	log.Println("Average grades retrieved successfully")
}

// Handler function for fetching average grades by classroom ID
func GetAveragesWithFactorsByClassroomID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classroomID := vars["classroomID"]

	id, err := strconv.Atoi(classroomID)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	termFactors := []models.TermFactor{}

	for queryParam, factorStrs := range queryParams {
		if len(factorStrs) > 0 {
			factor, err := strconv.ParseFloat(factorStrs[0], 32)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid factor for term %s", queryParam), http.StatusBadRequest)
				return
			}
			termFactors = append(termFactors, models.TermFactor{
				Term:   queryParam,
				Factor: float32(factor),
			})
		}
	}

	// Debug print to verify term factors
	fmt.Printf("Parsed term factors: %+v\n", termFactors)

	averagesData, err := database.FetchAveragesWithFactorsByClassroomID(id, termFactors)
	if err != nil {
		http.Error(w, "Error fetching average grades with factors", http.StatusInternalServerError)
		log.Printf("Error fetching average grades with factors: %v\n", err)
		return
	}

	responseData, err := json.Marshal(averagesData)
	if err != nil {
		http.Error(w, "Error encoding response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

	log.Println("Average grades with factors retrieved successfully")
}
