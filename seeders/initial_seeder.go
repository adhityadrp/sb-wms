package seeders

import (
	"log"
	"sb-wms/config"
	"sb-wms/models"

	"golang.org/x/crypto/bcrypt"
)

func Run() {
	log.Println("🌱 Menjalankan seeder awal...")

	// Seeder user admin
	pw, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "admin",
		Password: string(pw),
		Role:     "admin",
	}
	config.DB.FirstOrCreate(&admin, models.User{Username: "admin"})

	// Seeder kategori
	category := models.Category{Name: "Electronics"}
	config.DB.FirstOrCreate(&category, models.Category{Name: category.Name})

	// Seeder item
	item := models.Item{
		Name:       "Keyboard",
		SKU:        "KB-001",
		CategoryID: category.ID,
		Qty:        10,
		Price:      150000,
	}
	config.DB.FirstOrCreate(&item, models.Item{SKU: item.SKU})

	log.Println("✅ Seeder selesai dijalankan tanpa duplikasi.")
}
