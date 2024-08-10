package activities

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

func DeleteActivity(c *gin.Context) {
	activity_to_delete := c.Param("id_activity")
	db, err := config.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Erro ao conectar com o banco de dados",
			"error_detail": err.Error(),
		})
		return
	}
	defer db.Close()

	var exist string
	err = db.QueryRow("SELECT 'Y' FROM activity_list WHERE id_activity = $1", activity_to_delete).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Nenhuma atividade com esse ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Erro ao verificar actividade",
			"error_detail": err.Error(),
		})
		return
	}

	_, err = db.Exec("DELETE FROM activity_list WHERE id_activity = $1 ", activity_to_delete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Erro ao deletar atividade",
			"error_detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Atividade deletada com sucesso",
	})
}
