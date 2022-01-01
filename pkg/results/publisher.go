package results

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"io"
	"log"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/rabbitmq"
	"neotype-backend/pkg/users"
	"net/http"
	"strconv"
)

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
	var result Result

	err := c.ShouldBindJSON(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unreadable request body."})
		return
	}

	if result.WPM == 0 || result.Time == 0 || result.Accuracy == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request body is missing fields."})
		return
	}

	userID, err := users.ShouldAuthorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in."})
		return
	}

	result.User, _ = strconv.Atoi(userID.(string))

	body, err := json.Marshal(result)
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

	callLeaderboards(c, body)
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

	userID, err := users.ShouldAuthorize(c)
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
