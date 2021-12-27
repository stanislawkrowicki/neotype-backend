package results

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"neotype-backend/pkg/mysql"
	"strconv"
)

type Result struct {
	gorm.Model

	User int     `json:"user"`
	WPM  float32 `json:"wpm"`
	Time int     `json:"time"`
}

var db = mysql.NewConnection()

func InitConsumer() {
	if err := db.AutoMigrate(&Result{}); err != nil {
		log.Fatal("Failed to migrate Result struct.")
	}
}

func ConsumeResult(body []byte) {
	var obj QueueObject

	err := json.Unmarshal(body, &obj)
	if err != nil {
		log.Printf("Failed to unmarshal result! %s", err)
		return
	}

	data := obj.Body
	var result Result

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Printf("Failed to unmarshal object body! %s", err)
		return
	}

	userString, ok := obj.User.(string)
	if !ok {
		log.Printf("Failed to get user ID string")
		return
	}

	userID, err := strconv.Atoi(userString)
	if err != nil {
		log.Printf("Failed to convert user ID string to int")
		return
	}

	result.User = userID

	resp := db.Create(&result)
	if resp.Error != nil {
		log.Printf("Failed to add the result to the database! %s", resp.Error)
	}
}
