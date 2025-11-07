package controllers

import (
	"net/http"
	"sb-wms/config"
	"sb-wms/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var cats []models.Category
	if err := config.DB.Find(&cats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func CreateCategory(c *gin.Context) {
	var in models.Category
	if err := c.ShouldBindJSON(&in); err != nil || in.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	cat := models.Category{Name: in.Name}
	if err := config.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.Preload("Items").First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	var in struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	cat.Name = in.Name
	if err := config.DB.Save(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Cek apakah ada transaksi yang terkait dengan item di kategori ini
	var count int64
	if err := config.DB.Model(&models.Transaction{}).
		Joins("JOIN items ON transactions.item_id = items.id").
		Where("items.category_id = ?", id).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Terdapat transaksi terkait, tidak dapat menghapus kategori"})
		return
	}

	// Jika tidak ada transaksi, hapus kategori (cascade akan hapus items)
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}


