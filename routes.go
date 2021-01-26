package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func printRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "/players\n")
	fmt.Fprintf(w, "/kills\n")
	fmt.Fprintf(w, "/pr/{player}/kills\n")
	fmt.Fprintf(w, "/pr/{player}/kills?limit=5\n")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkLimit(param string) int {
	var limit int
	var err error
	if param != "" {
		limit, err = strconv.Atoi(param)
		if err != nil || limit < 1 {
			limit = 1
		}
		if limit >= 100 {
			limit = 99
		}
	} else {
		limit = 5
	}
	return limit
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commentID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
}

func players(w http.ResponseWriter, r *http.Request) {
	var stats []string
	q := `
    SELECT player_id FROM vw_core_players;
        `
	// open up database
	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(q)
	checkErr(err)
	defer rows.Close()
	var player string
	for rows.Next() {
		err = rows.Scan(&player)
		checkErr(err)
		stats = append(stats, player)
	}
	jsonData, err := json.MarshalIndent(stats, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func mostKills(w http.ResponseWriter, r *http.Request) {
	// go run --tags json1 web.go
	// https://github.com/mattn/go-sqlite3/issues/710
	var stats []stat
	q := `
    SELECT 
        date_key, game_mode_sub, player_id, kills 
    FROM 
        vw_stats_wz 
    WHERE 
        player_id 
    IN (SELECT * FROM vw_core_players) AND 1 
    ORDER BY 
        kills DESC LIMIT ?;
        `
	lim := 10

	// open up database
	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(q, lim)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		stat := stat{}
		err = rows.Scan(&stat.Date, &stat.Mode, &stat.Player, &stat.Kills)
		stats = append(stats, stat)
		checkErr(err)
	}
	jsonData, err := json.MarshalIndent(stats, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func prKills(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var userName string
	if val, ok := pathParams["userName"]; ok {
		userName = val
		//needs validation
	}
	query := r.URL.Query()
	param := query.Get("limit")
	limit := checkLimit(param)

	var stats []stat
	q := `
    SELECT 
        player_id, date_key, game_mode_sub, kills 
    FROM 
        vw_stats_wz 
    WHERE 
        player_id = '` + userName + `'
    ORDER BY 
        kills DESC LIMIT ?;
        `
	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(q, limit)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		stat := stat{}
		err = rows.Scan(&stat.Player, &stat.Date, &stat.Mode, &stat.Kills)
		stats = append(stats, stat)
		checkErr(err)
	}

	jsonData, err := json.MarshalIndent(stats, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

/*
func prKD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queries := mux.Vars(r)
    var stats []stat
    q := `
    SELECT
        date_key, game_mode_sub, kills
    FROM
        vw_stats_wz
    WHERE
        player_id = '?'
    ORDER BY
        kills DESC LIMIT 10;
        `
    lim := 10

    db, err := sql.Open("sqlite3", "./data.sqlite")
    checkErr(err)
    defer db.Close()

    rows, err := db.Query(q, lim)
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        stat := stat{}
        err = rows.Scan(&stat.Date, &stat.Mode, &stat.Kills)
        stats = append(stats, stat)
        checkErr(err)
    }

    json.NewEncoder(w).Encode(stats)
} */
