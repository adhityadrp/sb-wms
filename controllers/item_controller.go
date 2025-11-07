package controllers

import (
	"net/http"
	"sb-wms/config"
	"sb-wms/models"

	"github.com/gin-gonic/gin"
)

func GetItems(c *gin.Context) {
	var items []models.Item
	if err := config.DB.Preload("Category").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func CreateItem(c *gin.Context) {
	var in models.Item
	if err := c.ShouldBindJSON(&in); err != nil || in.Name == "" || in.SKU == "" || in.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if in.Qty < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qty cannot be negative"})
		return
	}

	if err := config.DB.Create(&in).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Preload("Category").First(&in, in.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, in)
}

func GetItemByID(c *gin.Context) {
	id := c.Param("id")
	var it models.Item
	if err := config.DB.Preload("Category").First(&it, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, it)
}

func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var it models.Item
	if err := config.DB.First(&it, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	var in struct {
		Name       string `json:"name"`
		SKU        string `json:"sku"`
		CategoryID uint   `json:"category_id"`
		Qty        int    `json:"qty"`
		Price      int    `json:"price"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.Name == "" || in.SKU == "" || in.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if in.Qty < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qty cannot be negative"})
		return
	}

	it.Name = in.Name
	it.SKU = in.SKU
	it.CategoryID = in.CategoryID
	it.Qty = in.Qty
	it.Price = in.Price

	if err := config.DB.Save(&it).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updated models.Item
	if err := config.DB.Preload("Category").First(&updated, it.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")

	var count int64
	if err := config.DB.Model(&models.Transaction{}).Where("item_id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ada transaksi terkait, Item tidak bisa dihapus"})
		return
	}

	if err := config.DB.Delete(&models.Item{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}


