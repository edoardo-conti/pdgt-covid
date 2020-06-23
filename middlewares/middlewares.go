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
func CreateToken(username string, isAdmin bool) (string, error) {
	var err error

	/*
	 * creazione token jwt della durata di 15 min
	 * token codificato con username dell'utente e scadenza della validità imposta
	 */
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["admin"] = isAdmin
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
func CheckTokenValid(r *http.Request) (bool, error) {
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

	/*
		// parsing del token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// verifica del metodo utilizzato per l'encoding, metodo richiesto: "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("wrong JWT method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
		})
	*/

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// verifica del metodo utilizzato per l'encoding, metodo richiesto: "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong JWT method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})

	// per scorrere i campi claims -->
	// for key, val := range claims {
	//	 fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	// fmt.Printf("role: %v\n", claims["role"])

	// rilascio dell'esito
	if err == nil {
		// token
		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			return false, err
		}
		return claims["admin"].(bool), nil
	} else {
		// error
		return false, err
	}
}

//AuthMiddleware middleware che funge da intermediario tra le richieste HTTP e gli handler ad essi associati
func AuthMiddleware(authorize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// si sfrutta 'CheckTokenValid' per verificare la validità del token rilevato
		isAdmin, err := CheckTokenValid(c.Request)
		if err != nil {
			// nel caso in cui non fosse valido viene restituita una risposta 403
			c.JSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "Errore: accesso alla risorsa non autenticato.",
			})
			c.Abort()
			return
		}

		// verifica autorizzazione
		if authorize == 1 {
			if isAdmin != true {
				// nel caso in cui non fosse un admin viene restituita una risposta 401
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  401,
					"message": "Errore: accesso alla risorsa non autorizzato. (richiesto lv.admin)",
				})
				c.Abort()
				return
			}
			// role valido -> richiesta autorizzata
			c.Next()
		} else {
			// token valido -> richiesta autenticata
			c.Next()
		}
	}
}

//IsAdmin ...
/*
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	user := c.Get("username").(*jwt.Token)
		//	claims := user.Claims.(jwt.MapClaims)
		//	role := claims["role"].(string)


		claims := jwt.ExtractClaims(c)
		user, _ := c.Get(identityKey)
		c.JSON(200, gin.H{
			"userID":   claims[identityKey],
			"userName": user.(*User).UserName,
			"text":     "Hello World.",
		})

		//c.Next()
	}
}
*/
