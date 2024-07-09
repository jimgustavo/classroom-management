package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	TeacherID int    `json:"teacher_id"`
	Role      string `json:"role"` // Add role to the claims
	jwt.StandardClaims
}

type ctxKey int

const (
	keyTeacherID ctxKey = iota
	keyRole
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), keyTeacherID, claims.TeacherID)
		ctx = context.WithValue(ctx, keyRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(keyRole).(string)
		if !ok || role != "admin" {
			http.Error(w, "Unauthorized to access this route", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ProTeacherOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(keyRole).(string)
		if !ok || (role != "proteacher" && role != "admin") { // Admins also have access
			http.Error(w, "Unauthorized to access this route", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GenerateToken(teacherID int, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		TeacherID: teacherID,
		Role:      role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
