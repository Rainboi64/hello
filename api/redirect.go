package api

import (
	"example/hello/db"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRedirectRoutes(r *mux.Router, queries *db.Queries) {
	r.HandleFunc("/{source}", createRedirectHandler(queries)).Methods("GET")
}

func createRedirectHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		source := vars["source"]
		if source == "" {
			http.Error(w, `{"error": "Invalid source"}`, http.StatusBadRequest)
			return
		}

		link, err := queries.GetLink(r.Context(), source)

		if err != nil {
			http.Error(w, `{"error": "Error fetching link"}`, http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, link.Destination, http.StatusFound)
	}
}
