package results

import (
	"encoding/json"
	"gorm.io/datatypes"
	"log"
	"neotype-backend/pkg/mysql"
	"neotype-backend/pkg/users"
	"strconv"
	"time"
)

type Result struct {
	ID        int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	User      int            `json:"user"`
	WPM       float32        `json:"wpm"`
	Accuracy  float32        `json:"accuracy"`
	Time      int            `json:"time"`
	CreatedAt datatypes.Date `json:"date"`
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
	result.CreatedAt = datatypes.Date(time.Now())

	resp := db.Create(&result)
	if resp.Error != nil {
		log.Printf("Failed to add the result to the database! %s", resp.Error)
	}

	var user users.User
	db.First(&user, "id = ?", userID)
	user.TestsTaken++
	user.AllTimeAvg = ((user.AllTimeAvg * float32(user.TestsTaken-1)) + result.WPM) / float32(user.TestsTaken)
	db.Save(&user)
}