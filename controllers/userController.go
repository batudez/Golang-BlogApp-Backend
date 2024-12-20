package controllers

import (
	"net/http"
	"os"
	"time"

	"blogapp.com/database"
	"blogapp.com/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	//"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//Get the email and password
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	//Create The user

	user := models.User{Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	//Get the email and pass off req body

	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//Look up requested user

	var user models.User
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or password",
		})

		return
	}

	//Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or password",
		})

		return
	}
	//generate jwttoken

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to create token",
		})
	}

	//send it back // TODO coookiye attık ileride sıkıntı çıkabilir
	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
