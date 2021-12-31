package leaderboards

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"neotype-backend/pkg/mysql"
	"neotype-backend/pkg/redis"
	"neotype-backend/pkg/results"
	"neotype-backend/pkg/users"
	"net/http"
	"strconv"
	"time"
)

var (
	rdb     = redis.NewConnection()
	ctx     = context.Background()
	mysqldb = mysql.NewConnection()
)

const (
	redisListKey        = "leaderboard"
	redisLeaderKey      = "leader"
	redisListExpiration = time.Hour * 72
)

func Entry(c *gin.Context) {
	var result results.Result

	err := c.ShouldBindJSON(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	userIDInterface, err := users.Authorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	userID := userIDInterface.(string)

	exists, err := rdb.Exists(ctx, redisListKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	var addError error

	score, err := rdb.ZScore(ctx, redisListKey, userID).Result()
	if err == goredis.Nil {
		addError = rdb.ZAdd(ctx, redisListKey, &goredis.Z{Member: userIDInterface, Score: float64(result.WPM)}).Err()
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	} else if score <= float64(result.WPM) {
		addError = rdb.ZAdd(ctx, redisListKey, &goredis.Z{Member: userIDInterface, Score: float64(result.WPM)}).Err()
	} else {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if addError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if exists == 0 {
		rdb.Expire(ctx, redisListKey, redisListExpiration)
	}

	var user users.User
	mysqldb.First(&user, "id = ?", userIDInterface)

	leaderKey := fmt.Sprintf("%s:%s", redisLeaderKey, userID)
	rdb.HSet(ctx, leaderKey,
		"accuracy", result.Accuracy, "username", user.Login, "time", result.Time)
	rdb.Expire(ctx, leaderKey, redisListExpiration)

	c.JSON(http.StatusOK, gin.H{})
}

func Leaders(c *gin.Context) {
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid number of leaders to fetch."})
		return
	}

	cardinals := int(rdb.ZCard(ctx, redisListKey).Val())
	if cardinals < count {
		count = cardinals
	}

	leaders, err := rdb.ZRevRangeWithScores(ctx, redisListKey, 0, int64(count)).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	var fetchedData []gin.H

	for _, leader := range leaders {
		hash := rdb.HGetAll(ctx, fmt.Sprintf("%s:%s", redisLeaderKey, leader.Member)).Val()
		fetchedData = append(fetchedData, gin.H{
			"username": hash["username"],
			"wpm":      leader.Score,
			"accuracy": hash["accuracy"],
			"time":     hash["time"],
		})
	}

	c.JSON(http.StatusOK, fetchedData)
}
