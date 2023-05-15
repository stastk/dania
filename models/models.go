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
		ingridient_variation_name VARCHAR(50) NOT NULL,
		count int NOT NULL,
		unit_id int NOT NULL,
		FOREIGN KEY (unit_id) REFERENCES Units(id) ON UPDATE CASCADE,
		unit_name VARCHAR(16) NOT NULL
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

func (i *Ingridients) Scan(value interface{}) error {
	var b []byte
	switch t := value.(type) {
	case []byte:
		b = t
	default:
		return errors.New("unknown type")
	}

	return json.Unmarshal(b, &i)
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

func (l *List) Scan(value interface{}) error {
	var b []byte
	switch t := value.(type) {
	case []byte:
		b = t
	default:
		return errors.New("unknown type")
	}

	return json.Unmarshal(b, &l)
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
			COALESCE(json_agg(DISTINCT li) FILTER (WHERE li.id IS NOT NULL), '[]') AS ingridients
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
				('Goat milk', 6),
				('Milk powder', 6),
				('Natural sugar', 5),
				('Stevia', 5),
				('Big onion', 4),
				('Purple TOR onion', 4),
				('Just glass, not a sugar', 5),
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
				(list_id, ingridient_id, ingridient_variation_id, ingridient_variation_name, count, unit_id, unit_name)
			VALUES 
				(1, 6, 1, 'Goat milk',  10, 1, 'g'),
				(1, 5, 2, 'Stevia',  1, 2, 'kg'),
				(1, 4, 1, 'Big onion',  2, 3, 'l'),
				(1, 3, 2, 'Chilli pepper',  1, 4, 'lb'),
				(1, 2, 1, 'Salt of the sea',  4, 5, 'piece'),
				(2, 1, 1, 'Ganash of the north',  100, 1, 'g'),
				(2, 2, 1, 'Salt of the sea',  3, 2, 'kg'),
				(2, 3, 2, 'Chilli pepper',  13, 3, 'l'),
				(2, 4, 1, 'Big onion',  100, 4, 'lb'),
				(2, 5, 3, 'Just glass, not a sugar',  33, 5, 'piece'),
				(3, 6, 2, 'Milk powder',  4000, 1, 'g'),
				(4, 1, 2, 'Ganash of the south',  200, 2, 'kg'),
				(4, 2, 1, 'Salt of the sea',  10, 3, 'l'),
				(5, 3, 2, 'Chilli pepper',  2, 4, 'lb'),
				(5, 4, 2, 'Purple TOR onion',  5, 5, 'piece'),
				(5, 5, 1, 'Natural sugar',  6, 1, 'g'),
				(5, 6, 2, 'Milk powder',  2000, 2, 'kg'),
				(5, 5, 3, 'Just glass, not a sugar',  10, 3, 'l'),
				(5, 4, 1, 'Big onion',  15, 4, 'lb'),
				(5, 3, 2, 'Chilli pepper',  10, 5, 'piece'),
				(5, 2, 1, 'Salt of the sea',  93, 1, 'g'),
				(5, 1, 1, 'Ganash of the north',  100, 2, 'kg'),
				(5, 2, 1, 'Salt of the sea',  450, 3, 'l'),
				(5, 3, 3, 'Pepper',  40, 4, 'lb'),
				(5, 4, 1, 'Big onion',  500, 5, 'piece'),
				(5, 5, 2, 'Stevia',  400, 1, 'g'),
				(6, 6, 2, 'Milk powder',  90, 2, 'kg'),
				(6, 5, 2, 'Stevia',  900, 3, 'l'),
				(6, 4, 2, 'Purple TOR onion',  20, 4, 'lb'),
				(6, 3, 1, 'Pepper of the Iron Man',  48, 5, 'piece'),
				(6, 2, 1, 'Salt of the sea',  600, 1, 'g'),
				(6, 1, 3, 'Ganash of the west',  200, 2, 'kg'),
				(7, 2, 1, 'Salt of the sea',  400, 3, 'l'),
				(7, 3, 2, 'Chilli pepper',  300, 4, 'lb'),
				(7, 4, 1, 'Big onion',  100, 5, 'piece'),
				(8, 5, 2, 'Stevia',  750, 1, 'g'),
				(8, 6, 1, 'Goat milk',  1, 2, 'kg'),
				(8, 1, 2, 'Ganash of the north',  1, 3, 'l'),
				(8, 2, 1, 'Salt of the sea',  1, 4, 'lb'),
				(8, 3, 2, 'Chilli pepper',  1, 5, 'piece'),
				(8, 4, 1, 'Big onion',  60, 2, 'kg'),
				(8, 5, 2, 'Stevia',  85, 3, 'l'),
				(8, 6, 2, 'Milk powder',  35, 4, 'lb'),
				(9, 1, 2, 'Ganash of the north',  110, 5, 'piece'),
				(10, 2, 1, 'Salt of the sea',  25, 2, 'kg'),
				(10, 3, 2, 'Chilli pepper',  92, 1, 'g'),
				(10, 4, 1, 'Big onion',  1, 4, 'lb');
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
