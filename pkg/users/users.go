package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"neotype-backend/pkg/mysql"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	gorm.Model

	Login              string `gorm:"unique"`
	Password           string
	LastSuccessLoginAt *string
	LastFailedLoginAt  *string
}

const bcryptCost = bcrypt.DefaultCost

var db = mysql.NewConnection()

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

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user session."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful.",
		"token":   token,
	})
}
