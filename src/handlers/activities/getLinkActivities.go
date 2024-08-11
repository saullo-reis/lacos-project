package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

type Response struct {
	Name_activity string `json:"name_activity"`
	Id_activity   int    `json:"id_activity"`
	Hour_start    string `json:"hour_start"`
	Hour_end      string `json:"hour_end"`
}

func GetLinkActivities(c *gin.Context) {
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

	cpf_person := c.Param("cpf_person")
	query :=
		`SELECT act.id_activities, actl.name, act.hour_start, act.hour_end
	FROM persons p
	JOIN activities act ON act.id_person = p.id_person
	JOIN activity_list actl ON act.id_activity = actl.id_activity   
	WHERE cpf = $1`

	var valuesQueriedActivities Response
	rows, err := db.Query(query, cpf_person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Erro ao buscar atividades",
			"error_detail": err.Error(),
		})
		return
	}
	var response []Response
	for rows.Next() {
		err = rows.Scan(&valuesQueriedActivities.Id_activity, &valuesQueriedActivities.Name_activity, &valuesQueriedActivities.Hour_start, &valuesQueriedActivities.Hour_end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":       500,
				"message":      "Erro ao escanear atividades",
				"error_detail": err.Error(),
			})
			return
		}
		response = append(response, valuesQueriedActivities)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Busca feita com sucesso",
		"data":    response,
	})
}
