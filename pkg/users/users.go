package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"neotype-backend/pkg/config"
	"neotype-backend/pkg/mysql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type User struct {
	gorm.Model

	Login              string  `json:"login" gorm:"unique"`
	Password           string  `json:"password"`
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

	if len(promptedUser.Login) < 4 || len(promptedUser.Login) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username should be between 4 and 16 characters"})
		return
	}

	if len(promptedUser.Password) < 8 || len(promptedUser.Login) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password should be between 8 and 32 characters"})
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
	userID, err := ShouldAuthorize(c)
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
		"createdAt":          user.CreatedAt,
	})
}

func Username(c *gin.Context) {
	userID, err := ShouldAuthorize(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user User

	result := db.First(&user, "id = ?", userID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": user.Login})
}

func Authorize(c *gin.Context) {
	providedToken := c.Param("token")
	if providedToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is missing"})
		return
	}

	token, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is broken."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "iss": claims["iss"]})
}

func ShouldAuthorize(c *gin.Context) (interface{}, error) {
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

	serviceURI, err := config.GetBaseURL("users")
	if err != nil {
		log.Println("Could not get users URI!")
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", serviceURI+"/authorize/"+tokenString, nil)
	if err != nil {
		log.Printf("Failed to create new http request to authorization: %s", err)
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		log.Printf("Got error while calling users. Is the service dead?")
		return nil, fmt.Errorf("got error while calling users. Check if service is up")
	}

	var respBody gin.H

	bytes, _ := io.ReadAll(response.Body)
	if err = json.Unmarshal(bytes, &respBody); err != nil {
		log.Printf("Failed to unmarshal response from authorization service!")
		return nil, fmt.Errorf("failed to unmarshal response from authorization service")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if respBody["iss"] == nil {
		return nil, fmt.Errorf("no issuer in token")
	}
	
	return respBody["iss"], nil
}
