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

// CreateToken funzione dedicata alla generazione di JSON Web Tokens
func CreateToken(username string) (string, error) {
	var err error

	/*
	 * creazione token jwt della durata di 15 min
	 * token codificato con username dell'utente e scadenza della validità imposta
	 */
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// encoding effettuato sfruttando la variabile d'ambiente 'JWT_ACCESS_SECRET' conservata in Heroku
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// CheckTokenValid verifica della validità del token in uso per l'inoltro di richieste HTTP
func CheckTokenValid(r *http.Request) error {
	tokenString := ""

	// si ricava il token dal campo Authorization dell'Header HTTP
	bearToken := r.Header.Get("Authorization")

	/*
	 * formato del token: "Authorization <token>"
	 * necessario split della stringa ricavata dall'header
	 */
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		tokenString = strArr[1]
	}

	// parsing del token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// verifica del metodo utilizzato per l'encoding, metodo richiesto: "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong JWT method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})

	// rilascio dell'esito
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

//AuthMiddleware middleware che funge da intermediario tra le richieste HTTP e gli handler ad essi associati
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// si sfrutta 'CheckTokenValid' per verificare la validità del token rilevato
		err := CheckTokenValid(c.Request)
		if err != nil {
			// nel caso in cui non fosse valido viene restituita una risposta 401
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "Errore, accesso alla risorsa non autorizzato.",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
