package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ittechman101/go-pos/models"

	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
)

type Cashier struct {
	Base
}

type JwtCustomClaims struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

type CashierResponse struct {
	CashierId int64  `json:"cashierId"`
	Name      string `json:"name"`
}

func (b *Cashier) Mount(group fiber.Router) {
	group.Post("", b.CreateCashier)
	group.Get("", b.GetAllCashier)
	group.Get("/:id", b.GetCashier)
	group.Put("/:id", b.UpdateCashier)
	group.Delete("/:id", b.DeleteCashier)

	group.Get("/:id/passcode", b.Passcode)
	group.Post("/:id/login", b.Login)
}

func (b *Cashier) GetAllCashier(c *fiber.Ctx) error {
	count := models.GetCashierCount()
	if count == 0 {
		for i := 1; i <= 10; i++ {
			data := new(models.Cashiers)
			data.Name = "kasir " + strconv.Itoa(i)
			models.CreateCashier(*data)
		}
	}

	count = models.GetCashierCount()
	var cashiers []models.Cashiers = models.FindAllCashier(c)
	cashiersResponses := make([]CashierResponse, len(cashiers))
	for i, element := range cashiers {
		cashiersResponses[i].CashierId = element.CashierId
		cashiersResponses[i].Name = element.Name
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"cashiers": cashiersResponses,
			"meta": fiber.Map{
				"total": count,
				"limit": limit,
				"skip":  skip,
			},
		},
	})
}

func (b *Cashier) GetCashier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	cashier, err := models.FindCashier(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashier,
	})
}

func (b *Cashier) Passcode(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	passCode, err := models.Passcode(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"passcode": passCode,
		},
	})
}

func (b *Cashier) Login(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "param validationError: \"cashierId\" is required"})
	}

	var p struct {
		Passcode string `json:"passcode"`
	}
	err = json.Unmarshal(c.Body(), &p)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "body validationError: \"passcode\" is required"})
	}

	passCode, err := models.Passcode(id)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
		})
	}

	if passCode == p.Passcode {
		//		secretKey := config.Config("SECRET_KEY")
		secretKey := "goPos"
		claims := jwt.MapClaims{
			"id":     id,
			"active": true,
			"exp":    time.Now().Add(time.Hour * 6).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "JWT Token Error",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": fiber.Map{
				"token": tokenString,
			},
		})
	} else {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Passcode Not Match",
		})
	}
}

func (b *Cashier) CreateCashier(c *fiber.Ctx) error {
	data := new(models.Cashiers)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	var p struct {
		Name string `json:"name"`
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

	item, err := models.CreateCashier(*data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed creating item",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    item,
	})
}

func (b *Cashier) UpdateCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	cashier, err := models.FindCashier(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	var p struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(c.Body(), &p)
	cashier.Name = p.Name
	if len(cashier.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	item, err := models.SaveCashier(cashier)
	item = item
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func (b *Cashier) DeleteCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting cashier",
			"err":     err,
		})
	}

	RowsAffected := models.DeleteCashier(id)
	if RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
