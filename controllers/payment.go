package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Payment struct {
	Base
}

func (b *Payment) Mount(group fiber.Router) {
	group.Get("", b.GetAllPayment)
	group.Get("/:id", b.GetPayment)
	group.Post("", b.CreatePayment)
	group.Put("/:id", b.UpdatePayment)
	group.Delete("/:id", b.DeletePayment)
}

func (b *Payment) GetAllPayment(c *fiber.Ctx) error {
	count := models.GetPaymentsCount()
	var payments []models.Payments = models.FindAllPayment(c)

	var paylentList []models.PaymentsList
	for _, val := range payments {

		var nominal1, nominal2 int
		if len(c.Query("subtotal")) > 0 {
			nominal1, _ = strconv.Atoi(c.Query("subtotal"))
			nominal2 = (nominal1 + 10000) - (nominal1 % 10000)
		}

		payment := models.PaymentsList{
			PaymentId: val.PaymentId,
			Name:      val.Name,
			Type:      val.Type,
			Logo:      val.Logo,
		}

		if val.Type == "CASH" && len(c.Query("subtotal")) > 0 {
			payment.Card = []int{nominal1, nominal2}
		}
		paylentList = append(paylentList, payment)
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"payments": paylentList,
			"meta": fiber.Map{
				"total": count,
				"limit": limit,
				"skip":  skip,
			},
		},
	})
}

func (b *Payment) GetPayment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	payment, err := models.FindPayment(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    payment,
	})
}

func (b *Payment) CreatePayment(c *fiber.Ctx) error {
	data := new(models.Payments)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	var p struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Logo string `json:"logo"`
	}
	err := json.Unmarshal(c.Body(), &p)
	if err != nil || len(p.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"success": false,
			"message": "body ValidationError: \"name\" is required",
			"error": fiber.Map{
				"message": "\"name\" is required",
				"path":    "name",
				"type":    "any.required",
				"context": fiber.Map{
					"label": "name",
					"key":   "name",
				},
			},
		})
	}

	data.Name = p.Name
	data.Type = p.Type
	data.Logo = p.Logo

	item, err := models.CreatePayment(*data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed creating item",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    item,
	})
}

func (b *Payment) UpdatePayment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Payment Not Found",
			"error":   fiber.Map{},
		})
	}

	payment, err := models.FindPayment(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
			"error":   fiber.Map{},
		})
	}

	var p struct {
		Name string `json:"name"`
	}
	json.Unmarshal(c.Body(), &p)

	payment.Name = p.Name

	if len(payment.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	_, err = models.SavePayment(payment)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func (b *Payment) DeletePayment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting payment",
			"err":     err,
		})
	}

	RowsAffected := models.DeletePayment(id)
	if RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
