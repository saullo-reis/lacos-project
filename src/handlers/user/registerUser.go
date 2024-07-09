package user

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "lacos.com/src/database/config"
)

type userStruch struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RegisterUser(c *gin.Context) {
	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":"Erro ao conectar com o banco de dados",
		})
		return
	}
	defer db.Close()

	var user userStruch
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": "400",
			"message":     "JSON Invalid",
		})
		return
	}

	rows, err := db.Query("SELECT 'Y' FROM users WHERE username = $1", user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message": "Error de servidor ao buscar pelo usuário",
		})
		return
	}
	defer rows.Close()

	var result sql.NullString
	if rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message": "Error de servidor ao buscar pelo usuário",
			})
			return
		}
	}

	if result.Valid && result.String == "Y" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message": "Esse usuário já existe!",
		})
		return
	} else {
		hashedPassword := hashPassword(user.Password)
		db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", user.Username, hashedPassword)
		c.JSON(http.StatusCreated, gin.H{
			"message": "Usuário "+ user.Username + " criado",
			"status_code": "201",
		})
		return
	}

}
