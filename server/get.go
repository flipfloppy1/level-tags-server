package main

import (
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

var pgMap = pgtype.NewMap()

type levelEntry struct {
	ID       int      `db:"id" json:"id"`
	Tags     []string `db:"tags" json:"tags"`
	Name     string   `db:"name" json:"name"`
	Category string   `db:"category" json:"category"`
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

	var levels []levelEntry
	err = c.db.Select(&levels, query, args...)

	if err != nil {
		log.Printf("error getting level info: %v", err)
		return
	}

	if len(levels) == 0 {
		writeErr(res, "No levels found", 404)
		return
	}

	tagNameSet := make(map[string]struct{})
	for _, l := range levels {
		for _, t := range l.Tags {
			tagNameSet[t] = struct{}{}
		}
	}
	tagNames := slices.Collect(maps.Keys(tagNameSet))

	query, args, err = sqlx.In("select * from tags where name = any(?)", tagNames)
	if err != nil {
		log.Printf("error constructing query for tags: %v\n", err)
		writeErr(res, "Internal server error", 500)
		return
	}
	query = c.db.Rebind(query)
	var tags []tagEntry
	err = c.db.Select(&tags, query, args...)
	if err != nil {
		log.Printf("error querying tags for levels: %v\n", err)
		writeErr(res, "Internal server error", 500)
		return
	}

	tagMap := make(map[string]tagEntry)
	for _, t := range tags {
		tagMap[t.Name] = t
	}

	returnTags := make(map[int]map[string][]string)

	for _, l := range levels {
		levelTags := make(map[string][]string)
		for _, t := range l.Tags {
			tag := tagMap[t]
			levelTags[tag.Category] = append(levelTags[tag.Category], tag.Name)
		}
		returnTags[l.ID] = levelTags
	}

	resp, err := json.Marshal(returnTags)
	if err != nil {
		log.Printf("error formatting level tag info: %v\n", err)
		writeErr(res, "Internal server error", 500)
		return
	}

	res.WriteHeader(200)
	res.Write(resp)
}
