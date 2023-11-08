package controllers

import (
	"context"
	"errors"
	"net/http"

	"time"

	"gorm.io/gorm"

	"github.com/bsanzhiev/tsurhai/auth"
	"github.com/bsanzhiev/tsurhai/database"
	"github.com/bsanzhiev/tsurhai/firebaseapp"
	"github.com/bsanzhiev/tsurhai/models"
	"github.com/gin-gonic/gin"
)

// RegisterUser
// ищем существующего пользователя
// если существует - проверяем что не мягко удален
// если он есть, но мягко удален - обновляем данные
// если его нет, создаем нового
func RegisterUser(c *gin.Context) {
	type RegisterResponse struct {
		UserID   uint   `json:"userID"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	var params models.RegisterRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var existingUser models.User
	err := database.Connect().Debug().Unscoped().Where("email = ?", params.Email).First(&existingUser).Error
	// Юзер найден
	if err == nil {
		// Если мягко удален, обновляем его данные
		if existingUser.DeletedAt.Valid {

			if len(params.Password) < 9 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be longer than 9 char."})
				return
			}

			if err := existingUser.HashPassword(params.Password); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}

			existingUser.FirstName = params.FirstName
			existingUser.SecondName = params.SecondName
			existingUser.Username = params.Username
			existingUser.DeletedAt = gorm.DeletedAt{}
			existingUser.CreatedAt = time.Now()

			// TODO Ебаная путаница с Instance и Connect - разобраться
			if updatesErr := database.Instance.Debug().Save(&existingUser).Error; updatesErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": updatesErr.Error()})
				return
			}

			c.JSON(http.StatusOK, RegisterResponse{
				UserID:   existingUser.ID,
				Email:    existingUser.Email,
				Username: existingUser.Username,
			})
		} else {
			// Иначе значит такой юзер уже создан
			c.JSON(http.StatusConflict, gin.H{"error": "user already exist"})
			c.Abort()
			return
		}

		// Если юзер не найден, создаем нового
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		var user models.User

		user.FirstName = params.FirstName
		user.SecondName = params.SecondName
		user.Username = params.Username
		user.Email = params.Email

		// проверяем пароль
		if len(params.Password) < 9 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be longer than 9 characters."})
			return
		}

		// хэшируем пароль
		if err := user.HashPassword(params.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		createErr := database.Connect().Debug().Create(&user).Error
		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": createErr.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusCreated, RegisterResponse{
			UserID:   user.ID,
			Email:    user.Email,
			Username: user.Username,
		})

		// Прочие ошибки
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
}

func LoginUser(c *gin.Context) {
	type LoginResponse struct {
		UserID      uint   `json:"userID"`
		Email       string `json:"email"`
		Username    string `json:"username"`
		AccessToken string `json:"token"`
	}
	var params models.LoginRequest
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

	c.JSON(http.StatusOK, LoginResponse{
		UserID:      user.ID,
		Email:       user.Email,
		Username:    user.Username,
		AccessToken: token,
	})
}

func VerifyToken(c *gin.Context) {
	type VerifyResponse struct {
		UserID      uint   `json:"userID"`
		Email       string `json:"email"`
		Username    string `json:"username"`
		AccessToken string `json:"token"`
	}

	type verifyResponse struct {
		IDToken string `json:"idToken"`
		Phone   string `json:"phone"`
	}

	var params verifyResponse
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Верификация токена аутентификации
	client, err := firebaseapp.FirebaseApp.Auth(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании клиента Firebase Auth"})
		return
	}

	_, err = client.VerifyIDToken(context.Background(), params.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен аутентификации"})
		c.Abort()
		return
	}

	// Токен верифицирован успешно, вы можете получить информацию о пользователе
	var user models.User
	if err := database.Connect().Debug().Where("phone = ?", params.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		c.Abort()
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create token."})
		return
	}

	c.JSON(http.StatusOK, VerifyResponse{
		UserID:      user.ID,
		Email:       user.Email,
		Username:    user.Username,
		AccessToken: token,
	})

}
