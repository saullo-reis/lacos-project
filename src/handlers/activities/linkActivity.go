package activities

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

type LinkActivityBody struct {
	Hour_start string `json:"hours_start"`
	Hour_end string `json:"hours_end"`
}

func LinkActivity(c *gin.Context) {
	cpf_person := c.Param("cpf_person")
	id_activity_list := c.Param("id_activity_list")
	id_period := c.Param("id_period")

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

	var body LinkActivityBody
	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Erro na leitura do JSON",
		})
		return
	}

	var error_details []string
	if body.Hour_start == "" || body.Hour_end == "" {
		if body.Hour_start == "" {
			error_details = append(error_details,
				"Horário Máximo obrigatório!")
		}
		if body.Hour_end == "" {
			error_details = append(error_details,
				"Horário Mínimo obrigatório!")
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":       400,
			"message":      "Algum parâmetro obrigatório faltou",
			"error_detail": error_details,
		})
		return
	}

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

	var name_activity string
	err = db.QueryRow("SELECT name FROM activity_list WHERE id_activity = $1", id_activity_list).Scan(&name_activity)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Atividade não encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Error ao verificar se a atividade existe",
			"error_detail": err.Error(),
		})
		return
	}
	id_activity_converted, err := strconv.Atoi(id_activity_list)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}

	id_period_converted, err := strconv.Atoi(id_period) 
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}

	_, err = db.Exec("INSERT INTO activities(id_activity, id_person, id_period, hour_start, hour_end) VALUES($1, $2, $3, $4, $5)", id_activity_converted, id_person, id_period_converted, body.Hour_start, body.Hour_end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       500,
			"message":      "Error ao linkar a atividade",
			"error_detail": err.Error(),
		})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": fmt.Sprintf("Atividade %s designada a %s", name_activity, name),
		},
	)
}
