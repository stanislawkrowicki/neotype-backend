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

var wordsArr []string

func loadWords() error {
	mainPath, _ := os.Getwd()
	var jsonFile *os.File
	var err error
	if !strings.Contains(mainPath, "tests") {
		jsonFile, err = os.Open(mainPath + "/pkg/words/words.json")
	} else {
		jsonFile, err = os.Open(mainPath + "/../../pkg/words/words.json")
	}

	if err != nil {
		return err
	}

	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &wordsArr)
	if err != nil {
		return err
	}

	return nil
}

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

	if len(wordsArr) == 0 {
		if err := loadWords(); err != nil {
			c.JSON(http.StatusInternalServerError, "Failed to load words")
		}
	}

	var selectedWords []string

	randomizerSource := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(randomizerSource)

	for i := 0; i < count; i++ {
		index := randomizer.Intn(len(wordsArr))
		selectedWords = append(selectedWords, wordsArr[index])
	}

	c.JSON(http.StatusOK, selectedWords)
}
