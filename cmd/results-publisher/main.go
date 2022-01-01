package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/results"
)

func main() {
	results.InitPublisher()

	router := gin.Default()

	router.POST("/result", results.QueueResult)
	router.GET("/results/:count", results.FetchResults)

	port, err := config.GetPort("results")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

}
