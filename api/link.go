package api

import (
	"encoding/json"
	"example/hello/db"
	"net/http"

	"example/hello/util"

	"github.com/gorilla/mux"
)

func SetupLinkRoutes(r *mux.Router, queries *db.Queries) {
	r.HandleFunc("/link/", createNewLinkHandler(queries)).Methods("POST")
}

func createNewLinkHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Destination string `json:"Destination"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid JSON request"}`, http.StatusBadRequest)
			return
		}

		if req.Destination == "" {
			http.Error(w, `{"error": "Missing required fields: Destination"}`, http.StatusBadRequest)
			return
		}

		destination, err := queries.CreateLink(r.Context(), db.CreateLinkParams{Source: util.RandStringBytes(16), Destination: req.Destination})

		if err != nil {
			http.Error(w, `{"error": "failed creating link"}`, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response := struct {
			Source string `json:"Source"`
		}{
			Source: destination.Source,
		}
		json.NewEncoder(w).Encode(response)
	}
}
