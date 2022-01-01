package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/leaderboards"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "https://neotype.herokuapp.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/leaderboards", leaderboards.Entry)
	router.GET("/leaderboards/:count", leaderboards.Leaders)

	port, err := config.Get("leaderboards", "port")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

}
