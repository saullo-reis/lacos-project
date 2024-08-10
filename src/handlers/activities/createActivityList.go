package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

type Activity_list struct{
	Name string `json:"name_activity"`
}

func CreateActivities(c *gin.Context){
	db, err := config.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message": "Falha na conexão com o banco de dados",
			"error": err.Error(),
		})
		return
	}

	var activity_list Activity_list

	if err := c.ShouldBindJSON(&activity_list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"message": "Falha na leitura do JSON",
		})
		return
	}

	if activity_list.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"message": "Um nome para a atividade é obrigatório",
		})
		return
	}

	_, err = db.Exec("INSERT INTO activity_list(name) VALUES($1)", activity_list.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message": "Falha ao inserir a atividade no banco de dados",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"message": "Atividade criada com sucesso",
	})
}