package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"blogapp.com/database"
	"blogapp.com/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		tokenString, _ = c.Cookie("Authorization")
	}

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRETKEY")), nil
	})

	// Eğer token parsing sırasında bir hata oluşmuşsa yetkilendirmeyi iptal et
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		database.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
