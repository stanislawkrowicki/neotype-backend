package web_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"neotype-backend/pkg/config"
	"net/http"
	"strconv"
)

func GetWords(c *gin.Context) {
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	wordsServiceURI, err := config.GetBaseURL("words")
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/words/%d", wordsServiceURI, count))
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}
