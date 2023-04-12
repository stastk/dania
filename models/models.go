package models

import (
	"database/sql"
)

var DB *sql.DB

type Ingridient struct {
	Name          string
	VariationName string
	UnitOfMeasure string
	Quantity      float32
}

func AllIngridients() ([]Ingridient, error) {
	rows, err := DB.Query("SELECT * FROM ingridients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ings []Ingridient

	for rows.Next() {
		var ing Ingridient

		err := rows.Scan(&ing.Name, &ing.VariationName, &ing.UnitOfMeasure, &ing.Quantity)
		if err != nil {
			return nil, err
		}

		ings = append(ings, ing)
	}

}
