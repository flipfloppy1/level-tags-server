package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type tagEntry struct {
	Username  string `db:"username" json:"username"`
	Role      string `db:"role" json:"role"`
	AccountID int    `db:"accountid" json:"accountid"`
}

type get struct {
	db *sqlx.DB
}

func (c *get) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var entries []tagEntry
	err := c.db.Select(&entries, `SELECT id, tags FROM levels WHERE id = ANY($1)`)

	if err != nil {
		log.Printf("error getting credits info: %v\n", err)
		resp, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{"Internal server error"})
		res.WriteHeader(500)
		res.Write(resp)
		return
	}

	resp, err := json.Marshal(entries)
	if err != nil {
		log.Printf("error formatting credits info: %v\n", err)
		resp, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{"Internal server error"})
		res.WriteHeader(500)
		res.Write(resp)
		return
	}

	res.WriteHeader(200)
	res.Write(resp)
}
