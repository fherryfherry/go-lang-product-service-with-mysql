package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

type Product struct {
	ID   int    `json:"id"`
	Sku string `json:"sku"`
	Name string `json:"name"`
	Description string `json:"description"`
	Stock int `json:"stock"`
	Price float32 `json:"price"`
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var products[] Product
var apiResponse ApiResponse

func detailProduct(c *gin.Context) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/sample_shop")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	var product Product
	var id = c.Param("id")

	err = db.QueryRow("SELECT id, sku, name, description, stock, price FROM products where id = ?", id).Scan(&product.ID,&product.Sku,&product.Name,&product.Description,&product.Stock,&product.Price)
	if err != nil {
		apiResponse.Status = 0
		apiResponse.Message = err.Error()
		c.IndentedJSON(http.StatusInternalServerError,apiResponse)
	} else {
		c.IndentedJSON(http.StatusOK, product)
	}
}

func listProduct(c *gin.Context) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/sample_shop")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query("SELECT id, sku, name, description, stock, price FROM products")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var product Product
		// for each row, scan the result into our tag composite object
		err = results.Scan(&product.ID, &product.Sku, &product.Name, &product.Description, &product.Stock, &product.Price)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		products = append(products, product)
	}

	c.IndentedJSON(http.StatusOK, products)
}

func createProduct(c *gin.Context) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/sample_shop")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO products (sku,name,description,stock,price) values (?,?,?,?,?)", c.Param("sku"), c.Param("name"), c.Param("description"), c.Param("stock"), c.Param("price"))

	if err != nil {
		apiResponse.Status = 0
		apiResponse.Message = err.Error()
		c.IndentedJSON(http.StatusInternalServerError, apiResponse)
		return
	}

	defer insert.Close()

	apiResponse.Status = 1
	apiResponse.Message = "Created!"
	c.IndentedJSON(http.StatusCreated, apiResponse)
}

func main() {
	fmt.Println("Go Product Service MySQL RAW")

	router := gin.Default()
	router.GET("/products", listProduct)
	router.GET("/products/detail/:id", detailProduct)
	router.GET("/products/create", createProduct)

	router.Run("localhost:8080")
}