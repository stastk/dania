package models

import (
	"database/sql"
)

var DB *sql.DB

var dbname = "prania_exp"

type Ingridient struct {
	Id   int
	Name string
}

type Ingvars struct {
	Id    int
	Name  string
	IngId int
}

var somevar []string

// Drop all tables #db
func DropDB() ([]string, error) {

	rows, err := DB.Query(
		`
		DROP TABLE IF EXISTS IngridientsToIngridientsVariations;
		
		DROP TABLE IF EXISTS Ingridients;
		DROP TABLE IF EXISTS IngridientsVariants;
		`)
	if err != nil {
		return somevar, err
	}

	defer rows.Close()

	return somevar, err
}

// Create all tables #db
func InitDB() ([]string, error) {
	rows, err := DB.Query(
		`
		CREATE TABLE IF NOT EXISTS Ingridients (
			id serial PRIMARY KEY,
			name VARCHAR(50) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS IngridientsVariants (
			id serial PRIMARY KEY,
			name VARCHAR(50) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS IngridientsToIngridientsVariations (
			ingridien_id int NOT NULL,
			ingridient_variation_id int NOT NULL,
			PRIMARY KEY (ingridien_id, ingridient_variation_id),
			FOREIGN KEY (ingridien_id) REFERENCES Ingridients(id) ON UPDATE CASCADE,
			FOREIGN KEY (ingridient_variation_id) REFERENCES IngridientsVariants(id) ON UPDATE CASCADE
		);
		`)

	if err != nil {
		return somevar, err
	}

	defer rows.Close()

	return somevar, err
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

		err := rows.Scan(&ing.Id, &ing.Name)
		if err != nil {
			return nil, err
		}

		ings = append(ings, ing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ings, nil

}

func AllIngvars() ([]Ingridient, error) {
	rows, err := DB.Query("SELECT * FROM ingridients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ings []Ingridient

	for rows.Next() {
		var ing Ingridient

		err := rows.Scan(&ing.Id, &ing.Name)
		if err != nil {
			return nil, err
		}

		ings = append(ings, ing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ings, nil

}
