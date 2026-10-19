package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type tagEntry struct {
	Name        string `db:"name" json:"name"`
	Category    string `db:"category" json:"category"`
	Mode        string `db:"mode" json:"mode"`
	Color       string `db:"color" json:"color"`
	Description string `db:"description" json:"description"`
}

type tags struct {
	db *sqlx.DB
	rl *rateLimiter
}

func (t *tags) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if t.rl.Limit(res, req) {
		return
	}

	var tags []tagEntry
	err := t.db.Select(&tags, "select * from tags")
	if err != nil {
		writeErr(res, "Internal server error", 500)
		log.Printf("error getting tags: %v", err)
		return
	}

	grouped := make(map[string]map[string][]string)
	for _, t := range tags {
		cat, ok := grouped[t.Category]
		if !ok {
			cat = make(map[string][]string)
		}
		cat[t.Name] = []string{t.Mode, t.Color, t.Description}
		grouped[t.Category] = cat
	}

	resp, _ := json.Marshal(struct {
		Style    map[string][]string `json:"style"`
		Theme    map[string][]string `json:"theme"`
		Meta     map[string][]string `json:"meta"`
		Gameplay map[string][]string `json:"gameplay"`
	}{grouped["style"], grouped["theme"], grouped["meta"], grouped["gameplay"]})
	res.WriteHeader(200)
	res.Write(resp)
}
