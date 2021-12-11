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

	port, err := config.Get("users", "port")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
