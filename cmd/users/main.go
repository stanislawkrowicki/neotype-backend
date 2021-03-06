package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/users"
)

func main() {
	users.Init()

	router := gin.Default()

	router.POST("/register", users.Register)
	router.POST("/login", users.Login)
	router.GET("/data", users.Data)
	router.GET("/username", users.Username)
	router.GET("/authorize/:token", users.Authorize)

	port, err := config.GetPort("users")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
