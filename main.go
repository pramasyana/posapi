package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ittechman101/go-pos/models"
	"github.com/ittechman101/go-pos/routes"
)

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Posapi",
	})
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "|${status}|${method} ${path} ${queryParams}| ${body} | ${resBody} |\n----------------------------------------\n",
	}))
	models.ConnectDB()
	//	defer database.DB.Close()
	routes.Register(app, models.DB)

	log.Fatal(app.Listen(":3030"))
}
