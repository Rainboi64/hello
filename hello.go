package main

import (
	"context"
	"example/hello/api"
	"example/hello/db"
	"example/hello/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	setupEnv()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv(("DB_CONNECTION_STRING")))
	println(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	r := mux.NewRouter()
	r.Use(middleware.Logging())

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	api.SetupAuthRoutes(r, queries)
	api.SetupRedirectRoutes(r, queries)

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.ProtectedHandler())

	api.SetupUserRoutes(protected, queries)
	api.SetupLinkRoutes(protected, queries)

	log.Println("Hello, Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
