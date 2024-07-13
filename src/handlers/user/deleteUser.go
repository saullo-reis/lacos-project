package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "lacos.com/src/database/config"
)

func DeleteUser(c *gin.Context) {
	usernameToDelete := c.Param("username")
	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao conectar ao banco de dados",
		})
		return
	}
	defer db.Close()

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE username = $1", usernameToDelete).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": 404,
				"message":     "Usuário não encontrado",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message":     "Erro ao buscar o usuário no banco de dados",
			})
		}
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE username = $1", usernameToDelete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao deletar usuário",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"message":     "Usuário com o nome de " + usernameToDelete + " deletado",
	})
}
