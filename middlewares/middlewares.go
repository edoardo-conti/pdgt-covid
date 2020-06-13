package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//CreateToken ...
func CreateToken(username string) (string, error) {
	var err error

	// creating jwt token (15 min)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

//CheckTokenValid ...
func CheckTokenValid(r *http.Request) error {
	tokenString := ""

	// check HTTP Header Authorization
	bearToken := r.Header.Get("Authorization")
	// token format: "Authorization <token>"
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		tokenString = strArr[1]
	}
	// parsing
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// verify token method "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong JWT method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err == nil {
		// token
		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			return err
		}
		return nil
	} else {
		// error
		return err
	}
}

//AuthMiddleware ...
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckTokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "accesso alla risorsa non autorizzato",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
