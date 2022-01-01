package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/users"
)

func main() {
	users.Init()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "https://neotype.herokuapp.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/register", users.Register)
	router.POST("/login", users.Login)
	router.GET("/data", users.Data)
	router.GET("/username", users.Username)

	port, err := config.Get("users", "port")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
