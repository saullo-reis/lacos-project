package persons

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "lacos.com/src/database/config"
)

func DeletePerson(c *gin.Context) {
	cpfToDelete := c.Param("cpf")

	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao conectar com o banco de dados: " + err.Error(),
		})
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao iniciar transação: " + err.Error(),
		})
		return
	}

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM persons WHERE cpf = $1)", cpfToDelete).Scan(&exists)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao verificar CPF: " + err.Error(),
		})
		return
	}

	if !exists {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "CPF não existe no nosso banco de dados",
		})
		return
	}

	_, err = tx.Exec("UPDATE persons SET active = 'N' WHERE cpf = $1", cpfToDelete)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao deletar pessoa: " + err.Error(),
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao confirmar transação: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"message":     "Pessoa deletada com sucesso do banco de dados",
	})
}
