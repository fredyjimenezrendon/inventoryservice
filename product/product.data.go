package product

import (
	"context"
	"database/sql"
	"inventoryservice/database"
	"sync"
	"time"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func getProduct(productId int) (*Product, error) {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := database.DbConn.QueryRowContext(context, `SELECT productId,
       manufacturer,
       sku,
       upc,
       pricePerUnit,
       quantityOnHand,
       productName
       FROM products
       WHERE productId =?`, productId)
	product := &Product{}
	err := result.Scan(&product.ProductId,
		&product.Manufacturer,
		&product.Sku,
		&product.UPC,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return product, nil
}

func getAllProducts() ([]Product, error) {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(context, `SELECT productId,
       manufacturer,
       sku,
       upc,
       pricePerUnit,
       quantityOnHand,
       productName
       FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductId,
			&product.Manufacturer,
			&product.Sku,
			&product.UPC,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		products = append(products, product)
	}
	return products, nil
}

func updateProduct(product Product) error {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(context, `UPDATE products SET
       manufacturer=?,
       sku=?,
       upc=?,
       pricePerUnit=?,
       quantityOnHand=?,
       productName=?
       WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.UPC,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductId)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(context, `INSERT INTO products (
       manufacturer,
       sku,
       upc,
       pricePerUnit,
       quantityOnHand,
       productName)
       VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.UPC,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		return 0, err
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(insertedId), nil
}

func deleteProduct(productId int) error {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(context, `DELETE FROM products WHERE 
       productId = ?`,
		productId)
	if err != nil {
		return err
	}
	return nil
}
