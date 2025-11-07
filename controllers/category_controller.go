package controllers

import (
	"net/http"
	"sb-wms/config"
	"sb-wms/models"

	"github.com/gin-gonic/gin"
)

// GetCategories mengambil semua kategori dari database dan mengembalikannya sebagai JSON
func GetCategories(c *gin.Context) {
	var cats []models.Category
	if err := config.DB.Find(&cats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// CreateCategory membuat kategori baru berdasarkan input JSON dan menyimpannya ke database
func CreateCategory(c *gin.Context) {
	var in models.Category
	if err := c.ShouldBindJSON(&in); err != nil || in.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	cat := models.Category{Name: in.Name}
	if err := config.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed create"})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// GetCategoryByID mengambil kategori berdasarkan ID dari URL parameter, termasuk item-item yang terkait
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.Preload("Items").First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// DeleteCategory menghapus kategori berdasarkan ID dari URL parameter
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// UpdateCategory memperbarui kategori berdasarkan ID dari URL parameter dengan data baru dari JSON
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, cat)
}
