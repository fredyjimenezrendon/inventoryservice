package product

import (
	"encoding/json"
	"fmt"
	"inventoryservice/cors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	http.Handle(fmt.Sprintf("%s/products", apiBasePath), cors.MiddlewareHandler(handleProducts))
	http.Handle(fmt.Sprintf("%s/product/", apiBasePath), cors.MiddlewareHandler(handleProduct))
}

func productsHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		products, err := getAllProducts()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJson, err := json.Marshal(products)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(productsJson)
	case http.MethodPost:
		bodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var newProduct Product
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = insertProduct(newProduct)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		products = append(products, newProduct)
		writer.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

func productHandler(writer http.ResponseWriter, request *http.Request) {
	urlPathSegments := strings.Split(request.URL.Path, "product/")
	productId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	switch request.Method {
	case http.MethodGet:
		product, _ := getProduct(productId)
		productJson, err := json.Marshal(product)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(productJson)
		return
	case http.MethodPut:
		productJson, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var updatedProduct Product
		err = json.Unmarshal(productJson, &updatedProduct)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		updateProduct(updatedProduct)
		writer.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		deleteProduct(productId)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
