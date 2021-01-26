package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type stat struct {
	Player string `json_object:"player_id"`
	Date   string `json:"date"`
	Mode   string `json:"game_mode"`
	Kills  string `json:"kills"`
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
	//api.HandleFunc("/pr/{userName:[a-z_]{3,10}}/kd", prKD).Methods(http.MethodGet)
	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":9990", r))
}
