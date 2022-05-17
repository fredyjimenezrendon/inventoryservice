package product

type Product struct {
	ProductId      int    `json: "productId"`
	Manufacturer   string `json: "manufacturer"`
	Sku            string `json: "sku"`
	UPC            string `json: "upc"`
	PricePerUnit   string `json: "pricePerUnit"`
	QuantityOnHand int    `json: "quantityOnHand"`
	ProductName    string `json: "productName"`
}

var products []Product
