package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func main() {
	connInfo := "postgres://postgres:1234@localhost:5432/pqgotest?sslmode=disable"

	db, err := sql.Open("postgres", connInfo)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// createProductTable(db)

	// product := Product{"panda", 9.99, true}
	// pk := insertProduct(db, product)

	// fmt.Printf("ID: %d\n", pk)

	// querySingleRow(db, 1)

	queryMultipleRow(db)

}

func queryMultipleRow(db *sql.DB) {
	data := []Product{}

	rows, err := db.Query(`SELECT name, available, price FROM product`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var name string
	var available bool
	var price float64

	for rows.Next() {
		err := rows.Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}

	fmt.Println(data)
}

func querySingleRow(db *sql.DB, index int) {
	var name string
	var available bool
	var price float64

	query := `SELECT name, available, price FROM product WHERE id = $1`
	err := db.QueryRow(query, index).Scan(&name, &available, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with ID %d", index)
		}
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\nAvailable: %t\nPrice: %f\n", name, available, price)
}

func insertProduct(db *sql.DB, product Product) int {
	query := `INSERT INTO product (name, price, available)
	VALUES ($1, $2, $3) RETURNING id`

	var pk int
	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}

	return pk
}

func createProductTable(db *sql.DB) {

	/*	product Table
		- ID
		- Name
		- Price
		- Available
		- Date Created
	*/

	query := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6, 2) NOT NULL,
		available BOOLEAN,
		created_at timestamp DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
