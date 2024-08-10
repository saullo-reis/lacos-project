package activities

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

type ActivityToUpdate struct {
	Name string `json:"name_activity"`
}

func UpdateActivity(c *gin.Context) {
	var activity_to_update ActivityToUpdate
	id_activity := c.Param("id_activity")
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

	if err = c.ShouldBindJSON(&activity_to_update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Erro ao ler JSON",
			"error_detail": err.Error(),
		})
		return
	}

	if activity_to_update.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Um novo nome é obrigatório",
		})
		return
	}

	var exist string
	err = db.QueryRow("SELECT 'Y' FROM activity_list WHERE id_activity = $1", id_activity).Scan(&exist)
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

	_, err = db.Exec("UPDATE activity_list SET name = $1 WHERE id_activity = $2 ", activity_to_update.Name, id_activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Erro ao atualizar atividade",
			"error_detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Nome da atividade atualizada com sucesso",
	})
}
