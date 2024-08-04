package persons

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func IfThenElseInt(condition bool, ifTrue, ifFalse int64) int64 {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func UpdatePersons(c *gin.Context) {
	db, err := config.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao conectar com o banco " + err.Error(),
		})
		return
	}
	defer db.Close()

	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error na leitura do json " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "JSON inválido " + err.Error(),
		})
		return
	}

	if body.CPF == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "CPF é obrigatório",
		})
		return
	}

	var responsiblePerson SearchFieldsResponsablePerson
	var person SearchFieldsPerson
	var idPerson int
	err = db.QueryRow(`
		SELECT 
			p.id_person, p.name, p.birth_date, p.rg, p.cad_unico, p.nis, p.school, p.address, p.address_number,
			p.blood_type, p.neighborhood, p.city, p.cep, p.home_phone, p.cell_phone, p.contact_phone, p.email, p.current_age,
			rp.id_person as rp_id_person, rp.name as rp_name, rp.relationship, rp.rg as rp_rg, rp.cpf as rp_cpf, rp.cell_phone as rp_cell_phone
		FROM 
			persons p
		LEFT JOIN 
			responsible_person rp ON p.id_person = rp.id_person
		WHERE p.cpf = $1`,
		body.CPF).Scan(&idPerson, &person.Name, &person.BirthDate, &person.RG, &person.CadUnico, &person.NIS, &person.School, &person.Address, &person.AddressNumber,
		&person.BloodType, &person.Neighborhood, &person.City, &person.CEP, &person.HomePhone, &person.CellPhone, &person.ContactPhone, &person.Email, &person.CurrentAge,
		&responsiblePerson.IDPerson, &responsiblePerson.Name, &responsiblePerson.Relationship, &responsiblePerson.RG, &responsiblePerson.CPF, &responsiblePerson.CellPhone)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": 400,
				"message":     "CPF não encontrado no banco de dados",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro na busca dos dados " + err.Error(),
		})
		return
	}

	personName := IfThenElse(body.Name == "", person.Name, body.Name).(string)
	personBirthDate := IfThenElse(body.BirthDate == "", person.BirthDate, body.BirthDate).(string)
	personRG := IfThenElse(body.RG == "", person.RG, body.RG).(string)
	personCadUnico := IfThenElse(body.CadUnico == "", person.CadUnico, body.CadUnico).(string)
	personNIS := IfThenElse(body.NIS == "", person.NIS, body.NIS).(string)
	personSchool := IfThenElse(body.School == "", person.School, body.School).(string)
	personAddress := IfThenElse(body.Address == "", person.Address, body.Address).(string)
	personAddressNumber := IfThenElse(body.AddressNumber == "", person.AddressNumber, body.AddressNumber).(string)
	personBloodType := IfThenElse(body.BloodType == "", person.BloodType, body.BloodType).(string)
	personNeighborhood := IfThenElse(body.Neighborhood == "", person.Neighborhood, body.Neighborhood).(string)
	personCity := IfThenElse(body.City == "", person.City, body.City).(string)
	personCEP := IfThenElse(body.CEP == "", person.CEP, body.CEP).(string)
	personHomePhone := IfThenElse(body.HomePhone == "", person.HomePhone, body.HomePhone).(string)
	personCellPhone := IfThenElse(body.CellPhone == "", person.CellPhone, body.CellPhone).(string)
	personContactPhone := IfThenElse(body.ContactPhone == "", person.ContactPhone, body.ContactPhone).(string)
	personEmail := IfThenElse(body.Email == "", person.Email, body.Email).(string)

	responsiblePersonName := IfThenElse(body.ResponsiblePerson.Name == "" , responsiblePerson.Name, body.ResponsiblePerson.Name).(string)
	responsiblePersonRelationship := IfThenElse(body.ResponsiblePerson.Relationship == "", responsiblePerson.Relationship, body.ResponsiblePerson.Relationship).(string)
	responsiblePersonRG := IfThenElse(body.ResponsiblePerson.RG == "", responsiblePerson.RG, body.ResponsiblePerson.RG).(string)
	responsiblePersonCPF := IfThenElse(body.ResponsiblePerson.CPF == "", responsiblePerson.CPF, body.ResponsiblePerson.CPF).(string)
	responsiblePersonCellPhone := IfThenElse(body.ResponsiblePerson.CellPhone == "", responsiblePerson.CellPhone, body.ResponsiblePerson.CellPhone).(string)

	query := `
		UPDATE persons
		SET 
			name = $1,
			birth_date = $2,
			rg = $3,
			cad_unico = $4,
			nis = $5,
			school = $6,
			address = $7,
			address_number = $8,
			blood_type = $9,
			neighborhood = $10,
			city = $11,
			cep = $12,
			home_phone = $13,
			cell_phone = $14,
			contact_phone = $15,
			email = $16
		WHERE cpf = $17
	`

	_, err = db.Exec(query,
		personName,
		personBirthDate,
		personRG,
		personCadUnico,
		personNIS,
		personSchool,
		personAddress,
		personAddressNumber,
		personBloodType,
		personNeighborhood,
		personCity,
		personCEP,
		personHomePhone,
		personCellPhone,
		personContactPhone,
		personEmail,
		body.CPF,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro na atualização da pessoa " + err.Error(),
		})
		return
	}

	query = `
		UPDATE responsible_person
		SET 
			name = $1,
			relationship = $2,
			rg = $3,
			cpf = $4,
			cell_phone = $5
		WHERE id_person = $6
	`
	_, err = db.Exec(query, responsiblePersonName, responsiblePersonRelationship, responsiblePersonRG, responsiblePersonCPF, responsiblePersonCellPhone, responsiblePerson.IDPerson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro na atualização da pessoa responsável " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"message":     "Pessoa atualizada com sucesso",
		"data":        body,
	})
}
