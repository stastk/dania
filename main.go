package prania

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"prania/models"

	_ "github.com/lib/pq"
)

func main() {
	var err error

	// Initalize the sql.DB connection pool and assign it to the models.DB
	// global variable.
	models.DB, err = sql.Open("postgres", "postgres://user:pass@localhost/bookstore")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ingridients", ingridientsIndex)
	http.ListenAndServe(":3000", nil)
}

// booksIndex sends a HTTP response listing all books.
func ingridientsIndex(w http.ResponseWriter, r *http.Request) {
	ings, err := models.AllIngridients()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, ing := range ings {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", ing.Name, ing.VariationName, ing.UnitOfMeasure, ing.Quantity)
	}
}
