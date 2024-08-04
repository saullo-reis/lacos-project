package persons

import (
	"database/sql"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func RegisterPersons(c *gin.Context) {
	db, err := config.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao conectar com o banco de dados",
		})
		return
	}
	defer db.Close()

	var messagesError []string
	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro na leitura do json",
			"error": err.Error(),
		})
		return
	}

	validEmail := IsValidEmail(body.Email)
	if !validEmail && body.Email != "" {
		messagesError = append(messagesError, "Email não valido")
	}

	if body.CPF == "" {
		messagesError = append(messagesError, "CPF obrigatório")
	}
	if body.BirthDate == "" {
		messagesError = append(messagesError, "Data de nascimento obrigatória")
	}
	if body.Name == "" {
		messagesError = append(messagesError, "Nome da pessoa é obrigatório")
	}

	if len(messagesError) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"messages":    messagesError,
		})
		return
	}

	var personID int
	err = db.QueryRow("SELECT id_person FROM persons WHERE name = $1 and cpf = $2", body.Name, body.CPF).Scan(&personID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro na verificação da pessoa",
		})
		return
	}

	if personID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "Pessoa já registrada com esse Nome e CPF",
		})
		return
	}

	err = db.QueryRow(`INSERT INTO persons (
		name, birth_date, rg, cpf, cad_unico, nis, school, address, address_number, blood_type, neighborhood, city, cep, home_phone, cell_phone, contact_phone, email
	, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, 'Y') RETURNING id_person`,
		body.Name, body.BirthDate, body.RG, body.CPF, body.CadUnico, body.NIS, body.School, body.Address, body.AddressNumber, body.BloodType, body.Neighborhood, body.City, body.CEP, body.HomePhone, body.CellPhone, body.ContactPhone, body.Email).Scan(&personID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao inserir a pessoa no banco de dados",
		})
		return
	}

	_, err = db.Exec(`INSERT INTO responsible_person (
		name, id_person, rg, cpf, relationship, cell_phone) VALUES ($1, $2, $3, $4, $5, $6)`,
		body.ResponsiblePerson.Name, personID, body.ResponsiblePerson.RG, body.ResponsiblePerson.CPF, body.ResponsiblePerson.Relationship, body.ResponsiblePerson.CellPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao inserir a pessoa responsável no banco de dados",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"message":     "Pessoa registrada com sucesso",
		"data":        body,
	})
}
