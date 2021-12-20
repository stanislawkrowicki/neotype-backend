package users

import "C"
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

	Login              string `gorm:"unique"`
	Password           string
	TestsTaken         int
	AllTimeAvg         *float32
	LastSuccessLoginAt *string
	LastFailedLoginAt  *string
}

const bcryptCost = bcrypt.DefaultCost

var (
	db     = mysql.NewConnection()
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func Init() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load environment variables.")
	}
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
	promptedTokens, ok := c.Request.Header["Authorization"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No auth selected."})
		return
	}

	auth := promptedTokens[0]
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unsupported authorization method."})
		return
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token."})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Could not get claims from token."})
		return
	}

	var user User

	result := db.Where("id = ?", claims["iss"]).First(&user)
	if result.Error != nil {
		// Token is checked earlier so there can't be a possibility where token has invalid user
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not select user from database."})
		return
	}

	c.JSON(http.StatusOK, user)
}
