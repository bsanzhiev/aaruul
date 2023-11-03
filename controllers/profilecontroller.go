package controllers

import (
	"net/http"

	"github.com/bsanzhiev/tsurhai/auth"
	"github.com/bsanzhiev/tsurhai/database"
	"github.com/bsanzhiev/tsurhai/models"
	"github.com/gin-gonic/gin"
)

func ProfileUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	email := claims.Email

	var user models.User
	if err := database.Connect().Debug().Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
