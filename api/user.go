package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example/hello/db"
	"example/hello/util"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

func SetupUserRoutes(r *mux.Router, queries *db.Queries) {
	r.HandleFunc("/users/{id}", createGetUserHandler(queries)).Methods("GET")
	r.HandleFunc("/users/{id}", createUpdateUserHandler(queries)).Methods("PUT")
	r.HandleFunc("/users/{id}", createDeleteUserHandler(queries)).Methods("DELETE")
}

func createGetUserHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		idStr := vars["id"]
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
			return
		}

		user, err := queries.GetUser(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
			return
		}

		response := struct {
			ID          int64  `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number,omitempty"`
		}{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber.String,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func createUpdateUserHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
			return
		}

		var req struct {
			FirstName   string `json:"first_name,omitempty"`
			LastName    string `json:"last_name,omitempty"`
			Email       string `json:"email,omitempty"`
			Password    string `json:"password,omitempty"`
			PhoneNumber string `json:"phone_number,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid JSON request"}`, http.StatusBadRequest)
			return
		}

		existingUser, err := queries.GetUser(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
			return
		}

		updatedFirstName := existingUser.FirstName
		updatedLastName := existingUser.LastName
		updatedEmail := existingUser.Email
		updatedPasshash := existingUser.Passhash
		updatedSalt := existingUser.Salt
		updatedPhoneNumber := existingUser.PhoneNumber

		if req.FirstName != "" {
			updatedFirstName = req.FirstName
		}
		if req.LastName != "" {
			updatedLastName = req.LastName
		}
		if req.Email != "" {
			updatedEmail = req.Email
		}
		if req.Password != "" {
			salt := "randomsalt"
			hashedPassword, err := util.HashAndSalt(req.Password, salt)
			if err != nil {
				http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
				return
			}
			updatedPasshash = hashedPassword
			updatedSalt = salt
		}
		if req.PhoneNumber != "" {
			updatedPhoneNumber = pgtype.Text{String: req.PhoneNumber, Valid: true}
		}

		err = queries.UpdateUser(r.Context(), db.UpdateUserParams{
			ID:          userID,
			FirstName:   updatedFirstName,
			LastName:    updatedLastName,
			Email:       updatedEmail,
			Passhash:    updatedPasshash,
			Salt:        updatedSalt,
			PhoneNumber: updatedPhoneNumber,
		})
		if err != nil {
			http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
			return
		}

		updatedUser, err := queries.GetUser(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error": "Failed to retrieve updated user"}`, http.StatusInternalServerError)
			return
		}

		response := struct {
			ID          int64  `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number,omitempty"`
		}{
			ID:          updatedUser.ID,
			FirstName:   updatedUser.FirstName,
			LastName:    updatedUser.LastName,
			Email:       updatedUser.Email,
			PhoneNumber: updatedUser.PhoneNumber.String,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func createDeleteUserHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		idStr := vars["id"]
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
			return
		}

		_, err = queries.GetUser(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
			return
		}

		err = queries.DeleteUser(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error": "Failed to delete user"}`, http.StatusInternalServerError)
			return
		}

		response := struct {
			Message string `json:"message"`
			ID      int64  `json:"id"`
		}{
			Message: "User deleted successfully",
			ID:      userID,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
