package models

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ListPaymentTypes struct {
	PaymentTypeID int    `json:"paymentTypeId"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	TotalAmount   int    `json:"totalAmount"`
}

func FindAllRevenue(c *fiber.Ctx) (int, []ListPaymentTypes) {
	var PaymentTypes []ListPaymentTypes
	var payments []Payments
	db := GetDB()
	db.Find(&payments)

	var total int
	for _, val := range payments {
		paymentType := ListPaymentTypes{
			PaymentTypeID: int(val.ID),
			Name:          val.Name,
			Logo:          val.Logo,
		}

		var totalAmount int
		db := GetDB().Table("orders").Select(
			"SUM(products.price * order_products.qty) AS totalAmount")
		db = db.Where("orders.payment_id = " + strconv.Itoa(int(val.PaymentId)))
		db = db.Joins("JOIN order_products on orders.id = order_products.orderid")
		db.Joins("INNER JOIN products on order_products.product_id = products.id").Row().Scan(&totalAmount)

		paymentType.TotalAmount = totalAmount
		total += totalAmount

		PaymentTypes = append(PaymentTypes, paymentType)
	}

	return total, PaymentTypes
}
