package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

type Activities struct {
	IdActivity string `json:"id_activity"`
	Name       string `json:"name"`
}

func GetActivitiesList(c *gin.Context) {
	nameToSearch := c.Param("name")
	db, err := config.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message": "Erro ao conectar com o banco de dados",
			"error_detail": err.Error(),
		})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id_activity, name FROM activity_list WHERE name LIKE '%' || $1 || '%'", nameToSearch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message": "Erro ao executar a consulta no banco de dados",
			"error_detail": err.Error(),
		})
		return
	}
	defer rows.Close()

	var activities []Activities
	for rows.Next() {
		var activity Activities
		err = rows.Scan(&activity.IdActivity, &activity.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"message": "Erro ao buscar atividade",
				"error_detail": err.Error(),
			})
			return
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message": "Erro ao iterar sobre as atividades",
			"error_detail": err.Error(),
		})
		return
	}

	if len(activities) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"message": "Não encontrou nenhuma atividade com esse nome",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"message": "Busca bem sucedida",
		"data": activities,
	})
}
