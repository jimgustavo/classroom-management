package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/middleware"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	log.Println("Original password:", creds.Password)
	log.Println("Hashed password:", string(hashedPassword))

	err = database.CreateTeacher(creds.Name, creds.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to create teacher", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teacherID, err := database.AuthenticateTeacher(credentials.Email, credentials.Password)
	if err != nil {
		log.Println("Error:", err) // Add logging for error
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	teacher, err := database.GetTeacherByID(teacherID)
	if err != nil {
		log.Println("Error:", err) // Add logging for error
		http.Error(w, "Failed to retrieve teacher data", http.StatusInternalServerError)
		return
	}

	token, err := middleware.GenerateToken(teacherID)
	if err != nil {
		log.Println("Error:", err) // Add logging for error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token":      token,
		"teacher_id": teacherID,
		"teacher":    teacher,
	}

	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the token from local storage or cookies
	// In case of local storage
	w.Header().Set("Set-Cookie", "token=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure; SameSite=Strict")
	// Optionally, you can clear the teacher ID as well
	w.Header().Set("Set-Cookie", "teacher_id=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure; SameSite=Strict")

	w.WriteHeader(http.StatusOK)
}

// GetAllTeachersHandler retrieves all teachers without requiring authorization
func GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teachers, err := database.GetAllTeachers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teachers)
}

// DeleteTeacherHandler deletes a specific teacher without requiring authorization
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	teacherID := params["id"]
	if teacherID == "" {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	err := database.DeleteTeacher(teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
