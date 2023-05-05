package models

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config //

var DBpr *sqlx.DB

// Test DB name
var DBname = "prania_exp"

// Test DB username
var DBusername = "gouser"

// Test DB password
var DBusepassword = "gopassword"

var err error

// TODO schema must be placed into another file
var schema = `
	CREATE TABLE IF NOT EXISTS Ingridients (
		id serial PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS IngridientsVariations (
		id serial PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		ingridient_id int NOT NULL,
		FOREIGN KEY (ingridient_id) REFERENCES Ingridients(id) ON UPDATE CASCADE
	);

	CREATE TABLE IF NOT EXISTS IngridientsCategories (
		id serial PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS IngridientsToCategories (
		id serial PRIMARY KEY,
		ingridient_id int NOT NULL,
		FOREIGN KEY (ingridient_id) REFERENCES Ingridients(id) ON UPDATE CASCADE,
		ingridient_category_id int NOT NULL,
		FOREIGN KEY (ingridient_category_id) REFERENCES IngridientsCategories(id) ON UPDATE CASCADE
	);

`

// Ingridient and all Variations attached to it
type Ingridient struct {
	Id                    int                   `db:"id" json:"id"`
	Name                  string                `db:"name" json:"name"`
	Variations            Variations            `db:"variations" json:"variations"`
	IngridientsCategories IngridientsCategories `db:"categories" json:"categories"`
}

// Variation of specific Ingridient
type Variation struct {
	Id           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	IngridientId int    `db:"ingridient_id" json:"ingridient_id"`
}

// Category of Ingridient
type IngridientsCategory struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Variations []Variation
type IngridientsCategories []IngridientsCategory

// Make the Variations type implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (v *Variations) Scan(value interface{}) error {
	var b []byte
	switch t := value.(type) {
	case []byte:
		b = t
	default:
		return errors.New("unknown type")
	}

	return json.Unmarshal(b, &v)
}

func (v *IngridientsCategories) Scan(value interface{}) error {
	var b []byte
	switch t := value.(type) {
	case []byte:
		b = t
	default:
		return errors.New("unknown type")
	}

	return json.Unmarshal(b, &v)
}

func GetIngridients(request string) ([]Ingridient, error) {
	ingridients := []Ingridient{}
	rows, err := DBpr.Queryx(request)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var record Ingridient
		if err := rows.StructScan(&record); err != nil {
			panic(err)
		}

		if err != nil {
			log.Print(err)
		}

		ingridients = append(ingridients, record)

	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return ingridients, nil
}

func GetIngridientsCategories(request string) ([]IngridientsCategory, error) {
	categories := []IngridientsCategory{}
	rows, err := DBpr.Queryx(request)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var record IngridientsCategory
		if err := rows.StructScan(&record); err != nil {
			panic(err)
		}

		if err != nil {
			log.Print(err)
		}

		categories = append(categories, record)

	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return categories, nil
}

func NewIngridient(name string) ([]Ingridient, error) {
	request := `
		INSERT INTO Ingridients(name)
		VALUES ('` + name + `')
		RETURNING *;
	`
	return GetIngridients(request)
}

func NewVarition(name string, parentId int) ([]Ingridient, error) {
	request := `
		INSERT INTO IngridientsVariations(name, ingridient_id)
		VALUES ('` + name + `', ` + strconv.Itoa(parentId) + `);
	
		SELECT i.id, i.name,
			COALESCE(json_agg(v) FILTER (WHERE v.id IS NOT NULL), '[]') AS variations
			FROM Ingridients i
			LEFT JOIN IngridientsVariations v ON v.ingridient_id = i.id
			WHERE i.id = ` + strconv.Itoa(parentId) + `
			GROUP BY i.id;
	`
	return GetIngridients(request)
}

func NewIngridientsCategory(name string) ([]IngridientsCategory, error) {
	request := `
		INSERT INTO IngridientsCategories(name)
		VALUES ('` + name + `')
		RETURNING *;;
	`
	return GetIngridientsCategories(request)
}

// Single Ingridient with Variations
func IngridientShow(id int) ([]Ingridient, error) {

	request := `
		SELECT
			i.id,
			i.name,
			COALESCE(json_agg(DISTINCT v) FILTER (WHERE v.id IS NOT NULL), '[]') AS variations,
			COALESCE(json_agg(DISTINCT ic) FILTER (WHERE ic.id IS NOT NULL), '[]') AS categories

		FROM Ingridients i

		LEFT JOIN (
			SELECT DISTINCT id, ingridient_id, name
			FROM IngridientsVariations
			GROUP BY name, id
		) AS v
		ON v.ingridient_id = i.id

		LEFT JOIN (
			SELECT DISTINCT id, name
			FROM IngridientsCategories
			GROUP BY name, id
		) AS ic
		ON ic.id IN (
			SELECT ingridient_category_id 
			FROM IngridientsToCategories 
			WHERE ingridient_id = i.id
		)
		WHERE i.id = ` + strconv.Itoa(id) + `
		GROUP BY i.id;
	`

	return GetIngridients(request)
}

