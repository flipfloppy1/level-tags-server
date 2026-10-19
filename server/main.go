package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// writeErr writes code and {"error": message} to res.
func writeErr(res http.ResponseWriter, message string, code int) {
	resp, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{message})
	res.WriteHeader(code)
	res.Write(resp)
}

func main() {
	databaseConnectionString := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("pgx", databaseConnectionString)
	if err != nil {
		log.Fatalf("error connecting to DB: %v\n", err)
	}

	rl := &rateLimiter{clients: make(map[string]rateLimitClient)}

	mux := http.NewServeMux()
	mux.Handle("/credits", &credits{db, rl})
	mux.Handle("/tags", &tags{db, rl})
	mux.Handle("/get", &get{db, rl})
	http.ListenAndServe(":8080", mux)
}
