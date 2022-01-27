# Mytheresa Challenge

By implementing MVC architecture in Go. This will allow us to expand the system more efficiently and become more organized.

Implementing services was the focus of the first version. We can improve performance with the following version.

## API

### PUT /api/v0/products
Update the list of products
Content Type `application/json`

@todo: add Authorization Bearer Token

Sample:
```sh
{
   "products":[
      {
         "sku":"000001",
         "name":"BV Lean leather ankle boots",
         "category":"boots",
         "price":89000
      },
      {
         "sku":"000002",
         "name":"BV Lean leather ankle boots",
         "category":"boots",
         "price":99000
      },
      {
         "sku":"000003",
         "name":"Ashlington leather ankle boots",
         "category":"boots",
         "price":71000
      },
      {
         "sku":"000004",
         "name":"Naima embellished suede sandals",
         "category":"sandals",
         "price":79500
      },
      {
         "sku":"000005",
         "name":"Nathane leather sneakers",
         "category":"sneakers",
         "price":59000
      }
   ]
}
```



#### Responses:

200 OK When is updated
400 Bad Request When there is a failure in the request format or the payload can't be unmarshalled.

@todo 401 StatusUnauthorized



### GET /api/v0/products
Can be filtered (query string) by:
- category
- priceLessThan

Returns a list of Product with the given discounts applied when necessary:
- Products in the boots category have a 30% discount.
- The product with sku = 000003 has a 15% discount.
- When multiple discounts collide, the biggest discount must be applied

The prices are integers for example, 100.00€ would be 10000.

#### Responses:

200 OK when product/s found
404 StatusNotFound if there is no product

- price.currency is always EUR
- When a product does not have a discount, price.final and price.original should be the same number and discount_percentage should be null.
- When a product has a discount price.original is the original price, price.final is the amount with the discount applied and discount_percentage represents the applied discount with the % sign.

#### Accept
Content Type `application/json`

Sample product with a discount of 30% and product without a discount:
```sh
{
   "products":[
      {
         "sku":"000001",
         "name":"BV Lean leather ankle boots",
         "category":"boots",
         "price":{
            "original":89000,
            "final":62300,
            "discount_percentage":"30%",
            "currency":"EUR"
         }
      },
      {
         "sku":"000002",
         "name":"BV Lean leather ankle boots 2",
         "category":"boots",
         "price":{
            "original":89000,
            "final":89000,
            "discount_percentage":null,
            "currency":"EUR"
         }
      }
   ]
}
```



## Development
#### Cache
As we do not use a database, we use a cache to save data. We can use DB in the following versions.

#### Rout

Routing uses the `gorilla/mux`. You can also add any new routing by adding it to this file: [rout.go](internal/server/rout.go)

Sample:

```sh
r.Handle("/status", http.HandlerFunc(controller.Status)).Methods("GET")

```
We can also add any middleware here, such as authentication.

Sample: *(Which `IsAuth` is a method of `auth` package.)*
```sh
r.Handle("/status", auth.IsAuth(http.HandlerFunc(controller.Status))).Methods("GET")

```


#### Config
Config is designed to add any configuration to the project. Here we can review all the env variables as well.

For example to change the running port for this project you can set env `PORT` to a custom port:
```sh
PORT=:8081

```

or change the default value in this file: [config.go](config/config.go)


#### CMD

- Run the server:

```sh
go run cmd/server/main.go
```

#### _Development Hints_

If you have [nodemon](https://www.npmjs.com/package/nodemon) on your local, then you can run the server by following command:
```sh
nodemon --exec go run cmd/server/main.go --signal SIGTERM
```


#### Test

Clean the cache

```sh
go clean -testcache
```

```sh
go test -cover ./...
```
OR
```sh
go test -cover -v ./...
```
