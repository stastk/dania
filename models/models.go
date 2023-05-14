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

	CREATE TABLE IF NOT EXISTS Units (
		id serial PRIMARY KEY,
		name VARCHAR(16) UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS Lists (
		id serial PRIMARY KEY,
		description VARCHAR(1024)
	);

	CREATE TABLE IF NOT EXISTS ListsIngridients (
		id serial PRIMARY KEY,
		list_id int NOT NULL,
		FOREIGN KEY (list_id) REFERENCES Lists(id) ON UPDATE CASCADE,
		ingridient_id int NOT NULL,
		FOREIGN KEY (ingridient_id) REFERENCES Ingridients(id) ON UPDATE CASCADE,
		ingridient_variation_id int NOT NULL,
		FOREIGN KEY (ingridient_variation_id) REFERENCES IngridientsVariations(id) ON UPDATE CASCADE,
		count int NOT NULL,
		unit_id int NOT NULL,
		FOREIGN KEY (unit_id) REFERENCES Units(id) ON UPDATE CASCADE
	);



`

//TODO Add to Lists jsonb ListsIngridients instance

// Ingridient and all Variations attached to it
type Ingridient struct {
	Id         int        `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Variations Variations `db:"variations" json:"variations"`
	Categories Categories `db:"categories" json:"categories"`
}

// IngridientInList
type IngridientInList struct {
	Id            int        `db:"id" json:"id"`
	Name          string     `db:"name" json:"name"`
	VariationName string     `db:"variation_name" json:"variation_name"`
	Categories    Categories `db:"categories" json:"categories"`
	Count         int        `db:"count" json:"count"`
	Unit          string     `db:"unit" json:"unit"`
}

// Variation of specific Ingridient
type VariationOfIngridient struct {
	Id           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	IngridientId int    `db:"ingridient_id" json:"ingridient_id"`
}

// Category of Ingridient
type CategoryOfIngridients struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// List of Ingridients
type List struct {
	Id          int         `db:"id" json:"id"`
	Description string      `db:"description" json:"description"`
	Ingridients Ingridients `db:"ingridients" json:"ingridients"`
}

type Variations []VariationOfIngridient
type Categories []CategoryOfIngridients
type Ingridients []IngridientInList

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

func (c *Categories) Scan(value interface{}) error {
	var b []byte
	switch t := value.(type) {
	case []byte:
		b = t
	default:
		return errors.New("unknown type")
	}

	return json.Unmarshal(b, &c)
}

func (c *List) Scan(value interface{}) error {
	var b []byte
	switch t := value.(type) {
	case []byte:
		b = t
	default:
		return errors.New("unknown type")
	}

	return json.Unmarshal(b, &c)
}

func Get_Lists(request string) ([]List, error) {
	lists := []List{}
	rows, err := DBpr.Queryx(request)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var record List
		if err := rows.StructScan(&record); err != nil {
			panic(err)
		}

		if err != nil {
			log.Print(err)
		}

		lists = append(lists, record)

	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return lists, nil
}

