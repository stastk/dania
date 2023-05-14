package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"net/http"
	"strconv"

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
	r.HandleFunc("/ingridients/category/{id}", categoryOfIngridients_Show)               // Show CategoryOfIngridients
	r.HandleFunc("/ingridients/category/new", categoryOfIngridients_New).Methods("POST") // New CategoryOfIngridients
	r.HandleFunc("/ingridients/categories", categoryOfIngridients_Show)                  // All CategoryOfIngridients
	r.HandleFunc("/ingridient/new/variation", variationOfIngridient_New).Methods("POST") // New VariationOfIngridient
	r.HandleFunc("/ingridient/new", ingridient_New).Methods("POST")                      // New Ingridient
	r.HandleFunc("/ingridient/{id}", ingridient_Show)                                    // Show Ingridient
	r.HandleFunc("/ingridients", ingridient_Show)                                        // All Ingridient

	r.HandleFunc("/list/{id}", list_Show) // Show List
	r.HandleFunc("/lists", list_Show)     // All Lists

	r.HandleFunc("/db/{action}", execDB)
	r.HandleFunc("/popdb/{count}/{minchild}/{maxchild}", populateDB)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}

// New Ingridient
func ingridient_New(w http.ResponseWriter, r *http.Request) {

	var answer []models.Ingridient
	r.ParseForm()
	err := r.ParseForm()

	f := r.Form
	name := f.Get("name")

	json.Marshal(answer)
	answer, err = models.Ingridient_New(name)
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

// New IngridientsCategory
func categoryOfIngridients_New(w http.ResponseWriter, r *http.Request) {

	var answer []models.CategoryOfIngridients
	r.ParseForm()
	err := r.ParseForm()

	f := r.Form
	name := f.Get("name")

	json.Marshal(answer)
	answer, err = models.CategoryOfIngridients_New(name)
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

// New VariationOfIngridient
func variationOfIngridient_New(w http.ResponseWriter, r *http.Request) {

	var answer []models.Ingridient
	r.ParseForm()
	err := r.ParseForm()
	f := r.Form
	name := f.Get("name")
	// TODO handle error
	parentId, err := strconv.Atoi(f.Get("parent_id"))
	json.Marshal(answer)
	//answer, err = models.NewIngridient(name)
	answer, err = models.Varition_New(name, parentId)
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

// Show List[s]
func list_Show(w http.ResponseWriter, r *http.Request) {
	var answer []models.List
	vars := mux.Vars(r)
	fmt.Println(vars)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		//TODO add paranoid error. Prevent DB request
	}
	answer, err = models.List_Show(id)

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
func ingridient_Show(w http.ResponseWriter, r *http.Request) {
	var answer []models.Ingridient
	vars := mux.Vars(r)
	fmt.Println(vars)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		//TODO add paranoid error. Prevent DB request
	}
	answer, err = models.Ingridient_Show(id)

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

// Show Ingridient[s] Categories
func categoryOfIngridients_Show(w http.ResponseWriter, r *http.Request) {
	var answer []models.CategoryOfIngridients
	vars := mux.Vars(r)
	fmt.Println(vars)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		//TODO add paranoid error. Prevent DB request
	}
	answer, err = models.CategoryOfIngridients_Show(id)

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
