package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	cache "github.com/patrickmn/go-cache"
)

//ImportProduct is struct for importing products
type ImportProduct struct {
	Products []iproduct `json:"products"`
}
type iproduct struct {
	Sku      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

//ExpoProduct is for exporting the Product list
type ExpoProduct struct {
	Products []Product `json:"products"`
}

//Product is Product model
type Product struct {
	Sku      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    Price  `json:"price"`
}

//Price is the price for product model
type Price struct {
	Original           int         `json:"original"`
	Final              int         `json:"final"`
	DiscountPercentage interface{} `json:"discount_percentage"`
	Currency           string      `json:"currency"`
}

//ImportProduct is going to get a list of products
func (bh *BaseHandler) ImportProduct(w http.ResponseWriter, r *http.Request) {
	log := log.Println
	log("ImportProduct")

	if r.Header.Get("Content-Type") != "application/json" {
		log("Content-Type is not application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var products ImportProduct
	err := json.NewDecoder(r.Body).Decode(&products)
	if err != nil {
		log(err)
		log("Payload can't be unmarshalled")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// log("products", products)

	//add to product to cache
	bh.cache.Set("products", &products, cache.DefaultExpiration)

}

//FilterProduct is going to get a list of products
func (bh *BaseHandler) FilterProduct(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	log := log.Println
	log("FilterProduct")
	// log("urlParams->", urlParams["category"])
	var cat string
	// var priceLessThan int
	priceLessThan := -1
	if len(urlParams["category"]) != 0 {
		cat = urlParams["category"][0]
	}

	if len(urlParams["priceLessThan"]) != 0 {
		var err error
		priceLessThan, err = strconv.Atoi(urlParams["priceLessThan"][0])
		if err != nil {
			log("priceLessThan is not int")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// plt = priceLessThan
	}

	products := bh.getProduct()
	if products == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var expProduct ExpoProduct

	for _, p := range products.Products {
		if search(p, cat, priceLessThan) {
			var price Price
			checkDiscount(&price, p)
			price.Currency = "EUR"

			expProduct.Products = append(expProduct.Products,
				Product{
					Sku:      p.Sku,
					Name:     p.Name,
					Category: p.Category,
					Price:    price,
				})
		}

	}

	// log("ExpoProduct", expProduct)

	if len(expProduct.Products) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	j, _ := json.Marshal(expProduct)
	io.WriteString(w, string(j))

}

func search(p iproduct, cat string, plt int) bool {
	if len(cat) == 0 && plt == -1 {
		return true
	}
	if len(cat) == 0 && plt != -1 {
		if p.Price <= plt {
			return true
		}
		return false
	}
	if len(cat) != 0 && plt != -1 {
		if p.Price <= plt && p.Category == cat {
			return true
		}
		return false
	}
	if len(cat) != 0 && plt == -1 {
		if p.Category == cat {
			return true
		}
		return false
	}
	return false
}

func checkDiscount(price *Price, p iproduct) {

	price.Original = p.Price
	price.Final = p.Price

	///check sku = 000003 15% discount
	///TThe product with sku = 000003 has a 15% discount.
	if p.Sku == "000003" {
		price.Original = p.Price
		price.Final = p.Price * (100 - 15) / 100
		price.DiscountPercentage = "15%"
	}

	///check boots category a 30% discount.
	///Products in the boots category have a 30% discount
	if p.Category == "boots" {
		price.Original = p.Price
		price.Final = p.Price * (100 - 30) / 100
		price.DiscountPercentage = "30%"
	}

	// log.Println("checkDiscount", price, p)
	///
}

func (bh *BaseHandler) getProduct() *ImportProduct {
	var products *ImportProduct
	data, found := bh.cache.Get("products")
	if found {
		products = data.(*ImportProduct)
	}
	return products
}
