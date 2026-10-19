package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type creditEntry struct {
	Username  string `db:"username" json:"username"`
	Role      string `db:"role" json:"role"`
	AccountID string `db:"accountid" json:"accountid"`
}

type credits struct {
	db *sqlx.DB
	rl *rateLimiter
}

func (c *credits) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if c.rl.Limit(res, req) {
		return
	}

	var entries []creditEntry
	err := c.db.Select(&entries, `select username, role, accountid from credits
		order by case role
			when 'developer' then 1
			when 'admin' then 2
			when 'moderator' then 3
			when 'helper' then 4
		else 5 end`)

	if err != nil {
		log.Printf("error getting credits info: %v\n", err)
		writeErr(res, "Internal server error", 500)
		return
	}

	resp, err := json.Marshal(entries)
	if err != nil {
		log.Printf("error formatting credits info: %v\n", err)
		writeErr(res, "Internal server error", 500)
		return
	}

	res.WriteHeader(200)
	res.Write(resp)
}
