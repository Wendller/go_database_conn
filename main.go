package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type Product struct {
	Id    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		Id:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func InsertProduct(db *sql.DB, product *Product) error {
	query, err := db.Prepare("INSERT INTO products(id, name, price) values($1, $2, $3)")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(product.Id, product.Name, product.Price)
	if err != nil {
		return err
	}

	log.Println("Insertion executed successfully")

	return nil
}

func UpdateProduct(db *sql.DB, product *Product) error {
	query, err := db.Prepare("UPDATE products SET name = $1, price = $2 WHERE id = $3")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(product.Name, product.Price, product.Id)
	if err != nil {
		return err
	}

	log.Println("Update executed successfully")

	return nil
}

func findById(db *sql.DB, productId string) (*Product, error) {
	query, err := db.Prepare("SELECT * FROM products WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var product Product

	err = query.QueryRow(productId).Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func main() {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/goexpert?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected successfully 🔌")

	product := NewProduct("Notebook", 1899.90)

	err = InsertProduct(db, product)
	if err != nil {
		log.Fatal(err)
	}

	product.Price = 2000.0
	err = UpdateProduct(db, product)
	if err != nil {
		log.Fatal(err)
	}

	queryProduct, err := findById(db, product.Id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Query result: %v", queryProduct)
}
