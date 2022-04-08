package controllers

import (
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
}

func (b *Payment) GetAllPayment(c *fiber.Ctx) error {
	count := models.GetPaymentsCount()
	if count == 0 {
		for i := 1; i <= 5; i++ {
			types := "E-WALLET"
			if i == 1 {
				types = "CASH"
			}
			data := new(models.Payments)
			data.Name = "Payment " + strconv.Itoa(i)
			data.Type = types
			data.Logo = "https://rm.id/images/berita/750x390/genjot-layanan-ovo-buka-peluang-kerja-sama-dengan-berbagai-pihak_22246.jpg"
			data.Card = "[" + strconv.Itoa(i*5000) + ", " + strconv.Itoa((i+1)*5000) + "]"
			models.CreatePayment(*data)
		}
	}

	count = models.GetPaymentsCount()
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

		if val.Type == "CASH" {
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
