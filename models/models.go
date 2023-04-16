package models

import (
	"database/sql"
)

var DB *sql.DB

var dbname = "prania_exp"

type Ingridient struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Variations string `json:"variations"`
}

var somevar []string

// TODO solve problem with #empty_variables (somevar)
// Drop all tables #db
func DropDB() ([]string, error) {
	rows, err := DB.Query(
		`
		DROP TABLE IF EXISTS IngridientsToIngridientsVariations;
		DROP TABLE IF EXISTS Ingridients;
		DROP TABLE IF EXISTS IngridientsVariations;
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

		CREATE TABLE IF NOT EXISTS IngridientsVariations (
			id serial PRIMARY KEY,
			name VARCHAR(50) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS IngridientsToIngridientsVariations (
			ingridien_id int NOT NULL,
			ingridient_variation_id int NOT NULL,
			PRIMARY KEY (ingridien_id, ingridient_variation_id),
			FOREIGN KEY (ingridien_id) REFERENCES Ingridients(id) ON UPDATE CASCADE,
			FOREIGN KEY (ingridient_variation_id) REFERENCES IngridientsVariations(id) ON UPDATE CASCADE
		);

		INSERT INTO Ingridients
			("id", "name")
			VALUES
			('1', 'Ganash'),
			('2', 'Salt'),
			('3', 'Pepper'),
			('4', 'Milk');
			
			INSERT INTO IngridientsVariations
			("id", "name")
			VALUES
			('1', 'Ganash of the north'),
			('2', 'Ganash of the south'),
			('3', 'Salt of the sea'),
			('4', 'Ganash of the west'),
			('5', 'Pepper of the Iron Man'),
			('6', 'Ganash of the north'),
			('7', 'Ganash of the east'),
			('8', 'Dr.Pepper');

			INSERT INTO IngridientsToIngridientsVariations
			("ingridien_id", "ingridient_variation_id")
			VALUES
			('3', '8'),
			('1', '2'),
			('1', '4'),
			('1', '6'),
			('1', '7'),
			('3', '5'),
			('2', '3');
		`)

	if err != nil {
		return somevar, err
	}

	defer rows.Close()

	return somevar, err
}

func AllIngridients() ([]Ingridient, error) {
	rows, err := DB.Query(`
	SELECT
		i.id,
		i.name,
		jsonb_agg(v) as variations
		FROM Ingridients i
		LEFT JOIN IngridientsToIngridientsVariations ivar ON i.id = ivar.ingridien_id
		LEFT JOIN IngridientsVariations v ON ivar.ingridient_variation_id = v.id
		GROUP BY i.id;
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ings []Ingridient
	for rows.Next() {
		var ing Ingridient

		err := rows.Scan(&ing.Id, &ing.Name, &ing.Variations)
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
