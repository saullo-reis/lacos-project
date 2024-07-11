package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "lacos.com/src/database/config"
)

type bodyToChangePassword struct {
	UserAdminId int    `json:"user_id"`
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(c *gin.Context) {
	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao conectar ao banco de dados",
		})
		return
	}
	defer db.Close()

	var body bodyToChangePassword
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "JSON invalido enviado",
		})
		return
	}

	if body.UserAdminId != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"status_code": 403,
			"message":     "Apenas o admin consegue modificar a senha de alguém",
		})
		return
	}

	if len(body.NewPassword) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "A senha deve ter no mínimo 8 chars!",
		})
		return
	}

	rows, err := db.Query("SELECT username FROM users WHERE username = $1", body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao buscar usuário no banco de dados",
		})
		return
	}

	var userToChangePassword string
	if rows.Next() {
		newPasswordHashed := hashPassword(body.NewPassword)
		err = rows.Scan(&userToChangePassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message":     "Erro ao escanear user",
			})
			return
		}
		_, err = db.Exec("UPDATE users SET password = $1 WHERE username = $2", newPasswordHashed, userToChangePassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message":     "Erro ao modificar senha: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"status_code": 200,
			"message":     "Senha de " + body.Username + " modificada com sucesso",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "Usuário não encontrado no banco de dados",
		})
		return
	}
}
