package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type vwStats struct {
	Player    string `json:"Player"`
	Date      string `json:"Date"`
	Mode      string `json:"Game_mode"`
	KD        string `json:"KDratio"`
	Kills     string `json:"Kills"`
	Deaths    string `json:"Deaths"`
	Headshots string `json:"Headshots"`
	Placement string `json:"Placement"`
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", printRoutes).Methods(http.MethodGet)
	api.HandleFunc("", printRoutes).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("", put).Methods(http.MethodPut)
	api.HandleFunc("", delete).Methods(http.MethodDelete)
	api.HandleFunc("/players", players).Methods(http.MethodGet)
	api.HandleFunc("/kills", mostKills).Methods(http.MethodGet)
	api.HandleFunc("/pr/{userName:[a-z_]{3,10}}/kills", prKills).Methods(http.MethodGet)
	api.HandleFunc("/pr/{userName:[a-z_]{3,10}}/kd", prKD).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":9990", r))
}
