package api

import (
	"encoding/json"
	"example/hello/db"
	"example/hello/util"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

func SetupAuthRoutes(r *mux.Router, queries *db.Queries) {
	r.HandleFunc("/login", createLoginHandler(queries)).Methods("POST")
	r.HandleFunc("/register", createRegisterHandler(queries)).Methods("POST")

}

func createLoginHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid JSON request"}`, http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, `{"error": "Missing required fields: email, password"}`, http.StatusBadRequest)
			return
		}

		val, err := queries.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, `{"error": "User not found"}`, http.StatusBadRequest)
		}

		salt := "randomsalt"

		if !util.VerifyPassword(req.Password, salt, val.Passhash) {
			http.Error(w, `{"error": "Invalid password"}`, http.StatusInternalServerError)
			return
		}

		tokenString, err := util.CreateToken(req.Email)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		response := struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func createRegisterHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req struct {
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			Password    string `json:"password"`
			PhoneNumber string `json:"phone_number,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid JSON request"}`, http.StatusBadRequest)
			return
		}

		if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
			http.Error(w, `{"error": "Missing required fields: first_name, last_name, email, password"}`, http.StatusBadRequest)
			return
		}

		salt := "randomsalt"
		hashedPassword, err := util.HashAndSalt(req.Password, salt)
		if err != nil {
			http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
			return
		}

		var phoneNumber pgtype.Text
		if req.PhoneNumber != "" {
			phoneNumber = pgtype.Text{String: req.PhoneNumber, Valid: true}
		}

		// Create user in database
		user, err := queries.CreateUser(r.Context(), db.CreateUserParams{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			Passhash:    hashedPassword,
			Salt:        salt,
			PhoneNumber: phoneNumber,
		})

		if err != nil {
			http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
			return
		}

		tokenString, err := util.CreateToken(req.Email)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		response := struct {
			ID          int64  `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number,omitempty"`
			Token       string `json:"token"`
		}{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber.String,
			Token:       tokenString,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
