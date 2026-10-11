package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	databaseConnectionString := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("pgx", databaseConnectionString)
	if err != nil {
		log.Fatalf("error connecting to DB: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/credits", &credits{db})
	http.ListenAndServe(":8080", mux)
}
