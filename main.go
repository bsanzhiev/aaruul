package main

import (
	"fmt"
	"github.com/bsanzhiev/tsurhai/controllers"
	"github.com/bsanzhiev/tsurhai/database"
	"github.com/bsanzhiev/tsurhai/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		panic(viperErr)
	}
	gin.SetMode(viper.GetString("APP_MODE"))

	database.Connect()
	database.Migrate()

	//app := gin.Default()
	app := initRouter()
	app.GET("/ping", Ping)

	log.Fatalln(
		app.Run(
			fmt.Sprintf("0.0.0.0:%d", viper.GetInt("APP_PORT")),
		),
	)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/pong", controllers.Pong)
		}
	}
	return router
}
