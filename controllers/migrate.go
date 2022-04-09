package controllers

import (
	"strconv"

	"github.com/ittechman101/go-pos/models"
)

type Migration struct {
	Base
}

func (b *Migration) Migrate() {
	go func() {
		count := models.GetCashierCount()
		if count == 0 {
			for i := 1; i <= 10; i++ {
				data := new(models.Cashiers)
				data.Name = "kasir " + strconv.Itoa(i)
				models.CreateCashier(*data)
			}

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

			for i := 1; i <= 10; i++ {
				data := new(models.Categories)
				data.Name = "Kategori " + strconv.Itoa(i)
				models.CreateCategory(*data)
			}

			for i := 1; i <= 5; i++ {
				data := new(models.Products)
				data.Name = "Produk " + strconv.Itoa(i)
				data.Stock = int64(i * 7)
				data.Price = float64(i * 78900)
				data.Image = "https://images.tokopedia.net/img/cache/500-square/hDjmkQ/2020/11/26/001f1c6e-d068-484f-9333-c3fa4129ef26.jpg"
				data.CategoriesId = int64(i)
				data.Discount = nil
				models.CreateProduct(*data)
			}
		}
	}()
}
