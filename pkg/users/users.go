package users

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"neotype-backend/pkg/mysql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type User struct {
	gorm.Model

	Login              string `json:"login" gorm:"unique"`
	Password           string
	TestsTaken         int     `json:"tests"`
	AllTimeAvg         float32 `json:"avg"`
	LastSuccessLoginAt *string `json:"lastSuccessLoginAt"`
	LastFailedLoginAt  *string `json:"LastFailedLoginAt"`
}

const bcryptCost = bcrypt.DefaultCost

var (
	db     = mysql.NewConnection()
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func Init() {
	_ = godotenv.Load("docker/.env")
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("Failed to migrate User struct.")
	}
}

func Register(c *gin.Context) {
	var promptedUser User

	err := c.BindJSON(&promptedUser)

	if promptedUser.Login == "" || promptedUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login or password can not be empty."})
		return
	}

	selected := db.Where("login = ?", promptedUser.Login).First(&User{})
	if selected.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exists."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(promptedUser.Password), bcryptCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash user password."})
		return
	}

	if promptedUser.Password == string(hashedPassword) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error storing your password safely."})
		return
	}

	promptedUser.Password = string(hashedPassword)

	result := db.Create(&promptedUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "There was an error while registering. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful."})
}

func Login(c *gin.Context) {
	var user User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse the data."})
		return
	}

	var found User

	result := db.Where("login = ?", user.Login).First(&found)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User does not exist."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect password."})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(found.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user session."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful.",
		"token":   token,
	})
}

func Data(c *gin.Context) {
	userID, err := Authorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err})
		return
	}

	var user User

	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		// Token is checked earlier so there can't be a possibility where token has invalid user
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not select user from database."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login":              user.Login,
		"tests":              user.TestsTaken,
		"avg":                user.AllTimeAvg,
		"lastSuccessLoginAt": user.LastSuccessLoginAt,
		"lastFailedLoginAt":  user.LastFailedLoginAt,
	})
}

func Authorize(c *gin.Context) (interface{}, error) {
	// TODO: change this to be an api call
	promptedTokens, ok := c.Request.Header["Authorization"]
	if !ok {
		return 0, fmt.Errorf("no auth method selected")
	}

	auth := promptedTokens[0]
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("unsupported auth method")
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("failed to get claims")
	}

	return claims["iss"], nil
}
