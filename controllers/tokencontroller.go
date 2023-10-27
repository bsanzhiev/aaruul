package controllers

import (
	"github.com/bsanzhiev/tsurhai/auth"
	"github.com/bsanzhiev/tsurhai/database"
	"github.com/bsanzhiev/tsurhai/models"

	"github.com/gin-gonic/gin"
	"net/http"
)

// TokenRequest
// Здесь мы определяем простую структуру, которая по сути будет тем,
// что конечная точка ожидает в качестве тела запроса.
// Он будет содержать идентификатор электронной почты и пароль пользователя.
type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GenerateToken
// Если все идет хорошо и пароль совпадает, генерируем JWT с помощью функции GenerateJWT().
// Это вернет подписанную строку токена со сроком действия 1 час,
// которая, в свою очередь, будет отправлена обратно клиенту в качестве ответа с кодом состояния 200.
func GenerateToken(c *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	// check if email exist and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}
	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
