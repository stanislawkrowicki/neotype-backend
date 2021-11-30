package words

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetRandomWords(c *gin.Context) {
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	jsonFile, err := os.Open("pkg/words/words.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var allWords []string
	var selectedWords []string

	err = json.Unmarshal(byteValue, &allWords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	randomizerSource := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(randomizerSource)

	for i := 0; i < count; i++ {
		index := randomizer.Intn(len(allWords))
		selectedWords = append(selectedWords, allWords[index])
	}

	c.JSON(http.StatusOK, selectedWords)
}
