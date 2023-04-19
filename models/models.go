package models

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DBpr *sqlx.DB
var DBname = "prania_exp"
var DBusername = "gouser"
var DBusepassword = "gopassword"

var schema = `
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
	`

// Temporary struct for getting Variations as string
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

	ingridients := []Ingridient{}
	rows, err := DBpr.Queryx(`
		SELECT
			i.id,
			i.name,
			jsonb_agg(v) as variations
			FROM Ingridients i
			LEFT JOIN IngridientsVariations v ON v.ingridient_id = i.id
			GROUP BY i.id;
	`)
	if err != nil {
		log.Fatal(err)
	}
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
		// TODO remove [null] as part of array
		if record.Variations != "[null]" {
			toappend := Ingridient{record.Id, record.Name, variations}
			ingridients = append(ingridients, toappend)
		}

	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return ingridients, nil

}

// Drop tables #db
func DropDB() (string, error) {
	answer, err := DBpr.Queryx(`
		DROP TABLE IF EXISTS IngridientsVariations;
		DROP TABLE IF EXISTS Ingridients;
	`)
	if err != nil {
		return "", err
	}

	defer answer.Close()

	return "Drop table -done", err
}

// Create tables without anything #db
func InitDB() (string, error) {
	DBpr.MustExec(schema)
	return "Create table -done", nil
}

// Populate tables with some data #db
func PopulateDB() (string, error) {
	tx := DBpr.MustBegin()
	tx.MustExec(`
		INSERT INTO Ingridients
			(name)
		VALUES
			('Ganash'),
			('Salt'),
			('Pepper'),
			('Milk');

		INSERT INTO IngridientsVariations
			(name, ingridient_id)
		VALUES 
			('Ganash of the north', 1),
			('Ganash of the south', 1),
			('Salt of the sea', 2),
			('Ganash of the west', 1),
			('Pepper of the Iron Man', 3),
			('Ganash of the north', 1),
			('Ganash of the east', 1),
			('Dr.Pepper', 3);
	`)
	tx.Commit()
	return "Populate -done", nil
}
