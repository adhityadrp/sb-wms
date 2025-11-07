package controllers

import (
	"net/http"
	"sb-wms/config"
	"sb-wms/models"

	// "strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTransaction(c *gin.Context) {
	var in struct {
		ItemID uint   `json:"item_id" binding:"required"`
		Qty    int    `json:"qty" binding:"required"`
		Type   string `json:"type" binding:"required"` // "in" or "out"
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// get user_id from context (set in JWT middleware)
	uid, _ := c.Get("user_id")
	userID := uid.(uint)

	// find item
	var item models.Item
	if err := config.DB.First(&item, in.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	if in.Type == "out" {
		if item.Qty < in.Qty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
			return
		}
		item.Qty -= in.Qty
	} else if in.Type == "in" {
		item.Qty += in.Qty
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'in' or 'out'"})
		return
	}

	// update item and create transaction
	tx := models.Transaction{
		UserID: userID,
		ItemID: in.ItemID,
		Qty:    in.Qty,
		Type:   in.Type,
		Note:   in.Note,
	}
	if err := config.DB.Transaction(func(txDB *gorm.DB) error {
		if err := txDB.Save(&item).Error; err != nil {
			return err
		}
		if err := txDB.Create(&tx).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}

	// 🔥 Setelah sukses create, preload relasi user & item
	if err := config.DB.
		Preload("User").
		Preload("Item").
		Preload("Item.Category").
		First(&tx, tx.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load transaction data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "transaction recorded",
		"transaction":  tx,
	})
}
