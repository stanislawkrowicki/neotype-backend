package results

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"log"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/rabbitmq"
	"neotype-backend/pkg/users"
	"net/http"
	"strconv"
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

	callLeaderboards(c, jsonData)
}

func FetchResults(c *gin.Context) {
	const (
		defaultLimit = 10
		minLimit     = 1
		maxLimit     = 50
	)

	limit := defaultLimit

	limitInt, err := strconv.Atoi(c.Param("count"))
	if err == nil && limitInt <= maxLimit && limitInt >= minLimit {
		limit = limitInt
	}

	userID, err := users.Authorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var results []Result
	db.Limit(limit).Order("id desc").Find(&results, "user = ?", userID)

	c.JSON(http.StatusOK, results)
}

func callLeaderboards(c *gin.Context, body []byte) {
	serviceURI, err := config.GetBaseURL("leaderboards")
	if err != nil {
		log.Println("Could not get leaderboards URI!")
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", serviceURI+"/leaderboards", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create new http request to leaderboards: %s", err)
		return
	}

	auth, ok := c.Request.Header["Authorization"]
	if !ok {
		return
	}

	req.Header.Set("User-Agent", "results service")
	req.Header.Set("Authorization", auth[0])

	response, err := client.Do(req)
	if err != nil {
		log.Printf("Got error while calling leaderboards. Is the service dead?")
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
}
