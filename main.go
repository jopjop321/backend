package main

import (
	"database/sql"
	"time"
	// "log"
	// "encoding/json"
	// "fmt"
	// "path/filepath"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Define a struct to represent your data model
type Product struct {
	ID           int    `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Desc         string `json:"description"`
	Image        string `json:"image"`
	Price_Cost   int    `json:"cost_price"`
	Price_Member int    `json:"member_price"`
	Price_Normal int    `json:"normal_price"`
	Sell         int    `json:"sell"`
	Amount       int    `json:"amount"`
}

type SellProduct struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Cost_Price int    `json:"cost_price"`
	Buy_Price  int    `json:"buy_price"`
	Sell       int    `json:"sell"`
	Amount     int    `json:"amount"`
	IsMember   bool   `json:"ismember"`
}
type SellProductDate struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Date		string `json:"date"`
	Cost_Price int    `json:"cost_price"`
	Buy_Price  int    `json:"buy_price"`
	Total_price       int    `json:"total_price"`
	Amount     int    `json:"amount"`
	IsMember   bool   `json:"ismember"`
}

type Sell struct {
	Code       string `json:"code"`
	Sell       int    `json:"sell"`
	Amount     int    `json:"amount"`
}

type AddStock struct {
	Code       string `json:"code"`
	Amount     int    `json:"amount"`
}


// Database connection
// var db *sql.DB
func createConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:jobben321@tcp(localhost:3306)/jstock")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	router.GET("/products", getProducts)
	router.GET("/products/nos", getProductsNos)
	router.GET("/sells", getsell)
	router.POST("/products", createProduct)
	router.PUT("/addstock", AddStockProduct)
	router.PUT("/products", updateProduct)
	router.PUT("/products/sell", sellProduct)
	router.DELETE("/products/:code", deleteProduct)

	router.Run(":8080")
}
func getProducts(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Code, &product.Name, &product.Price_Cost, &product.Price_Member, &product.Price_Normal, &product.Desc, &product.Image, &product.Sell, &product.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}
func getsell(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM sells ORDER BY date DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []SellProductDate{}
	for rows.Next() {
		var product SellProductDate
		err := rows.Scan(&product.Code, &product.Name, &product.Date, &product.Buy_Price, &product.Cost_Price, &product.Total_price, &product.Amount, &product.IsMember)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}  

func getProductsNos(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM product WHERE amount < 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Code, &product.Name, &product.Price_Cost, &product.Price_Member, &product.Price_Normal, &product.Desc, &product.Image, &product.Sell, &product.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func createProduct(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var product Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO product (code, name, cost_price, member_price, normal_price, description, image , sell, amount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", product.Code, product.Name, product.Price_Cost, product.Price_Member, product.Price_Normal, product.Desc, product.Image, product.Sell, product.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "id": lastInsertID})
}

func updateProduct(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var product Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("UPDATE product SET name = ?, cost_price = ?, member_price = ?, normal_price = ?, description = ?, image = ? , sell = ?, amount = ? WHERE code = ?", product.Name, product.Price_Cost, product.Price_Member, product.Price_Normal, product.Desc, product.Image, product.Sell, product.Amount, product.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
func AddStockProduct(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var addstock AddStock
	if err := c.BindJSON(&addstock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("UPDATE product SET amount = ? WHERE code = ?", addstock.Amount,addstock.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func deleteProduct(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	code := c.Param("code")

	result, err := db.Exec("DELETE FROM product WHERE code = ?", code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
type sella struct {
	sell     int    `json:"sell"`
}
func sellProduct(c *gin.Context) {
	db, err := createConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	
	var sellproduct SellProduct
	if err := c.BindJSON(&sellproduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var _sell sella
	query := "SELECT sell FROM product WHERE code = ?"
	err = db.QueryRow(query, sellproduct.Code).Scan(&_sell.sell)
    if err != nil {
        panic(err.Error())
    }

	result, err := db.Exec("UPDATE product SET sell = ?, amount = ? WHERE code = ?", sellproduct.Sell+_sell.sell,sellproduct.Amount-sellproduct.Sell, sellproduct.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	result1, err := db.Prepare("INSERT INTO sells (code, name, date, buy_price, cost_price, total_price, amount , isMember) VALUES (?, ?, ?, ?, ?, ?, ?,?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer result1.Close()

	_,err = result1.Exec(sellproduct.Code, sellproduct.Name, time.Now().Format("2006-01-02 15:04:05"), sellproduct.Buy_Price, sellproduct.Cost_Price, sellproduct.Buy_Price*sellproduct.Sell, sellproduct.Sell, sellproduct.IsMember)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
