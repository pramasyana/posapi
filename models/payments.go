package models

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Payments struct {
	gorm.Model
	PaymentId int64  `gorm:"Not Null" json:"paymentId"`
	Name      string `gorm:"Not Null" json:"name"`
	Type      string `gorm:"Not Null" json:"type"`
	Logo      string `gorm:"Not Null" json:"logo"`
	Card      string `gorm:"Not Null" json:"card"`
}

type PaymentsList struct {
	PaymentId int64  `json:"paymentId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Logo      string `json:"logo"`
	Card      []int  `json:"card,omitempty"`
}

func FindAllPayment(c *fiber.Ctx) []Payments {
	var payments []Payments

	db := GetDB()

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

	db.Find(&payments)

	return payments
}

func GetPaymentsCount() int64 {
	var payments []Payments

	count := GetDB().Find(&payments).RowsAffected

	return count
}

func FindPayment(id int) (Payments, error) {
	var payment Payments

	err := GetDB().Where("payment_id = ?", id).First(&payment).Error
	if err != nil {
		err = errors.New("payment not found")
	}

	return payment, err
}

func CreatePayment(payments Payments) (Payments, error) {
	// Get Max cashierId
	// var id int
	// GetDB().Raw(`
	// 	SELECT AUTO_INCREMENT as payment_id FROM information_schema.TABLES WHERE TABLE_SCHEMA = "` + config.Config("DB_NAME") + `" AND TABLE_NAME = "payments"
	// 	`).Row().Scan(
	// 	&id,
	// )

	var maxPayment Payments

	GetDB().Raw(`
		SELECT COALESCE(MAX(payment_id) + 1, 1) as payment_id
		FROM payments
		`).Scan(
		&maxPayment,
	)

	payments.PaymentId = maxPayment.PaymentId

	err := GetDB().Create(&payments).Error
	if err != nil {
		return payments, err
	}

	return payments, nil
}

func SavePayment(payment Payments) (Payments, error) {
	err := GetDB().Table("payments").Where("payment_id = ?", payment.PaymentId).Update(payment).Error

	return payment, err
}

func DeletePayment(id int) int64 {
	count := GetDB().Where("payment_id = ?", id).Delete(&Payments{}).RowsAffected

	return count
}
