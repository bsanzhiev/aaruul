package controllers

import (
	"github.com/bsanzhiev/tsurhai/database"
	"github.com/bsanzhiev/tsurhai/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// abort?
		c.Abort()
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"userID": user.ID, "email": user.Email, "username": user.Username})
}
