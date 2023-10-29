package controllers

import (
	"net/http"

	"github.com/bsanzhiev/tsurhai/auth"
	"github.com/bsanzhiev/tsurhai/database"
	"github.com/bsanzhiev/tsurhai/models"
	"github.com/gin-gonic/gin"
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

func LoginUser(c *gin.Context) {
	type ResponseData struct {
		UserID   uint   `json:"userID"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	var params models.UserLogin
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var user models.User
	if err := database.Connect().Debug().Where("email = ?", params.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		c.Abort()
		return
	}

	// if err := models.CheckPassword(params.Password); err != nil {
	credentialError := user.CheckPassword(params.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create token."})
		return
	}

	c.JSON(http.StatusOK, ResponseData{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Token:    token,
	})
}
