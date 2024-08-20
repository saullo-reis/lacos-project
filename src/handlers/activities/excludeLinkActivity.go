package activities

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

func ExcludeLinkActivity(c *gin.Context) {
	cpf_person := c.Param("cpf_person")
	id_activities := c.Param("id_activities")

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

	var name string
	var id_person int
	err = db.QueryRow("SELECT name, id_person FROM persons WHERE cpf = $1", cpf_person).Scan(&name, &id_person)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Pessoa não encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Error ao verificar se a pessoa existe",
			"error_detail": err.Error(),
		})
		return
	}

	//FAZER UMA VERIFICAÇÃO SE O ID DIGITADO NA URL EXISTE MESMO

	// var name_activity string
	// err = db.QueryRow("SELECT id_activities name FROM activities WHERE id_activities = $1", id_activities).Scan(&name_activity)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"status":  400,
	// 			"message": "Atividade não encontrada",
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":       500,
	// 		"message":      "Error ao verificar se a atividade existe",
	// 		"error_detail": err.Error(),
	// 	})
	// 	return
	// }

	_, err = db.Exec("DELETE FROM activities WHERE id_activities = $1 AND id_person = $2", id_activities, id_person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Error ao deletar atividade de uma pessoa",
			"error_detail": err.Error(),
		})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": fmt.Sprintf("Atividade retirada de %s", name),
		},
	)
}
