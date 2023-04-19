package models

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

// Temporary struct for getting Variations
type IngridientWithVariations struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Variations string `db:"variations" json:"variations"`
}

// Ingridient and all Variations attached to it
type Ingridient struct {
	Id         int         `db:"id" json:"id"`
	Name       string      `db:"json" json:"name"`
	Variations []Variation `db:"variations" json:"variations"`
}

// Variation of specific Ingridient
type Variation struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// All Ingridients with their Variations
func AllIngridients() ([]Ingridient, error) {
	db, err := sqlx.Connect("postgres", "user=gouser password=gopass dbname=prania_exp sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	allthings := []Ingridient{}
	rows, err := db.Queryx(`
		SELECT
			i.id,
			i.name,
			jsonb_agg(v) as variations
			FROM Ingridients i
			LEFT JOIN IngridientsVariations v ON v.ingridient_id = i.id
			GROUP BY i.id;
	`)
	defer rows.Close()

	for rows.Next() {
		var record IngridientWithVariations

		if err := rows.StructScan(&record); err != nil {
			panic(err)
		}

		in := []byte(record.Variations)
		variations := []Variation{}
		err := json.Unmarshal(in, &variations)
		if err != nil {
			log.Print(err)
		}

		toappend := Ingridient{record.Id, record.Name, variations}
		allthings = append(allthings, toappend)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return allthings, nil

}

var DB *sql.DB

var dbname = "prania_exp"

var somevar []string

// TODO solve problem with #empty_variables (somevar)
// Drop all tables #db
func DropDB() ([]string, error) {
	rows, err := DB.Query(`
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
	rows, err := DB.Query(`
	CREATE TABLE IF NOT EXISTS Ingridients (
		id serial PRIMARY KEY,
		name VARCHAR(50) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS IngridientsVariations (
		id serial PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		  ingridient_id int NOT NULL,
		FOREIGN KEY (ingridient_id) REFERENCES Ingridients(id) ON UPDATE CASCADE
	);

	INSERT INTO Ingridients
		("id", "name")
		VALUES
		(1, 'Ganash'),
		(2, 'Salt'),
		(3, 'Pepper'),
		(4, 'Milk');
		
	INSERT INTO IngridientsVariations
		("id", "name", "ingridient_id")
		VALUES
		(1, 'Ganash of the north', 1),
		(2, 'Ganash of the south', 1),
		(3, 'Salt of the sea', 2),
		(4, 'Ganash of the west', 1),
		(5, 'Pepper of the Iron Man', 3),
		(6, 'Ganash of the north', 1),
		(7, 'Ganash of the east', 1),
		(8, 'Dr.Pepper', 3);
		
	`)

	if err != nil {
		return somevar, err
	}

	defer rows.Close()

	return somevar, err
}
