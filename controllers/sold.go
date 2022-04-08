package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Sold struct {
	Base
}

func (b *Sold) Mount(group fiber.Router) {
	group.Get("", b.GetAllRevenue)
}

func (b *Sold) GetAllRevenue(c *fiber.Ctx) error {
	orderProducts := models.FindAllSold(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"orderProducts": orderProducts,
		},
	})

}
