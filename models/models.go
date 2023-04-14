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
		DROP TABLE IF EXISTS test1;
		DROP TABLE IF EXISTS test2;
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
		CREATE TABLE test1 (
			id INT PRIMARY KEY,
			name VARCHAR,
			description VARCHAR,
			manufacturer VARCHAR,
			color VARCHAR,
			inventory int CHECK (inventory > 0)
		  );
		  CREATE TABLE test2 (
			id INT PRIMARY KEY,
			name VARCHAR,
			description VARCHAR,
			manufacturer VARCHAR,
			color VARCHAR,
			inventory int CHECK (inventory > 0)
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
