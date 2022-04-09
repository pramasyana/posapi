package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Order struct {
	Base
}

func (b *Order) Mount(group fiber.Router) {
	group.Post("", b.CreateOrder)
	group.Post("/subtotal", b.GetSubTotal)
	group.Get("", b.GetAllOrder)
	group.Get("/:id", b.GetOrder)
	group.Get("/:id/download", b.GetOrderDownload)
	group.Get("/:id/check-download", b.GetOrderStatus)
}

func (b *Order) GetAllOrder(c *fiber.Ctx) error {
	orderList := models.FindAllOrders(c)

	count := models.GetOrderCount(c)
	if len(orderList) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "order Not Found",
			"error":   fiber.Map{},
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": fiber.Map{
				"orders": orderList,
				"meta": fiber.Map{
					"total": count,
					"limit": c.Query("limit"),
					"skip":  c.Query("skip"),
				},
			},
		})
	}
}

func (b *Order) GetSubTotal(c *fiber.Ctx) error {

	data := []models.ReqDetailSubtotalOrder{}

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	json.Unmarshal(c.Body(), &data)

	// respCh := make(chan interface{}, 1)
	// go func() {
	result := models.FindSubTotal(data)

	// 	if respCh != nil {
	// 		respCh <- resOrderSubTotal
	// 	}
	// }()

	// res := <-respCh
	// result := res.(models.ResSubTotal)

	if len(result.SubProducts) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Empty product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    result,
	})
}

func (b *Order) GetOrder(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var data = new([]models.ReqDetailSubtotalOrder)

	json.Unmarshal(c.Body(), &data)

	respCh := make(chan interface{}, 1)
	go func() {
		resOrderDetail, _ := models.FindOrder(id, *data)

		if respCh != nil {
			respCh <- resOrderDetail
		}
	}()

	result := <-respCh

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    result.(models.ResOrderDetail),
	})
}

func (b *Order) GetOrderDownload(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var data = new([]models.ReqDetailSubtotalOrder)

	json.Unmarshal(c.Body(), &data)

	resOrderDetail, _ := models.FindOrder(id, *data)

	models.UpdateIsDownload(id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    resOrderDetail,
	})
}

func (b *Order) GetOrderStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	resOrderDetail, err := models.FindOrderDownload(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found",
		})
	}
	var isDownload struct {
		IsDownload bool `json:"isDownload"`
	}

	isDownload.IsDownload = false
	if resOrderDetail.IsDownload {
		isDownload.IsDownload = true
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    isDownload,
	})
}

func (b *Order) CreateOrder(c *fiber.Ctx) error {
	data := new(models.ReqCreateOrder)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	err := json.Unmarshal(c.Body(), &data)
	if err != nil || len(data.Products) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create order",
			"error":   fiber.Map{},
		})
	}

	respCh := make(chan models.ResCreateOrder, 1)
	respErrorCh := make(chan error, 1)
	go func() {
		item, err := models.CreateOrder(*data)
		if respErrorCh != nil {
			respErrorCh <- err
		}

		if respCh != nil {
			respCh <- item
		}
	}()

	err = <-respErrorCh
	result := <-respCh

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create order : " + err.Error(),
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    result,
	})
}
