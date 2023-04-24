package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"main/models"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	models.DBpr, err = sqlx.Connect("postgres", "user="+models.DBusername+" password="+models.DBusepassword+" dbname="+models.DBname+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/initdb", initDB)
	r.HandleFunc("/dropdb", dropDB)
	r.HandleFunc("/popdb/{count}/{minchild}/{maxchild}", populateDB)
	r.HandleFunc("/ingridients", ingridientsIndex)
	r.HandleFunc("/ingridient/{id}", ingridientShow)
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)

}

func ingridientShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	id, err := strconv.Atoi(vars["id"])

	answer, err := models.IngridientShow(id)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	jsonResp, err := json.Marshal(answer)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonResp)

}

func ingridientsIndex(w http.ResponseWriter, r *http.Request) {
	answer, err := models.AllIngridients()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	jsonResp, err := json.Marshal(answer)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "application/json")

	w.Write(jsonResp)

}

// ingridientsIndex sends a HTTP response listing all ingridients.
func dropDB(w http.ResponseWriter, r *http.Request) {
	answer, err := models.DropDB()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s\n", answer)
}

// ingridientsIndex sends a HTTP response listing all ingridients.
func initDB(w http.ResponseWriter, r *http.Request) {
	answer, err := models.InitDB()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s\n", answer)
}

// ingridientsIndex sends a HTTP response listing all ingridients.
func populateDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	count, err := strconv.Atoi(vars["count"])
	minchild, err := strconv.Atoi(vars["minchild"])
	if minchild < 1 {
		minchild = 1
	}
	maxchild, err := strconv.Atoi(vars["maxchild"])

	fmt.Fprintf(w, "Count: %v\n", count)
	fmt.Println("count =>", vars["count"])
	answer, err := models.PopulateDB(count, minchild, maxchild)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s\n", answer)

}