// All Ingridients with Variations
func AllIngridients() ([]Ingridient, error) {
	request := `
		SELECT
			i.id,
			i.name,
			COALESCE(json_agg(DISTINCT v) FILTER (WHERE v.id IS NOT NULL), '[]') AS variations,
			COALESCE(json_agg(DISTINCT ic) FILTER (WHERE ic.id IS NOT NULL), '[]') AS categories

		FROM Ingridients i

		LEFT JOIN (
			SELECT DISTINCT id, ingridient_id, name
			FROM IngridientsVariations
			GROUP BY name, id
		) AS v
		ON v.ingridient_id = i.id

		LEFT JOIN (
			SELECT DISTINCT id, name
			FROM IngridientsCategories
			GROUP BY name, id
		) AS ic
		ON ic.id IN (
			SELECT ingridient_category_id 
			FROM IngridientsToCategories 
			WHERE ingridient_id = i.id
		)
		GROUP BY i.id;
	`
	return GetIngridients(request)
}

// All IngridientsCategories // TODO not working, fix that
func AllIngridientsCategories() ([]IngridientsCategory, error) {
	request := `
		SELECT
			ic.id,
			ic.name

		FROM IngridientsCategories ic
		ORDER BY ic.id;
	`
	return GetIngridientsCategories(request)
}

// Single IngridientsCategory
func IngridientsCategoryShow(id int) ([]IngridientsCategory, error) {
	request := `
		SELECT
			ic.id,
			ic.name

		FROM IngridientsCategories ic
		WHERE ic.id = ` + strconv.Itoa(id) + `
		ORDER BY ic.id;
	`
	return GetIngridientsCategories(request)
}

// Init/Drop tables #db
func ExecDB(action string) (string, error) {
	if action == "init" {
		// Create tables
		DBpr.MustExec(schema)

		// Populate DB
		tx := DBpr.MustBegin()
		tx.MustExec(`
			INSERT INTO Ingridients
				(name)
			VALUES
				('Ganash'),
				('Salt'),
				('Pepper'),
				('Onion'),
				('Sugar'),
				('Milk');

			INSERT INTO IngridientsVariations
				(name, ingridient_id)
			VALUES 
				('Ganash of the north', 1),
				('Ganash of the south', 1),
				('Salt of the sea', 2),
				('Ganash of the west', 1),
				('Pepper of the Iron Man', 3),
				('Ganash of the east', 1),
				('Chilli pepper', 3),
				('Pepper', 3),
				('Red hot', 3),
				('Spicy pepper', 3),
				('Dr.Pepper', 3);

			INSERT INTO IngridientsCategories
				(name)
			VALUES
				('Bakery'),
				('Spice'),
				('Meat'),
				('Fish'),
				('Vegan'),
				('Contain sugar');

			INSERT INTO IngridientsToCategories
				(ingridient_id, ingridient_category_id)
			VALUES 
				(1, 1),
				(1, 2),
				(1, 3),
				(2, 1),
				(3, 2),
				(4, 5),
				(4, 6);
		`)
		tx.Commit()
		return "Create tables", err

	} else if action == "drop" {
		answer, err := DBpr.Queryx(`
			DROP TABLE IF EXISTS IngridientsVariations;
			DROP TABLE IF EXISTS IngridientsToCategories;
			DROP TABLE IF EXISTS IngridientsCategories;
			DROP TABLE IF EXISTS Ingridients;
		`)
		if err != nil {
			return "", err
		}
		defer answer.Close()
		return "Drop tables", err
	}

	return "Executed", err
}

// Populate tables with some data #db
func PopulateDB(count int, minchild int, maxchild int) (string, error) {

	// TODO refactoring needed
	min := 4
	max := 16

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	tx := DBpr.MustBegin()

	for i := 0; i < count; i++ {
		randomChildCount := rand.Intn(maxchild-minchild) + minchild
		randomName := make([]byte, rand.Intn(max-min)+min)
		for li := range randomName {
			randomName[li] = letterBytes[rand.Int63()%int64(len(letterBytes))]
		}

		var ingridient Ingridient
		var variation Variation

		// Create some Ingridients
		err = tx.QueryRowx(`INSERT INTO Ingridients (name) VALUES ($1) RETURNING *;`, randomName).StructScan(&ingridient)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		log.Println("Ingridient:: ", ingridient.Id, ingridient.Name)

		for ich := 0; ich < randomChildCount; ich++ {
			// Add to created earlier Ingridients Variations
			err = tx.QueryRowx(`INSERT INTO IngridientsVariations (name, ingridient_id) VALUES ($1, $2) RETURNING *;`, string(randomName)+"_"+strconv.Itoa(ich), ingridient.Id).StructScan(&variation)
			if err != nil {
				tx.Rollback()
				return "", err
			}

		}

	}

	tx.Commit()

	return "Populate -done", nil
}
