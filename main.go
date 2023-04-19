package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"main/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	// Initalize the sql.DB connection pool and assign it to the models.DB
	// global variable.
	models.DB, err = sql.Open("postgres", "user=gouser password=gopass dbname=prania_exp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	models.DBpr, err = sqlx.Connect("postgres", "user="+models.DBusername+" password="+models.DBusepassword+" dbname="+models.DBname+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/initdb", initDB)
	http.HandleFunc("/dropdb", dropDB)
	http.HandleFunc("/popdb", populateDB)
	http.HandleFunc("/ingridients", ingridientsIndex)
	http.ListenAndServe(":3000", nil)

}

func ingridientsIndex(w http.ResponseWriter, r *http.Request) {
	answer, err := models.AllIngridients()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(answer)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return

}

// ingridientsIndex sends a HTTP response listing all ingridients.
func dropDB(w http.ResponseWriter, r *http.Request) {
	answer, err := models.DropDB()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, answer := range answer {
		fmt.Fprintf(w, "%s\n", answer)
	}
}

// ingridientsIndex sends a HTTP response listing all ingridients.
func initDB(w http.ResponseWriter, r *http.Request) {
	answer, err := models.InitDB()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, answer := range answer {
		fmt.Fprintf(w, "%s\n", answer)
	}
}

// ingridientsIndex sends a HTTP response listing all ingridients.
func populateDB(w http.ResponseWriter, r *http.Request) {
	answer, err := models.PopulateDB()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, answer := range answer {
		fmt.Fprintf(w, "%s\n", answer)
	}
}
