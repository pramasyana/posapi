package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
	"gorm.io/datatypes"
)

type Product struct {
	Base
}

func (b *Product) Mount(group fiber.Router) {
	group.Post("", b.CreateProduct)
	group.Get("", b.GetAllProduct)
	group.Get("/:id", b.GetProduct)
	group.Put("/:id", b.UpdateProduct)
	group.Delete("/:id", b.DeleteProduct)
}

func (b *Product) GetAllProduct(c *fiber.Ctx) error {
	if c.GetReqHeaders()["Authorization"] == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   fiber.Map{},
		})
	}

	count := models.GetProductCount(c)
	if count == 0 && len(c.Query("categoryId")) == 0 && len(c.Query("q")) == 0 {
		for i := 1; i <= 10; i++ {
			data := new(models.Categories)
			data.Name = "Kategori " + strconv.Itoa(i)
			models.CreateCategory(*data)
		}

		for i := 1; i <= 5; i++ {
			data := new(models.Products)
			data.Name = "Produk " + strconv.Itoa(i)
			data.Stock = int64(i * 7)
			data.Price = float64(i * 78900)
			data.Image = "https://images.tokopedia.net/img/cache/500-square/hDjmkQ/2020/11/26/001f1c6e-d068-484f-9333-c3fa4129ef26.jpg"
			data.CategoryId = int64(i)
			data.Discount = nil
			models.CreateProduct(*data)
		}
	}

	count = models.GetProductCount(c)
	var products []models.ProductList = models.FindAllProduct(c)
	if len(products) == 0 {
		products = []models.ProductList{}
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"products": products,
			"meta": fiber.Map{
				"total": count,
				"limit": limit,
				"skip":  skip,
			},
		},
	})
}

func (b *Product) GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product, err := models.FindProductCategory(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    product,
	})
}

func (b *Product) CreateProduct(c *fiber.Ctx) error {
	data := new(models.Products)
	parsingData := new(models.ProductsParsing)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	err := json.Unmarshal(c.Body(), &parsingData)
	if err != nil || len(parsingData.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create product",
			"error":   fiber.Map{},
		})
	}

	switch i := parsingData.Price.(type) {
	case float64:
		data.Price = i
	case float32:
		data.Price = float64(i)
	case int64:
		data.Price = float64(i)
	case int:
		data.Price = float64(i)
	case string:
		floatNum, _ := strconv.ParseFloat(i, 64)
		data.Price = floatNum
	}

	data.Name = parsingData.Name
	data.Stock = parsingData.Stock
	data.Image = parsingData.Image
	data.Sku = parsingData.Sku
	data.CategoryId = parsingData.CategoryId

	discount, err := json.Marshal(parsingData.Discount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create product",
			"error":   fiber.Map{},
		})
	}
	data.Discount = datatypes.JSON(string(discount))

	item, err := models.CreateProduct(*data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create product",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    item,
	})
}

func (b *Product) UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product, err := models.FindProduct(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   fiber.Map{},
		})
	}

	var p struct {
		CategoryId int64   `json:"categoryId,omitempty"`
		Name       string  `json:"name,omitempty"`
		Image      string  `json:"image,omitempty"`
		Price      float64 `json:"price,omitempty"`
		Stock      int64   `json:"stock,omitempty"`
	}

	err = json.Unmarshal(c.Body(), &p)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   fiber.Map{},
		})
	}

	if p.CategoryId != 0 {
		product.CategoryId = p.CategoryId
	}

	if p.Name != "" {
		product.Name = p.Name
	} else {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "product data doesn't match",
			"error":   fiber.Map{},
		})
	}

	if p.Image != "" {
		product.Image = p.Image
	}

	if p.Price != 0 {
		product.Price = p.Price
	}

	if p.Stock != 0 {
		product.Stock = p.Stock
	}

	err = models.SaveProduct(product)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func (b *Product) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed deleting product",
			"error":   fiber.Map{},
		})
	}

	_, err = models.FindProduct(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   fiber.Map{},
		})
	}

	models.DeleteProduct(id)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
