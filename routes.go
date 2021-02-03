package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func jsonWrite(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonData := json.NewEncoder(w)
	jsonData.SetIndent("", "  ")
	jsonData.Encode(data)
}

func printRoutes(w http.ResponseWriter, r *http.Request) {
	var plist []string
	plist = append(plist, "/players")
	plist = append(plist, "/kills")
	plist = append(plist, "/pr/{player}/kills")
	plist = append(plist, "/pr/{player}/kills?limit=5")
	plist = append(plist, "/pr/{player}/kd")
	jsonWrite(w, r, plist)
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
	jsonWrite(w, r, stats)
}

func mostKills(w http.ResponseWriter, r *http.Request) {
	var stats []vwStats
	q := `
    SELECT 
	player_id, date_key, game_mode_sub, kdRatio, kills, deaths, headshots, teamPlacement 
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
		stat := vwStats{}
		err = rows.Scan(&stat.Player, &stat.Date, &stat.Mode, &stat.KD, &stat.Kills, &stat.Deaths, &stat.Headshots, &stat.Placement)
		stats = append(stats, stat)
		checkErr(err)
	}
	jsonWrite(w, r, stats)
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

	var stats []vwStats
	q := `
    SELECT 
        player_id, date_key, game_mode_sub, kdRatio, kills, deaths, headshots, teamPlacement 
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
		stat := vwStats{}
		err = rows.Scan(&stat.Player, &stat.Date, &stat.Mode, &stat.KD, &stat.Kills, &stat.Deaths, &stat.Headshots, &stat.Placement)
		stats = append(stats, stat)
		checkErr(err)
	}
	jsonWrite(w, r, stats)
}

func prKD(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var userName string
	if val, ok := pathParams["userName"]; ok {
		userName = val
		//needs validation
	}
	query := r.URL.Query()
	param := query.Get("limit")
	limit := checkLimit(param)

	var stats []vwStats
	q := `
    SELECT 
		player_id, date_key, game_mode_sub, kdRatio, kills, deaths, headshots, teamPlacement 
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
		stat := vwStats{}
		err = rows.Scan(&stat.Player, &stat.Date, &stat.Mode, &stat.KD, &stat.Kills, &stat.Deaths, &stat.Headshots, &stat.Placement)
		stats = append(stats, stat)
		checkErr(err)
	}
	jsonWrite(w, r, stats)
}
