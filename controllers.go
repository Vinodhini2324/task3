package controllers

import (
	"fmt"
	"jwtapi/database"
	"jwtapi/models"
	"strconv"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//Encryption Key

const secretKey = "secret"

// Schema for registering
type RData struct {
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// Schema for login

type LData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	_, err := c.Cookie("jwt")
	if err != http.ErrNoCookie {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You are already Logged In",
		})
		return
	}
	var person RData
	err = c.ShouldBindJSON(&person)
	if err != nil {
		fmt.Println(err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(person.Password), 14)

	// Need To Specify All the fields
	if person.UserName == "" || person.FirstName == "" || person.LastName == "" || person.Password == "" || person.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Specify All The Fields in the below values",
			"values":  "username, firstname, lastname, password, email",
		})

		return
	}

	user := models.Users{
		UserName:  person.UserName,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Email:     person.Email,
		Password:  password,
	}

	// Verifying the username and email

	// Checking Email
	database.DB.Where("email = ?", person.Email).First(&user)

	if user.Id > 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Email Already Exist",
		})
		return
	}

	// Checking Username
	database.DB.Where("user_name = ?", person.UserName).First(&user)

	if user.Id > 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "User Already Exist",
		})
		return
	}

	database.DB.Create(&user)

	c.JSON(200, user)
}

func Login(c *gin.Context) {
	_, err := c.Cookie("jwt")
	if err != http.ErrNoCookie {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You are already Logged In",
		})
		return
	}
	var person LData

	err = c.ShouldBind(&person)
	if err != nil {
		fmt.Println(err)
	}

	if person.Email == "" || person.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Enter all the fields", "values": "email , password"})
		return
	}
	var user models.Users
	database.DB.Where("email = ?", person.Email).First(&user)
	database.DB.Where("user_name = ?", person.Email).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User Not found"})
		return

	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(person.Password)); err != nil {
		c.Status(http.StatusNotFound)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Password"})
		return
	}
	clamins := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := clamins.SignedString([]byte(secretKey))
	fmt.Println("Token: ", token)
	fmt.Println("Error I am Facing: ", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not login",
		})
		return
	}

	// Storing In cookies
	c.SetCookie("jwt", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Logged In successfully",
	})
}

func Profile(c *gin.Context) {
	// Retriving The Cookie
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No Cookies Found",
		})
		return
	}
	// Validating The Token
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Go To Login Page",
		})
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	var user models.Users
	// Searching Into the database
	database.DB.Where("id = ?", claims.Issuer).First(&user)
	c.JSON(http.StatusAccepted, user)
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -12, "/", "localhost", false, true)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Logged Out Successfully",
	})
}
