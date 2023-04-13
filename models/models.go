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

func InitDB() error {
	result, err := DB.Query("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		return nil, err
	}
	defer result.Close()

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
