package main

import (
	"context"
	"example/hello/api"
	"example/hello/db"
	"example/hello/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func run() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "dbname=mydb")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	users, err := queries.ListUsers(ctx)

	if err != nil {
		return err
	}

	log.Println(users)

	password := "topsecret"
	salt := ""

	passhash, _ := util.HashAndSalt(password, salt)

	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		FirstName: "Yaman",
		LastName: "Alhalabi",
		Email: "yaman.alhalabi@maddo.com",
		Passhash: passhash,
		Salt: salt,
		PhoneNumber: pgtype.Text{String: "+9639872134###", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(user)

	return nil
}

func main() {
	// Set up database connection
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "dbname=mydb")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	r := mux.NewRouter()

	api.SetupUserRoutes(r, queries)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Hello, Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}