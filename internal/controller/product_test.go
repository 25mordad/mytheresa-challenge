package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func TestImportProduct400(t *testing.T) {
	log := log.Println
	log("TestImportProduct400")
	assert := assert.New(t)
	/////test: wrong Content-Type
	body := map[string]interface{}{}
	req, rr := fakePost("/api/v0/products", body, "POST")
	var Cache = cache.New(18*time.Hour, 18*time.Hour)
	h := NewBaseHandler(Cache)
	router := mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.ImportProduct)
	router.ServeHTTP(rr, req)
	assert.Equal(400, rr.Code)

	/////test: price string

	body = map[string]interface{}{
		"products": []map[string]interface{}{
			map[string]interface{}{
				"sku":      "000001",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    "89000",
			},
		},
	}
	req, rr = fakePost("/api/v0/products", body, "POST")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rr, req)
	assert.Equal(400, rr.Code)

}

func TestImportProduct200(t *testing.T) {
	log := log.Println
	log("TestImportProduct200")
	assert := assert.New(t)
	/////test: 200 product updated
	body := map[string]interface{}{
		"products": []map[string]interface{}{
			map[string]interface{}{
				"sku":      "000001",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    89000,
			},
			map[string]interface{}{
				"sku":      "000002",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    99000,
			},
			map[string]interface{}{
				"sku":      "000003",
				"name":     "Ashlington leather ankle boots",
				"category": "boots",
				"price":    71000,
			},
			map[string]interface{}{
				"sku":      "000004",
				"name":     "Naima embellished suede sandals",
				"category": "sandals",
				"price":    79500,
			},
			map[string]interface{}{
				"sku":      "000005",
				"name":     "Nathane leather sneakers",
				"category": "sneakers",
				"price":    59000,
			},
		},
	}

	req, rr := fakePost("/api/v0/products", body, "POST")
	req.Header.Set("Content-Type", "application/json")

	var Cache = cache.New(18*time.Hour, 18*time.Hour)
	h := NewBaseHandler(Cache)
	router := mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.ImportProduct)
	router.ServeHTTP(rr, req)
	assert.Equal(200, rr.Code)

	products := h.getProduct()
	assert.Equal(products.Products[1].Sku, "000002")
	assert.Equal(products.Products[1].Name, "BV Lean leather ankle boots")
	assert.Equal(products.Products[1].Category, "boots")
	assert.Equal(products.Products[1].Price, 99000)
	assert.Equal(products.Products[4].Sku, "000005")
	assert.Equal(products.Products[4].Price, 59000)

}

