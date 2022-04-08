package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Revenue struct {
	Base
}

func (b *Revenue) Mount(group fiber.Router) {
	group.Get("", b.GetAllRevenue)
}

func (b *Revenue) GetAllRevenue(c *fiber.Ctx) error {
	var paymentTypes []models.ListPaymentTypes
	var totalRevenue int
	totalRevenue, paymentTypes = models.FindAllRevenue(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"totalRevenue": totalRevenue,
			"paymentTypes": paymentTypes,
		},
	})

}
