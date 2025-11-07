package main

import (
	"log"
	"os"
	"sb-wms/config"
	"sb-wms/middlewares"
	"sb-wms/models"
	"sb-wms/routes"
	"sb-wms/seeders"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.Init()
	middlewares.InitJWT()

	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Item{},
		&models.Transaction{},
	); err != nil {
		log.Fatal("❌ Migrate failed:", err)
	}
	
	// Run seeders
	seeders.Run()

	r := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Server running on port", port)
	r.Run(":" + port)
}