func TestFilterProduct200(t *testing.T) {
	log := log.Println
	log("TestFilterProduct200")
	assert := assert.New(t)
	/////First import product
	body := map[string]interface{}{
		"products": []map[string]interface{}{
			map[string]interface{}{
				"sku":      "000001",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    89000,
			},
			map[string]interface{}{
				"sku":      "000002",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    99000,
			},
			map[string]interface{}{
				"sku":      "000003",
				"name":     "Ashlington leather ankle boots",
				"category": "boots",
				"price":    71000,
			},
			map[string]interface{}{
				"sku":      "000004",
				"name":     "Naima embellished suede sandals",
				"category": "sandals",
				"price":    79500,
			},
			map[string]interface{}{
				"sku":      "000005",
				"name":     "Nathane leather sneakers",
				"category": "sneakers",
				"price":    59000,
			},
		},
	}

	req, rr := fakePost("/api/v0/products", body, "POST")
	req.Header.Set("Content-Type", "application/json")

	var Cache = cache.New(18*time.Hour, 18*time.Hour)
	h := NewBaseHandler(Cache)
	router := mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.ImportProduct)
	router.ServeHTTP(rr, req)
	assert.Equal(200, rr.Code)

	products := h.getProduct()
	assert.Equal(products.Products[1].Sku, "000002")

	//////////////////////////// end import product
	//////second test filter product get all products
	body = map[string]interface{}{}
	req, rr = fakePost("/api/v0/products", body, "GET")
	router = mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.FilterProduct)
	router.ServeHTTP(rr, req)
	var plist ExpoProduct
	err := json.NewDecoder(rr.Body).Decode(&plist)
	assert.Equal(err, nil)
	assert.Equal(len(plist.Products), 5)
	assert.Equal(plist.Products[1].Sku, "000002")
	//check discount
	assert.Equal(plist.Products[0].Price.Original, 89000)
	assert.Equal(plist.Products[0].Price.Final, 62300)
	assert.Equal(plist.Products[0].Price.DiscountPercentage, "30%")
	assert.Equal(plist.Products[0].Price.Currency, "EUR")

	/////filter product by category
	req, rr = fakePost("/api/v0/products?category=boots", body, "GET")
	router = mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.FilterProduct)
	router.ServeHTTP(rr, req)
	err = json.NewDecoder(rr.Body).Decode(&plist)
	assert.Equal(err, nil)
	assert.Equal(len(plist.Products), 3)
	assert.Equal(plist.Products[0].Sku, "000001")
	assert.Equal(plist.Products[1].Sku, "000002")
	assert.Equal(plist.Products[2].Sku, "000003")

	/////filter product by price less than
	req, rr = fakePost("/api/v0/products?priceLessThan=87000", body, "GET")
	router = mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.FilterProduct)
	router.ServeHTTP(rr, req)
	err = json.NewDecoder(rr.Body).Decode(&plist)
	assert.Equal(err, nil)
	assert.Equal(len(plist.Products), 3)
	assert.Equal(plist.Products[0].Sku, "000003")
	assert.Equal(plist.Products[1].Sku, "000004")
	assert.Equal(plist.Products[2].Sku, "000005")

	/////filter product by price less than and category
	req, rr = fakePost("/api/v0/products?category=boots&priceLessThan=87000", body, "GET")
	router = mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.FilterProduct)
	router.ServeHTTP(rr, req)
	err = json.NewDecoder(rr.Body).Decode(&plist)
	assert.Equal(err, nil)
	assert.Equal(len(plist.Products), 1)
	assert.Equal(plist.Products[0].Sku, "000003")

}

func TestFilterProduct404(t *testing.T) {
	log := log.Println
	log("TestFilterProduct200")
	assert := assert.New(t)
	/////First import product
	body := map[string]interface{}{
		"products": []map[string]interface{}{
			map[string]interface{}{
				"sku":      "000001",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    89000,
			},
			map[string]interface{}{
				"sku":      "000002",
				"name":     "BV Lean leather ankle boots",
				"category": "boots",
				"price":    99000,
			},
			map[string]interface{}{
				"sku":      "000003",
				"name":     "Ashlington leather ankle boots",
				"category": "boots",
				"price":    71000,
			},
			map[string]interface{}{
				"sku":      "000004",
				"name":     "Naima embellished suede sandals",
				"category": "sandals",
				"price":    79500,
			},
			map[string]interface{}{
				"sku":      "000005",
				"name":     "Nathane leather sneakers",
				"category": "sneakers",
				"price":    59000,
			},
		},
	}

	req, rr := fakePost("/api/v0/products", body, "POST")
	req.Header.Set("Content-Type", "application/json")

	var Cache = cache.New(18*time.Hour, 18*time.Hour)
	h := NewBaseHandler(Cache)
	router := mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.ImportProduct)
	router.ServeHTTP(rr, req)
	assert.Equal(200, rr.Code)

	products := h.getProduct()
	assert.Equal(products.Products[1].Sku, "000002")

	//////////////////////////// end import product
	//////second test filter not found
	body = map[string]interface{}{}
	req, rr = fakePost("/api/v0/products?category=test", body, "GET")
	router = mux.NewRouter()
	router.HandleFunc("/api/v0/products", h.FilterProduct)
	router.ServeHTTP(rr, req)
	assert.Equal(404, rr.Code)

}

//////fakePost
func fakePost(url string, body map[string]interface{}, method string) (*http.Request, *httptest.ResponseRecorder) {

	jsonBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewReader(jsonBytes))
	rr := httptest.NewRecorder()

	return req, rr
}
