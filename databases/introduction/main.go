package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	ID    string
	Name  string
	Price float32
}

func NewProduct(name string, price float32) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	product := NewProduct("Product 1", float32(10.6))
	err = insertProduct(db, product)

	if err != nil {
		log.Fatal(err)
	}

	product.Price = float32(20.6)

	err = updateProduct(db, product)

	if err != nil {
		log.Fatal(err)
	}

	result, err := getProduct(context.Background(), db, product.ID)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

	products, err := listAllProducts(db)

	if err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		log.Println(product)
	}

	err = deleteProduct(db, product.ID)

	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB, product *Product) error {
	statement, err := db.Prepare("INSERT INTO tb_products (id, name, price) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(product.ID, product.Name, product.Price)
	return err
}

func updateProduct(db *sql.DB, product *Product) error {
	statement, err := db.Prepare("UPDATE tb_products SET name = ?, price = ? WHERE id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(product.Name, product.Price, product.ID)

	return err
}

func getProduct(ctx context.Context, db *sql.DB, id string) (*Product, error) {
	statement, err := db.Prepare("SELECT * FROM tb_products WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	var product Product

	err = statement.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Price)

	if err != nil {
		return nil, err
	}

	return &product, err
}

func listAllProducts(db *sql.DB) ([]*Product, error) {
	statement, err := db.Prepare("SELECT * FROM tb_products")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	var products []*Product

	rows, err := statement.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product Product

		err = rows.Scan(&product.ID, &product.Name, &product.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func deleteProduct(db *sql.DB, id string) error {
	statement, err := db.Prepare("DELETE FROM tb_products WHERE id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)

	return err
}
