package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/words"
)

func main() {
	router := gin.Default()

	router.GET("/words/:count", words.GetRandomWords)

	port, err := config.GetPort("words")
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
