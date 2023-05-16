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
	r.HandleFunc("/mylist/{list_id}", embededList)                                       // Show CategoryOfIngridients
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

func embededList(w http.ResponseWriter, r *http.Request) {

	var answer []models.List
	vars := mux.Vars(r)
	fmt.Println(vars)
	list_id, err := strconv.Atoi(vars["list_id"])
	if err != nil {
		//TODO add paranoid error. Prevent DB request
	}

	if list_id > 0 {
		answer, err = models.List_Show(list_id)
	}

	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	list_html := `
		<style>
			body{
				padding: 100px 630px;
				margin: 0;
			}
			.list table,
			.list thead,
			.list tbody,
			.list tfoot,
			.list tr,
			.list th,
			.list td {
				width: auto;
				height: auto;
				margin: 0;
				padding: 0;
				border: none;
				border-collapse: inherit;
				border-spacing: 0;
				border-color: inherit;
				vertical-align: inherit;
				text-align: left;
				font-weight: inherit;
				-webkit-border-horizontal-spacing: 0;
				-webkit-border-vertical-spacing: 0;
			}
			.list{
				background: #ffffff;
				border-radius: 4px;
				overflow: hidden;
				display: block;
				font-family: Arial, sans-serif;
				border: 1px solid #eee;
				padding: 16px;
			}
			.list .description{
				padding: 16px 8px;
				display: block;
				color: #888;
				text-align: center;
			}
			.list table{
				width: 100%;
				font-size: 18px;
			}
			tr:hover td{
				background: #efefef;
			}
			.list td{
				border: 0px solid;
			}
			.list td:first-of-type{
				border-radius: 6px 0 0 6px;
			}
			.list td:last-of-type{
				border-radius: 0 6px 6px 0;
			}
			.list td.ingridient{
				padding: 0 8px;
			}
			.list td.count{
				ax-width: 100px;
				min-width: 64px;
				width: 100px;
				padding: 12px 0;
			}
			.line{
				display: flex;
				align-items: center;
			}
			.line .dots{
				background-image: linear-gradient(to right, black 33%, rgba(255,255,255,0) 0%);
				background-position: bottom;
				background-size: 3px 1px;
				background-repeat: repeat-x;
				height: 1px;
				width: 100%;
			}
			.line .name{
				display: flex;
				color: #222;
				flex: 0 0 auto;
				flex-wrap: nowrap;
				padding-right: 8px;
			}
		</style>
		<div class="list">
			<table>`

	i := 0

	for ingridients_count := range answer[0].Ingridients {
		// TODO use flex instead of table
		for ingridients_count >= i {
			list_html += `
				<tr for="v_` + strconv.Itoa(answer[0].Ingridients[i].IngridientVariationId) + strconv.Itoa(answer[0].Ingridients[i].IngridientId) + strconv.Itoa(answer[0].Ingridients[i].Count) + `">
					<td class="ingridient">
						<label for="v_` + strconv.Itoa(answer[0].Ingridients[i].IngridientVariationId) + strconv.Itoa(answer[0].Ingridients[i].IngridientId) + strconv.Itoa(answer[0].Ingridients[i].Count) + `" class="line">
							<input type="checkbox" id="v_` + strconv.Itoa(answer[0].Ingridients[i].IngridientVariationId) + strconv.Itoa(answer[0].Ingridients[i].IngridientId) + strconv.Itoa(answer[0].Ingridients[i].Count) + `">
							<span class="name">` + answer[0].Ingridients[i].VariationName + `</span>
							<span class="dots"></span>
						</label>
					</td>
					<td class="count">
						<strong>` + strconv.Itoa(answer[0].Ingridients[i].Count) + `</strong> ` + answer[0].Ingridients[i].UnitName + `
					</td>
				</tr>`
			i++
		}
	}

	list_html += `
		</table>
	`

	if len(answer[0].Description) > 0 {
		list_html += `<span class="description">` + answer[0].Description + `</span>`
	}

	list_html += `
			</div>
		</body>
	`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, list_html)
	return
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
