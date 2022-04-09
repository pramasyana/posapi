package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	OrderId     int64   `gorm:"Not Null" json:"orderId"`
	TotalPaid   float64 `json:"totalPaid"`
	TotalPrice  float64 `json:"totalPrice"`
	TotalReturn float64 `json:"totalReturn"`
	ReceiptId   string  `json:"receiPtid"`
	CashierId   int64   `json:"cashierId"`
	PaymentId   int64   `json:"paymentId"`
	IsDownload  bool    `gorm:"default:false" json:"isDownload"`
}

type OrderProducts struct {
	gorm.Model
	ProductId        int64   `gorm:"Not Null" json:"productId"`
	Qty              int64   `json:"qty"`
	Orderid          int64   `json:"orderId"`
	Price            float64 `json:"price"`
	TotalNormalPrice float64 `json:"totalNormalPrice"`
	TotalFinalPrice  float64 `json:"totalFinalPrice"`
}

type ReqCreateOrder struct {
	PaymentId int64   `json:"paymentId"`
	TotalPaid float64 `json:"totalPaid"`
	Products  []struct {
		ProductId int64 `json:"productId"`
		Qty       int64 `json:"qty"`
	} `json:"products"`
}

type CreateOrderStruct struct {
	ID             int64     `json:"id"`
	CashiersId     int64     `json:"cashiersId"`
	PaymentTypesId int64     `json:"paymentTypesId"`
	TotalPrice     float64   `json:"totalPrice"`
	TotalPaid      float64   `json:"totalPaid"`
	TotalReturn    float64   `json:"totalReturn"`
	ReceiptId      string    `json:"receiptId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateProductStruct struct {
	ProductId        int64       `json:"productId"`
	Name             string      `json:"name"`
	Price            float64     `json:"price"`
	Discount         interface{} `json:"discount"`
	Qty              int64       `json:"qty"`
	TotalNormalPrice float64     `json:"totalNormalPrice"`
	TotalFinalPrice  float64     `json:"totalFinalPrice"`
}

type ResCreateOrder struct {
	CreatedOrder    CreateOrderStruct     `json:"order"`
	CreatedProducts []CreateProductStruct `json:"products"`
}

type OrderList struct {
	OrderId       int64     `json:"orderId"`
	CashiersId    int64     `json:"cashiersId"`
	PaymentTypeId int64     `json:"paymentTypesId"`
	TotalPrice    float64   `json:"totalPrice"`
	TotalPaid     float64   `json:"totalPaid"`
	TotalReturn   float64   `json:"totalReturn"`
	ReceiptId     string    `json:"receiptId"`
	CreatedAt     time.Time `json:"createdAt"`
	Cashier       struct {
		CashierId int64  `json:"cashierId"`
		Name      string `json:"name"`
	} `json:"cashier"`
	PaymentType struct {
		PaymentTypeId int64  `json:"paymentTypeId"`
		Name          string `json:"name"`
		Logo          string `json:"logo"`
		Type          string `json:"type"`
	} `json:"payment_type"`
}

type ReqDetailSubtotalOrder struct {
	ProductId int64 `json:"productid"`
	Qty       int64 `json:"qty"`
}

type ResDetailProduct struct {
	ProductId        int64       `json:"productId"`
	Name             string      `json:"name"`
	DiscountsId      interface{} `json:"discountsId"`
	Price            float64     `json:"price"`
	Discount         interface{} `json:"discount"`
	Qty              int64       `json:"qty"`
	TotalNormalPrice float64     `json:"totalNormalPrice"`
	TotalFinalPrice  float64     `json:"totalFinalPrice"`
}

type DetailOrder struct {
	OrderId       int64     `json:"orderId"`
	CashiersId    int64     `json:"cashiersId"`
	PaymentTypeId int64     `json:"paymentTypesId"`
	TotalPrice    float64   `json:"totalPrice"`
	TotalPaid     float64   `json:"totalPaid"`
	TotalReturn   float64   `json:"totalReturn"`
	ReceiptId     string    `json:"receiptId"`
	CreatedAt     time.Time `json:"createdAt"`
	Cashier       struct {
		CashierId int64  `json:"cashierId"`
		Name      string `json:"name"`
	} `json:"cashier"`
	PaymentType struct {
		PaymentTypeId int64  `json:"paymentTypeId"`
		Name          string `json:"name"`
		Logo          string `json:"logo"`
		Type          string `json:"type"`
	} `json:"payment_type"`
}

type ResOrderDetail struct {
	DetailOrder    DetailOrder        `json:"order"`
	DetailProducts []ResDetailProduct `json:"products"`
}

type ResSubTotalProduct struct {
	ProductId        int64       `json:"productId"`
	Name             string      `json:"name"`
	Stock            int64       `json:"stock"`
	Price            float64     `json:"price"`
	Image            string      `json:"image"`
	Sku              string      `json:"sku"`
	Qty              int64       `json:"qty"`
	CategoryId       int64       `json:"categoryId"`
	Discount         interface{} `json:"discount"`
	TotalNormalPrice float64     `json:"totalNormalPrice"`
	TotalFinalPrice  float64     `json:"totalFinalPrice"`
}
type ResSubTotal struct {
	SubTotal    float64              `json:"subtotal"`
	SubProducts []ResSubTotalProduct `json:"products"`
}

func FindOrderProductsByOrderId(id int64) []OrderProducts {
	var orderProducts []OrderProducts

	db := GetDB()
	err := GetDB().Where("orderid = ?", id).First(&orderProducts).Error
	if err != nil {
		return orderProducts
	}

	db.Where("orderid = ?", id).Find(&orderProducts)

	return orderProducts
}

type Discount struct {
	Qty       int64  `json:"qty"`
	Type      string `json:"type"`
	Result    int64  `json:"result"`
	ExpiredAt int64  `json:"expiredAt"`
}

func FindSubTotal(reqSubtotalOrder []ReqDetailSubtotalOrder) ResSubTotal {
	var resSubTotal ResSubTotal
	resSubTotal.SubProducts = make([]ResSubTotalProduct, len(reqSubtotalOrder))

	var subtotal float64 = 0
	for i := 0; i < len(reqSubtotalOrder); i++ {
		products, _ := FindProduct(int(reqSubtotalOrder[i].ProductId))
		resSubTotal.SubProducts[i].CategoryId = products.CategoriesId

		respDiscount := (*json.RawMessage)(&products.Discount)
		if respDiscount != nil {
			json.Unmarshal([]byte(*respDiscount), &resSubTotal.SubProducts[i].Discount)
		}

		resSubTotal.SubProducts[i].Image = products.Image
		resSubTotal.SubProducts[i].Name = products.Name
		resSubTotal.SubProducts[i].Price = products.Price
		resSubTotal.SubProducts[i].ProductId = products.ProductId
		resSubTotal.SubProducts[i].Sku = products.Sku
		resSubTotal.SubProducts[i].Stock = products.Stock
		resSubTotal.SubProducts[i].Qty = reqSubtotalOrder[i].Qty
		resSubTotal.SubProducts[i].TotalNormalPrice = products.Price * float64(reqSubtotalOrder[i].Qty)

		result := resSubTotal.SubProducts[i].TotalNormalPrice
		if resSubTotal.SubProducts[i].Discount != nil {
			discount := resSubTotal.SubProducts[i].Discount.(map[string]interface{})
			now := time.Now()
			exp := time.Unix(int64(discount["expiredAt"].(float64)), 0)
			fmt.Println(now.AddDate(1, 0, 0).Unix())
			formatString := "---"
			if now.Before(exp) {
				if discount["type"] == "PERCENT" {
					if resSubTotal.SubProducts[i].Qty >= int64(discount["qty"].(float64)) {
						result = resSubTotal.SubProducts[i].TotalNormalPrice - (resSubTotal.SubProducts[i].TotalNormalPrice / 100 * discount["result"].(float64))
						formatString = "Discount " + strconv.Itoa(int(discount["result"].(float64))) + "% Rp. " + strconv.Itoa(int(result))
					}

				} else if discount["type"] == "BUY_N" {
					if resSubTotal.SubProducts[i].Qty >= int64(discount["qty"].(float64)) {
						result = resSubTotal.SubProducts[i].TotalNormalPrice - discount["result"].(float64)
						formatString = "Buy " + strconv.Itoa(int(discount["qty"].(float64))) + " only Rp. " + strconv.Itoa(int(result))
					}
				}
			}

			resSubTotal.SubProducts[i].Discount = map[string]interface{}{
				"discountId":      products.ProductId,
				"qty":             discount["qty"],
				"type":            discount["type"],
				"result":          discount["result"],
				"expiredAt":       exp,
				"expiredAtFormat": exp.Format("02 Jan 2006"),
				"stringFormat":    formatString,
			}
		}

		resSubTotal.SubProducts[i].TotalFinalPrice = result
		subtotal += resSubTotal.SubProducts[i].TotalFinalPrice
	}
	resSubTotal.SubTotal = subtotal
	return resSubTotal
}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func GenerateReceipt() string {
	return "S" + strconv.Itoa(rangeIn(100, 999)) + RandStringBytes(1)
}

func FindAllOrders(c *fiber.Ctx) []OrderList {
	// var orders []Order
	var listOrder []OrderList
	var dataOrder OrderList

	db := GetDB().Table("orders")
	db = db.Select(`orders.order_id, orders.cashier_id, orders.payment_id, orders.total_paid, 
				orders.total_price, orders.total_return, orders.receipt_id, orders.created_at, 
				cashiers.cashier_id as cashier_cashier_id, cashiers.name as cashier_name, 
				payments.payment_id as payment_payment_id, payments.name as payment_name, payments.logo, payments.type`)
	db = db.Joins("JOIN cashiers ON orders.cashier_id = cashiers.cashier_id")
	db = db.Joins("JOIN payments ON orders.payment_id = payments.payment_id")

	limit := 10
	if len(c.Query("limit")) > 0 {
		limit, _ = strconv.Atoi(c.Query("limit"))
	}

	skip := 0
	if len(c.Query("skip")) > 0 {
		skip, _ = strconv.Atoi(c.Query("skip"))
	}

	db = db.Limit(limit)
	db = db.Offset(skip)

	rows, err := db.Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&dataOrder.OrderId,
			&dataOrder.CashiersId,
			&dataOrder.PaymentTypeId,
			&dataOrder.TotalPaid,
			&dataOrder.TotalPrice,
			&dataOrder.TotalReturn,
			&dataOrder.ReceiptId,
			&dataOrder.CreatedAt,
			&dataOrder.Cashier.CashierId,
			&dataOrder.Cashier.Name,
			&dataOrder.PaymentType.PaymentTypeId,
			&dataOrder.PaymentType.Name,
			&dataOrder.PaymentType.Logo,
			&dataOrder.PaymentType.Type,
		)

		listOrder = append(listOrder, dataOrder)
	}

	return listOrder
}

func GetOrderCount(c *fiber.Ctx) int64 {
	var orders []Order

	db := GetDB()

	count := db.Find(&orders).RowsAffected

	return count
}

func FindOrder(id int, reqDetailOrder []ReqDetailSubtotalOrder) (ResOrderDetail, error) {
	var order Order

	err := GetDB().Where("order_id = ?", id).First(&order).Error
	if err != nil {
		err = errors.New("order not found")
	}

	orderProducts := FindOrderProductsByOrderId(int64(id))
	var resOrderDetail ResOrderDetail

	resOrderDetail.DetailOrder.OrderId = order.OrderId
	resOrderDetail.DetailOrder.CashiersId = order.CashierId
	resOrderDetail.DetailOrder.PaymentTypeId = order.PaymentId
	resOrderDetail.DetailOrder.TotalPaid = order.TotalPaid
	resOrderDetail.DetailOrder.CreatedAt = order.CreatedAt

	resOrderDetail.DetailProducts = make([]ResDetailProduct, len(orderProducts))
	for i := 0; i < len(orderProducts); i++ {
		products, getProductErr := FindProduct(int(orderProducts[i].ProductId))
		if getProductErr == nil {
			resOrderDetail.DetailProducts[i].ProductId = products.ProductId
			resOrderDetail.DetailProducts[i].Name = products.Name
			resOrderDetail.DetailProducts[i].Price = products.Price
			resOrderDetail.DetailProducts[i].DiscountsId = nil
			if products.Discount != nil {
				resOrderDetail.DetailProducts[i].DiscountsId = products.ProductId
			}

			respDiscount := (*json.RawMessage)(&products.Discount)
			if respDiscount != nil {
				json.Unmarshal([]byte(*respDiscount), &resOrderDetail.DetailProducts[i].Discount)
			}

			resOrderDetail.DetailProducts[i].Qty = orderProducts[i].Qty
			resOrderDetail.DetailProducts[i].TotalNormalPrice = orderProducts[i].TotalNormalPrice

			// result := int(resOrderDetail.DetailProducts[i].TotalNormalPrice)
			// if resOrderDetail.DetailProducts[i].Discount != nil {
			// 	discount := resOrderDetail.DetailProducts[i].Discount.(map[string]interface{})
			// 	now := time.Now()
			// 	exp := time.Unix(int64(discount["expiredAt"].(float64)), 0)
			// 	if now.Before(exp) {
			// 		if discount["type"] == "PERCENT" {
			// 			if resOrderDetail.DetailProducts[i].Qty >= int64(discount["qty"].(float64)) {
			// 				result = int(resOrderDetail.DetailProducts[i].TotalNormalPrice) - (int(resOrderDetail.DetailProducts[i].TotalNormalPrice) / 100 * int(discount["result"].(float64)))
			// 			}
			// 		} else if discount["type"] == "BUY_N" {
			// 			if resOrderDetail.DetailProducts[i].Qty >= int64(discount["qty"].(float64)) {
			// 				result = int(resOrderDetail.DetailProducts[i].TotalNormalPrice) - int(discount["result"].(float64))
			// 			}
			// 		}
			// 	}
			// }
			resOrderDetail.DetailProducts[i].TotalFinalPrice = orderProducts[i].TotalFinalPrice
			resOrderDetail.DetailOrder.TotalPrice += orderProducts[i].TotalFinalPrice
		}
	}
	resOrderDetail.DetailOrder.TotalReturn = resOrderDetail.DetailOrder.TotalPaid - resOrderDetail.DetailOrder.TotalPrice
	resOrderDetail.DetailOrder.ReceiptId = order.ReceiptId
	cashier, _ := FindCashier(int(order.CashierId))
	payment, _ := FindPayment(int(order.PaymentId))
	resOrderDetail.DetailOrder.Cashier.CashierId = cashier.CashierId
	resOrderDetail.DetailOrder.Cashier.Name = cashier.Name
	resOrderDetail.DetailOrder.PaymentType.PaymentTypeId = payment.PaymentId
	resOrderDetail.DetailOrder.PaymentType.Type = payment.Type
	resOrderDetail.DetailOrder.PaymentType.Name = payment.Name
	resOrderDetail.DetailOrder.PaymentType.Logo = payment.Logo

	return resOrderDetail, err
}

func CreateOrder(reqOrder ReqCreateOrder) (ResCreateOrder, error) {
	var payment Payments
	order := new(Order)
	var orderProducts = make([]OrderProducts, len(reqOrder.Products))

	// var id int
	// GetDB().Raw(`
	// 	SELECT AUTO_INCREMENT as order_id FROM information_schema.TABLES WHERE TABLE_SCHEMA = "` + config.Config("DB_NAME") + `" AND TABLE_NAME = "orders"
	// 	`).Row().Scan(
	// 	&id,
	// )

	var maxOrder Order

	GetDB().Raw(`
		SELECT COALESCE(MAX(order_id) + 1, 1) as order_id
		FROM orders
		`).Scan(
		&maxOrder,
	)

	order.CashierId = 1
	order.PaymentId = reqOrder.PaymentId
	order.TotalPaid = reqOrder.TotalPaid
	order.OrderId = maxOrder.OrderId
	order.ReceiptId = GenerateReceipt()

	var resCreateOrder ResCreateOrder
	resCreateOrder.CreatedOrder.ID = maxOrder.OrderId
	resCreateOrder.CreatedOrder.PaymentTypesId = order.PaymentId
	resCreateOrder.CreatedOrder.CashiersId = order.CashierId
	resCreateOrder.CreatedOrder.TotalPaid = order.TotalPaid
	resCreateOrder.CreatedOrder.CreatedAt = time.Now()
	resCreateOrder.CreatedOrder.UpdatedAt = time.Now()
	resCreateOrder.CreatedOrder.ReceiptId = order.ReceiptId

	count := GetDB().Table("payments").Where("payment_id = ?", reqOrder.PaymentId).Find(&payment).RowsAffected
	if count == 0 {
		return resCreateOrder, errors.New("PaymentType not found")
	}

	resCreateOrder.CreatedProducts = make([]CreateProductStruct, len(reqOrder.Products))
	for i := 0; i < len(reqOrder.Products); i++ {
		products, getProductErr := FindProduct(int(reqOrder.Products[i].ProductId))
		if getProductErr == nil {
			resCreateOrder.CreatedProducts[i].ProductId = products.ProductId
			resCreateOrder.CreatedProducts[i].Name = products.Name
			resCreateOrder.CreatedProducts[i].Price = products.Price

			respDiscount := (*json.RawMessage)(&products.Discount)
			if respDiscount != nil {
				json.Unmarshal([]byte(*respDiscount), &resCreateOrder.CreatedProducts[i].Discount)
			}
			resCreateOrder.CreatedProducts[i].Qty = reqOrder.Products[i].Qty
			resCreateOrder.CreatedProducts[i].TotalNormalPrice = products.Price * float64(reqOrder.Products[i].Qty)

			result := float64(resCreateOrder.CreatedProducts[i].TotalNormalPrice)
			if resCreateOrder.CreatedProducts[i].Discount != nil {
				discount := resCreateOrder.CreatedProducts[i].Discount.(map[string]interface{})
				now := time.Now()
				exp := time.Unix(int64(discount["expiredAt"].(float64)), 0)
				if now.Before(exp) {
					if discount["type"] == "PERCENT" {
						if resCreateOrder.CreatedProducts[i].Qty >= int64(discount["qty"].(float64)) {
							result = float64(resCreateOrder.CreatedProducts[i].TotalNormalPrice) - (float64(resCreateOrder.CreatedProducts[i].TotalNormalPrice) / 100 * discount["result"].(float64))
						}
					} else if discount["type"] == "BUY_N" {
						if resCreateOrder.CreatedProducts[i].Qty >= int64(discount["qty"].(float64)) {
							result = float64(resCreateOrder.CreatedProducts[i].TotalNormalPrice) - discount["result"].(float64)
						}
					}
				}
			}
			resCreateOrder.CreatedProducts[i].TotalFinalPrice = result
			resCreateOrder.CreatedOrder.TotalPrice += resCreateOrder.CreatedProducts[i].TotalFinalPrice
		}

		orderProducts[i].ProductId = reqOrder.Products[i].ProductId
		orderProducts[i].Qty = reqOrder.Products[i].Qty
		orderProducts[i].Orderid = order.OrderId
		orderProducts[i].Price = resCreateOrder.CreatedProducts[i].Price
		orderProducts[i].TotalNormalPrice = resCreateOrder.CreatedProducts[i].TotalNormalPrice
		orderProducts[i].TotalFinalPrice = resCreateOrder.CreatedProducts[i].TotalFinalPrice
		orderProductErr := GetDB().Create(&orderProducts[i]).Error
		if orderProductErr != nil {
			return resCreateOrder, orderProductErr
		}
	}
	resCreateOrder.CreatedOrder.TotalReturn = resCreateOrder.CreatedOrder.TotalPaid - resCreateOrder.CreatedOrder.TotalPrice

	order.TotalPrice = resCreateOrder.CreatedOrder.TotalPrice
	order.TotalReturn = resCreateOrder.CreatedOrder.TotalReturn
	orderErr := GetDB().Create(order).Error
	if orderErr != nil {
		return resCreateOrder, orderErr
	}
	return resCreateOrder, nil
}

func FindOrderDownload(id int) (Order, error) {
	var order Order

	err := GetDB().Where("order_id = ?", id).First(&order).Error
	if err != nil {
		err = errors.New("order not found")
	}

	return order, err
}

func UpdateIsDownload(id int) error {
	err := GetDB().Table("orders").Where("order_id = ?", id).Update("is_download", true).Error

	return err
}
