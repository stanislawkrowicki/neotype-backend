package results

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"neotype-backend/pkg/rabbitmq"
	"neotype-backend/pkg/users"
	"net/http"
)

type QueueObject struct {
	User interface{} `json:"user"`
	Body []byte      `json:"result"`
}

const queueName = "results"

var (
	rabbit *rabbitmq.RabbitMQ
)

func InitPublisher() {
	rabbit = rabbitmq.New()
	err := rabbit.Connect(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %s", err)
	}
}

func QueueResult(c *gin.Context) {
	// TODO: some validation
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	userID, err := users.Authorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in."})
		return
	}

	obj := QueueObject{userID, jsonData}
	body, err := json.Marshal(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal request data."})
		return
	}

	err = rabbit.Channel.Publish(
		"",
		rabbit.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add the result to the queue."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added result to the queue."})
}

func FetchResults(c *gin.Context) {
	userID, err := users.Authorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var results []Result
	db.Limit(10).Order("id desc").Find(&results, "user = ?", userID)

	c.JSON(http.StatusOK, results)
}
