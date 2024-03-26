package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string
	Products []Product
}

type Product struct {
	gorm.Model
	Name         string
	Price        float64
	CategoryID   uint
	Category     Category
	SerialNumber SerialNumber
}

type SerialNumber struct {
	gorm.Model
	Number    string
	ProductID uint
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=goexpert port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	db.Exec("DELETE FROM serial_numbers")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM categories")

	// BELONGS TO
	category := Category{Name: "Eletronics"}
	db.Create(&category)

	product := Product{Name: "Dell notebook", Price: 3000.00, CategoryID: category.ID}
	db.Create(&product)

	// HAS ONE
	serialNumber := SerialNumber{Number: "123456", ProductID: product.ID}
	db.Create(&serialNumber)

	// HAS MANY
	var categories []Category
	err = db.Model(&Category{}).Preload("Products.SerialNumber").Find(&categories).Error
	if err != nil {
		log.Fatal(err)
	}

	for _, category := range categories {
		fmt.Printf("Category %v products:\n", category.Name)
		for _, product := range category.Products {
			fmt.Printf("-> %v - Serial Number: %v\n",
				product.Name, product.SerialNumber.Number)
		}
	}

	// db.Exec("DELETE FROM products")

	// db.Create([]Product{
	// 	{Name: "Dell", Price: 2000.00},
	// 	{Name: "Lenovo", Price: 3500.00},
	// 	{Name: "Mac", Price: 8000.00},
	// })

	// var product Product
	// db.First(&product, "name = ?", "Dell")
	// fmt.Println(product)

	// var products []Product
	// db.Limit(3).Find(&products)
	// for _, product := range products {
	// 	fmt.Printf("Product %d: %v\n", product.ID, product.Name)
	// }

	// db.Where("price > 3000").Find(&products)
	// for _, product := range products {
	// 	fmt.Printf("Product %d: %v\n", product.ID, product.Name)
	// }
}
