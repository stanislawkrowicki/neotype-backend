package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	webapi "neotype-backend/pkg/web-api"
)

func main() {
	router := gin.Default()

	router.GET("/words/:count", webapi.GetWords)

	port, err := config.Get("web-api", "port")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
