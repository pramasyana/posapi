package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

type Products struct {
	gorm.Model
	ProductId    int64          `gorm:"Not Null" json:"productId"`
	Name         string         `json:"name"`
	Stock        int64          `json:"stock"`
	Price        float64        `json:"price"`
	Image        string         `json:"image"`
	Sku          string         `json:"sku"`
	CategoriesId int64          `json:"categoriesId"`
	Discount     datatypes.JSON `json:"discount"`
}

type ProductsParsing struct {
	Name       string      `json:"name"`
	Stock      int64       `json:"stock"`
	Price      interface{} `json:"price"`
	Image      string      `json:"image"`
	Sku        string      `json:"sku"`
	CategoryId int64       `json:"categoryId"`
	Discount   interface{} `json:"discount"`
}

type ProductList struct {
	ProductId int64        `json:"productId"`
	Sku       string       `json:"sku"`
	Name      string       `json:"name"`
	Stock     int64        `json:"stock"`
	Price     float64      `json:"price"`
	Image     string       `json:"image"`
	Category  CategoryList `json:"category"`
	Discount  interface{}  `json:"discount"`
}

func FindAllProduct(c *fiber.Ctx) []ProductList {
	var products []ProductList
	var product ProductList

	db := GetDB().Table("products").Select(
		"products.product_id, products.sku, products.name, products.stock, products.price, products.image, products.discount, products.categories_id, categories.name as category_name")
	db = db.Joins("left join categories on products.categories_id=categories.category_id")
	db = db.Where("products.deleted_at is NULL")

	if len(c.Query("categoryId")) > 0 {
		db = db.Where("products.categories_id = ?", c.Query("categoryId"))
	}

	if len(c.Query("q")) > 0 {
		db = db.Where("products.name LIKE ?", `%`+c.Query("q")+`%`)
	}

	limit := 10
	if len(c.Query("limit")) > 0 {
		limit, _ = strconv.Atoi(c.Query("limit"))
	}

	skip := 0
	if len(c.Query("skip")) > 0 {
		skip, _ = strconv.Atoi(c.Query("skip"))
	}

	db = db.Limit(limit)
	db = db.Offset(skip)

	rows, _ := db.Rows()
	defer rows.Close()

	for rows.Next() {
		var (
			discount *json.RawMessage
		)
		rows.Scan(
			&product.ProductId,
			&product.Sku,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Image,
			&discount,
			&product.Category.CategoryId,
			&product.Category.Name,
		)

		if discount != nil {
			json.Unmarshal([]byte(*discount), &product.Discount)
		}

		products = append(products, product)
	}

	return products
}

func GetProductCount(c *fiber.Ctx) int64 {
	var products []Products

	db := GetDB()

	if len(c.Query("categoryId")) > 0 {
		db = db.Where("categories_id = ?", c.Query("categoryId"))
	}

	if len(c.Query("q")) > 0 {
		db = db.Where("name LIKE ?", `%`+c.Query("q")+`%`)
	}

	count := db.Find(&products).RowsAffected

	return count
}

func FindProduct(id int) (Products, error) {
	var product Products

	err := GetDB().Where("product_id = ?", id).First(&product).Error
	if err != nil {
		err = errors.New("product not found")
	}

	return product, err
}

func FindProductCategory(id int) (ProductList, error) {
	var product ProductList

	db := GetDB().Table("products").Select(
		"products.product_id, products.sku, products.name, products.stock, products.price, products.image, products.categories_id, categories.name as category_name")
	db = db.Where("products.deleted_at is NULL AND products.product_id = ?", id)

	rows, err := db.Joins("left join categories on products.categories_id=categories.category_id").Rows()
	if err != nil {
		return product, err
	}

	defer rows.Close()

	if !rows.Next() {
		return product, errors.New("Product Not Found")
	}

	rows.Scan(
		&product.ProductId,
		&product.Sku,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.Category.CategoryId,
		&product.Category.Name,
	)

	return product, err
}

func CreateProduct(product Products) (Products, error) {
	// Get Max productId
	var maxProduct Products
	var category Categories

	count := GetDB().Table("categories").Where("category_id = ?", product.CategoriesId).Find(&category).RowsAffected
	if count == 0 {
		return product, errors.New("category not found")
	}

	GetDB().Raw(`
		SELECT COALESCE(MAX(product_id) + 1, 1) as product_id
		FROM products
		`).Scan(
		&maxProduct,
	)

	product.Sku = fmt.Sprintf("ID%03d", maxProduct.ProductId)
	product.ProductId = maxProduct.ProductId

	err := GetDB().Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func SaveProduct(product Products) error {
	var category Categories

	count := GetDB().Table("categories").Where("category_id = ?", product.CategoriesId).Find(&category).RowsAffected
	if count == 0 {
		return errors.New("category not found")
	}

	err := GetDB().Table("products").Where("product_id = ?", product.ProductId).Update(product).Error

	return err
}

func DeleteProduct(id int) int64 {
	count := GetDB().Where("product_id = ?", id).Delete(&Products{}).RowsAffected

	return count
}