func Get_Ingridients(request string) ([]Ingridient, error) {
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

func Get_CategoriesOfIngridients(request string) ([]CategoryOfIngridients, error) {
	categories := []CategoryOfIngridients{}
	rows, err := DBpr.Queryx(request)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var record CategoryOfIngridients
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

func Ingridient_New(name string) ([]Ingridient, error) {
	request := `
		INSERT INTO Ingridients(name)
		VALUES ('` + name + `')
		RETURNING *;
	`
	return Get_Ingridients(request)
}

func Varition_New(name string, parentId int) ([]Ingridient, error) {
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
	return Get_Ingridients(request)
}

func CategoryOfIngridients_New(name string) ([]CategoryOfIngridients, error) {
	request := `
		INSERT INTO IngridientsCategories(name)
		VALUES ('` + name + `')
		RETURNING *;
	`
	return Get_CategoriesOfIngridients(request)
}

// Show Ingridient with Variations
func Ingridient_Show(id int) ([]Ingridient, error) {

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
	`

	if id > 0 {
		request += `WHERE i.id = ` + strconv.Itoa(id)
	}

	request += `	
		GROUP BY i.id;
	`

	return Get_Ingridients(request)
}

// Show CategoryOfIngridients
func CategoryOfIngridients_Show(id int) ([]CategoryOfIngridients, error) {
	request := `
		SELECT
			ic.id,
			ic.name

		FROM IngridientsCategories ic
	`
	if id > 0 {
		request += `WHERE ic.id = ` + strconv.Itoa(id)
	}

	request += `
		ORDER BY ic.id;
	`
	return Get_CategoriesOfIngridients(request)
}

// Show Ingridient with Variations
func List_Show(id int) ([]List, error) {

	request := `
		SELECT
			l.id,
			l.description,
			COALESCE(json_agg(DISTINCT li) FILTER (WHERE li.id IS NOT NULL), '[]') AS ingridietnts
		FROM Lists l
		LEFT JOIN (
			SELECT 
				id, 
				list_id, 
				ingridient_id, 
				ingridient_variation_id, 
				count, 
				unit_id
			FROM ListsIngridients
		) AS li ON li.list_id = l.id
	`

	if id > 0 {
		request += `WHERE l.id = ` + strconv.Itoa(id)
	}

	request += `	
		GROUP BY l.id;
	`

	return Get_Lists(request)
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

			INSERT INTO Units
				(name)
			VALUES 
				('g'),
				('kg'),
				('l'),
				('lb'),
				('piece');

			INSERT INTO Lists
				(description)
			VALUES 
				('Simple egg and milk omlet'),
				('Just a burger with bread, cheese and meat'),
				('Bread, cheese and meat... probably this is cheesburger'),
				('Milkshake'),
				('Some candy'),
				('Desert with... idk...'),
				('Soup'),
				('Soup'),
				(''),
				('');

			INSERT INTO ListsIngridients
				(list_id, ingridient_id, ingridient_variation_id, count, unit_id )
			VALUES 
				(1, 1, 1, 1, 1),
				(1, 2, 1, 1, 2),
				(1, 3, 1, 1, 3),
				(1, 4, 1, 1, 4),
				(2, 2, 1, 1, 5),
				(2, 3, 1, 1, 2),
				(2, 4, 1, 1, 3),
				(3, 5, 1, 1, 4),
				(4, 2, 1, 1, 1),
				(4, 3, 1, 1, 2),
				(4, 1, 1, 1, 1),
				(4, 1, 1, 1, 1),
				(4, 4, 1, 1, 5),
				(4, 2, 1, 1, 4),
				(4, 2, 1, 1, 3),
				(5, 3, 1, 1, 5),
				(5, 4, 1, 1, 2),
				(5, 6, 1, 1, 1),
				(6, 5, 1, 1, 4),
				(6, 6, 1, 2, 5),
				(7, 3, 1, 1, 2),
				(8, 2, 1, 1, 2),
				(9, 4, 1, 1, 2),
				(10, 1, 1, 1, 1),
				(10, 2, 1, 2, 1),
				(10, 6, 1, 2, 3),
				(10, 5, 1, 1, 1);
		`)

		tx.Commit()
		return "Create tables", err

	} else if action == "drop" {
		answer, err := DBpr.Queryx(`
			DROP TABLE IF EXISTS ListsIngridients;
			DROP TABLE IF EXISTS Lists;
			DROP TABLE IF EXISTS IngridientsVariations;
			DROP TABLE IF EXISTS IngridientsToCategories;
			DROP TABLE IF EXISTS IngridientsCategories;
			DROP TABLE IF EXISTS Ingridients;
			DROP TABLE IF EXISTS Units;
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
		var variation VariationOfIngridient

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
