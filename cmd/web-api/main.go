package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/gateway"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "UPDATE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// words
	router.GET("/words/:count", func(c *gin.Context) {
		gateway.Proxy(c, "words", fmt.Sprintf("/words/%s", c.Param("count")))
	})

	// users
	router.POST("/login", func(c *gin.Context) {
		gateway.Proxy(c, "users", "/login")
	})
	router.POST("/register", func(c *gin.Context) {
		gateway.Proxy(c, "users", "/register")
	})
	router.GET("/data", func(c *gin.Context) {
		gateway.Proxy(c, "users", "/data")
	})
	router.GET("/username", func(c *gin.Context) {
		gateway.Proxy(c, "users", "/username")
	})

	// results
	router.POST("/result", func(c *gin.Context) {
		gateway.Proxy(c, "results", "/result")
	})
	router.GET("/results/:count", func(c *gin.Context) {
		gateway.Proxy(c, "results", fmt.Sprintf("/results/%s", c.Param("count")))
	})

	// leaderboards
	router.GET("/leaderboards/:count", func(c *gin.Context) {
		gateway.Proxy(c, "leaderboards", fmt.Sprintf("/leaderboards/%s", c.Param("count")))
	})
	port, err := config.Get("web-api", "port")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
