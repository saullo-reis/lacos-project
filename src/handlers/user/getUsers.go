package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "lacos.com/src/database/config"
)

func GetUsers(c *gin.Context) {
	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Problema ao conectar ao banco de dados",
		})
		return
	}
	defer db.Close()

	username := c.Param("username")

	var users []string
	var rows *sql.Rows

	if username == "All" {
		rows, err = db.Query("SELECT username FROM users")
	} else {
		rows, err = db.Query("SELECT username FROM users WHERE username LIKE '%' || $1 || '%'", username)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Problema ao buscar o usuário",
		})
		return
	}
	defer rows.Close() 

	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message":     "Problema ao escanear o usuário",
			})
			return
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "Nenhum usuário encontrado com esse nome",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"message":     "Sucesso",
		"data":        users,
	})
}
