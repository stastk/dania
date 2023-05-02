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
	r.HandleFunc("/db/{action}", execDB)
	r.HandleFunc("/popdb/{count}/{minchild}/{maxchild}", populateDB)
	r.HandleFunc("/ingridients", showIngridients)
	r.HandleFunc("/ingridients/{id}", showIngridients)
	r.HandleFunc("/ingridient/new", newIngridient).Methods("POST")
	r.HandleFunc("/ingridient/new/variation", newVariation).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)

}

// New Ingridient
func newIngridient(w http.ResponseWriter, r *http.Request) {

	var answer []models.Ingridient
	r.ParseForm()
	err := r.ParseForm()
	f := r.Form
	name := f.Get("name")
	json.Marshal(answer)
	//answer, err = models.NewIngridient(name)
	answer, err = models.NewIngridient(name)
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
	w.Write(jsonResp)
}

// New Ingridient
func newVariation(w http.ResponseWriter, r *http.Request) {

	var answer []models.Ingridient
	r.ParseForm()
	err := r.ParseForm()
	f := r.Form
	name := f.Get("name")
	// TODO handle error
	parentId, err := strconv.Atoi(f.Get("parent_id"))
	json.Marshal(answer)
	//answer, err = models.NewIngridient(name)
	answer, err = models.NewVarition(name, parentId)
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
	w.Write(jsonResp)
}

// Show Ingridient[s]
func showIngridients(w http.ResponseWriter, r *http.Request) {
	var answer []models.Ingridient
	vars := mux.Vars(r)
	fmt.Println(vars)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		//TODO add paranoid error. Prevent DB request
	}
	if len(vars) == 0 {
		answer, err = models.AllIngridients()
	} else if len(vars) > 0 || id >= 0 {
		answer, err = models.IngridientShow(id)
	}

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
	w.Write(jsonResp)

}

// Init/Drop DB
func execDB(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fmt.Println(vars)

	action := vars["action"]
	answer, err := models.ExecDB(action)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s\n", answer)
}

// Populate DB
func populateDB(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	count, err := strconv.Atoi(vars["count"])
	minchild, err := strconv.Atoi(vars["minchild"])
	if minchild < 1 {
		minchild = 1
	}
	maxchild, err := strconv.Atoi(vars["maxchild"])

	answer, err := models.PopulateDB(count, minchild, maxchild)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s\n", answer)

}
