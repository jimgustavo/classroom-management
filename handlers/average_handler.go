// handlers/average_handlers.go

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
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
