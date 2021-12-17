package words

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetRandomWords(c *gin.Context) {
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.String(http.StatusBadRequest, "Not a number")
		return
	}

	if count <= 0 {
		c.String(http.StatusBadRequest, "Number should be greater than 0")
		return
	}

	mainPath, _ := os.Getwd()
	var jsonFile *os.File
	if !strings.Contains(mainPath, "tests") {
		jsonFile, err = os.Open(mainPath + "/pkg/words/words.json")
	} else {
		jsonFile, err = os.Open(mainPath + "/../../pkg/words/words.json")
	}

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
