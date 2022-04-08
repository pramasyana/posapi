package models

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ListSold struct {
	ProductID   int    `json:"productId"`
	Name        string `json:"name"`
	TotalQty    int    `json:"totalQty"`
	TotalAmount int    `json:"totalAmount"`
}

func FindAllSold(c *fiber.Ctx) []ListSold {
	var Sold []ListSold
	var products []Products
	db := GetDB()
	db.Select("id, name, product_id").Find(&products)

	for _, val := range products {
		paymentType := ListSold{
			ProductID: int(val.ID),
			Name:      val.Name,
		}

		var totalQty int
		var totalAmount int
		db := GetDB().Table("order_products").Select(
			"SUM(products.price * order_products.qty) AS totalAmount, SUM(order_products.qty) AS totalQty")
		db = db.Where("order_products.product_id = " + strconv.Itoa(int(val.ProductId)))
		db.Joins("JOIN products on order_products.product_id = products.id").Row().Scan(&totalAmount, &totalQty)

		fmt.Println(totalAmount)
		fmt.Println(totalQty)
		paymentType.TotalAmount = totalAmount
		paymentType.TotalQty = totalQty

		Sold = append(Sold, paymentType)
	}

	return Sold
}
