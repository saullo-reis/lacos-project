package user

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	dbconfig "lacos.com/src/database/config"
)

type body struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string
	Password string
}

var jwtKeySecret string
var jwtKey = []byte(jwtKeySecret)

func generateJWT(usename string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), 
		Issuer:    "secure-chat",
		Subject:   usename,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func hasherPassword(password string) string{
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func LoginUser(c *gin.Context) {
	jwtKeySecret = os.Getenv("JWTKEYSECRET")

	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":"Erro ao conectar com o banco de dados",
		})
		return
	}
	defer db.Close()

	var body body
	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message": "JSON Invalido",
		})
		return
	}

	passwordHashed := hasherPassword(body.Password)

	var userQueried User
	rows, err := db.Query("SELECT username, password FROM users WHERE username = $1", body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 500,
			"message":    "Erro ao ler dados do usuário no banco de dados",
		})
		return
	}
	defer rows.Close()

	if rows.Next(){
		var username, password string
		if err = rows.Scan(&username, &password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": 400,
				"message":    "Erro ao escanear usuário",
			})
			return
		}
		userQueried = User{
			Username: username,
			Password: password,
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"status_code": 404,
			"message":    "Usuário não registrado",
		})
		return
	}

	if passwordHashed == userQueried.Password {
		token, err := generateJWT(userQueried.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message": "Error ao gerar o token",
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"status_code": 202,
			"message":    "Usuário autorizado",
			"token": token,
		})
		return
	}else{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code": 401,
			"message": "Usuário não autorizado",
		})
		return
	}

}
