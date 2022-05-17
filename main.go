package main

import (
	"inventoryservice/database"
	"inventoryservice/product"
	"inventoryservice/receipt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	receipt.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
