package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"main/models"

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

	http.HandleFunc("/ingridients", ingridientsIndex)
	http.ListenAndServe(":3000", nil)
}

// ingridientsIndex sends a HTTP response listing all ingridients.
func ingridientsIndex(w http.ResponseWriter, r *http.Request) {
	ings, err := models.AllIngridients()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, ing := range ings {
		fmt.Fprintf(w, "%d : %s\n", ing.Id, ing.Name)
	}
}
