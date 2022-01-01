package results

import (
	"encoding/json"
	"gorm.io/datatypes"
	"log"
	"neotype-backend/pkg/mysql"
	"neotype-backend/pkg/users"
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
	var result Result

	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Failed to unmarshal result! %s", err)
		return
	}

	result.CreatedAt = datatypes.Date(time.Now())

	resp := db.Create(&result)
	if resp.Error != nil {
		log.Printf("Failed to add the result to the database! %s", resp.Error)
	}

	var user users.User
	db.First(&user, "id = ?", result.User)
	user.TestsTaken++
	user.AllTimeAvg = ((user.AllTimeAvg * float32(user.TestsTaken-1)) + result.WPM) / float32(user.TestsTaken)
	db.Save(&user)
}
