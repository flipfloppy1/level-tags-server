package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type levelEntry struct {
	ID       int    `db:"id" json:"id"`
	Tags     string `db:"tags" json:"tags"`
	Name     string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
}

type get struct {
	db *sqlx.DB
	rl *rateLimiter
}

func (c *get) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if c.rl.Limit(res, req) {
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		writeErr(res, "Missing IDs", 400)
		return
	}
	idStrs := strings.Split(id, ",")
	ids := make([]int, 0, len(idStrs))
	for _, idStr := range idStrs {
		parsed, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		ids = append(ids, parsed)
	}

	if len(ids) == 0 {
		writeErr(res, "Invalid IDs", 400)
		return
	}

	query, args, err := sqlx.In(`select id, tags from levels where id = any(?)`, ids)
	if err != nil {
		writeErr(res, "Internal server error", 500)
		log.Printf("error constructing query for levels: %v", err)
		return
	}
	query = c.db.Rebind(query)

	var entries []levelEntry
	err = c.db.Select(&entries, query, args...)

	if err != nil {
		log.Printf("error getting level info: %v", err)
		return
	}

	if len(entries) == 0 {
		writeErr(res, "No levels found", 404)
		return
	}

	resp, err := json.Marshal(entries)
	if err != nil {
		log.Printf("error formatting level info: %v\n", err)
		writeErr(res, "Internal server error", 500)
		return
	}

	res.WriteHeader(200)
	res.Write(resp)
}
