package main

import (
	"os"
	// "fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const productEndpoint = "/products/:id"



type Product struct {
    gorm.Model
    Name  string         `json:"name"`
    Price float64        `json:"price"`
	CategoryID uint      `json:"category_id"`
    Category   Category  `json:"category" gorm:"foreignKey:id"` 
}

type Category struct {
    gorm.Model
    Name  string  `json:"name"`
}

type Payment struct {
    gorm.Model
    amount float64 `json:"amount"`
	name string `json:"name"`
    surname string `json:"surname"`
    email string `json:"email"`
}

var products = map[string]*Product{}
var db *gorm.DB = nil
var err error

func main() {
	e := echo.New()

	if _, err = os.Stat("./test.db"); err == nil {
		os.Remove("./test.db")
	}

    db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    
    db.AutoMigrate(&Product{}, &Category{}, &Payment{})

	category := Category{Name: "Laptops"}
    db.FirstOrCreate(&category)

	category2 := Category{Name: "Mobile Phones"}
    db.Create(&category2)

    laptop := Product{Name: "Laptop", Price: 1500.0, CategoryID: category.ID, Category: category,}
    db.FirstOrCreate(&laptop)

	phone1 := Product{Name: "Samsung G5", Price: 5100.0, CategoryID: category2.ID, Category: category2,}
    db.Create(&phone1)

	phone2 := Product{Name: "Apple Iphone 15", Price: 3494.0, CategoryID: category2.ID, Category: category2,}
    db.Create(&phone2)

	phone3 := Product{Name: "Xiaomi Redmi Note 13", Price: 1299.0, CategoryID: category2.ID, Category: category2,}
    db.Create(&phone3)


	e.GET("/products", GetAllProducts)
    e.GET(productEndpoint, GetProduct)
    e.POST("/products", CreateProduct)
    e.PUT(productEndpoint, UpdateProduct)
    e.DELETE(productEndpoint, DeleteProduct)

	e.GET("/pay", GetAllPayments)
	e.POST("/pay", AddPayment)


	//  <<CORS config>> -- start
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return next(c)
		}
	})
	e.Use(middleware.CORS())
	// <<CORS config>> -- end
	e.Logger.Fatal(e.Start(":1323"))
}


func GetAllProducts(c echo.Context) error {
	var products []Product
	db.Preload("Category").Find(&products)
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, product)
}

func CreateProduct(c echo.Context) error {
	product := new(Product)
	if err := c.Bind(product); err != nil {
		return err
	}
	db.Create(&product)
	return c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	if err := c.Bind(&product); err != nil {
		return err
	}
	db.Save(&product)
	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	db.Delete(&product)
	return c.NoContent(http.StatusNoContent)
}


func GetAllPayments(c echo.Context) error {
	var payments []Payment
	db.Find(&payments)
	// fmt.Println(payments)
	return c.JSON(http.StatusOK, payments)
}

func AddPayment(c echo.Context) error {
	payment := new(Payment)
	if err := c.Bind(payment); err != nil {
		return err
	}
	db.Create(&payment)
	// fmt.Println(payment)
	return c.JSON(http.StatusCreated, payment)
}