package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/leaderboards"
)

func main() {
	router := gin.Default()

	router.POST("/leaderboards", leaderboards.Entry)
	router.GET("/leaderboards/:count", leaderboards.Leaders)

	port, err := config.GetPort("leaderboards")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

}
