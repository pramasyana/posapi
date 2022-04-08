package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/controllers"
	"github.com/jinzhu/gorm"
)

func Register(router fiber.Router, database *gorm.DB) {

	cashier := new(controllers.Cashier)
	routeCashier := router.Group("/cashiers")
	cashier.Mount(routeCashier)

	category := new(controllers.Category)
	routeCategory := router.Group("/categories")
	category.Mount(routeCategory)

	product := new(controllers.Product)
	routeProduct := router.Group("/products")
	product.Mount(routeProduct)

	payment := new(controllers.Payment)
	routePayment := router.Group("/payments")
	payment.Mount(routePayment)

	order := new(controllers.Order)
	routeOrder := router.Group("/orders")
	order.Mount(routeOrder)

	revenue := new(controllers.Revenue)
	routeRevenue := router.Group("/revenues")
	revenue.Mount(routeRevenue)

	sold := new(controllers.Sold)
	routeSold := router.Group("/solds")
	sold.Mount(routeSold)
}
