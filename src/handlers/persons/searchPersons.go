package persons

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lacos.com/src/database/config"
)

func SearchPersons(c *gin.Context) {
	db, err := config.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao conectar com o banco de dados",
		})
		return
	}
	defer db.Close()

	var params Params
	if err = c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"message":     "JSON inválido",
		})
		return
	}

	query := `
	SELECT 
		p.name, p.birth_date, p.rg, p.cpf, p.cad_unico, p.nis, p.school, p.address, p.address_number,
		p.blood_type, p.neighborhood, p.city, p.cep, p.home_phone, p.cell_phone, p.contact_phone, p.email, p.current_age,
		rp.id_person as rp_id_person, rp.name as rp_name, rp.relationship, rp.rg as rp_rg, rp.cpf as rp_cpf, rp.cell_phone as rp_cell_phone, p.active
	FROM 
		persons p
	LEFT JOIN 
		responsible_person rp ON p.id_person = rp.id_person
	WHERE 1=1`

	var args []interface{}
	paramIndex := 1

	if params.Name != "" {
		query += " AND p.name ILIKE $" + strconv.Itoa(paramIndex)
		args = append(args, "%"+params.Name+"%")
		paramIndex++
	}
	if params.CPF != "" {
		query += " AND p.cpf ILIKE $" + strconv.Itoa(paramIndex)
		args = append(args, "%"+params.CPF+"%")
		paramIndex++
	}
	if params.School != "" {
		query += " AND p.school ILIKE $" + strconv.Itoa(paramIndex)
		args = append(args, "%"+params.School+"%")
		paramIndex++
	}
	if params.RG != "" {
		query += " AND p.rg ILIKE $" + strconv.Itoa(paramIndex)
		args = append(args, "%"+params.RG+"%")
		paramIndex++
	}
	if params.Active != "" {
		query += " AND p.active = $" + strconv.Itoa(paramIndex)
		args = append(args, params.Active)
		paramIndex++
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao buscar pessoas: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var bodies []Response

	for rows.Next() {
		var body SearchFieldsPerson
		var responsiblePerson SearchFieldsResponsablePerson
		var response Response

		err := rows.Scan(
			&body.Name, &body.BirthDate, &body.RG, &body.CPF, &body.CadUnico, &body.NIS, &body.School, &body.Address, &body.AddressNumber,
			&body.BloodType, &body.Neighborhood, &body.City, &body.CEP, &body.HomePhone, &body.CellPhone, &body.ContactPhone, &body.Email, &body.CurrentAge,
			&responsiblePerson.IDPerson, &responsiblePerson.Name, &responsiblePerson.Relationship, &responsiblePerson.RG, &responsiblePerson.CPF, &responsiblePerson.CellPhone, &body.Active,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"message":     "Erro ao escanear resultado: " + err.Error(),
			})
			return
		}
		response.Address = body.Address.String
		response.AddressNumber = body.AddressNumber.String
		response.BirthDate = body.BirthDate.String
		response.BloodType = body.BloodType.String
		response.CEP = body.CEP.String
		response.CPF = body.CPF.String
		response.CadUnico = body.CadUnico.String
		response.CellPhone = body.CellPhone.String
		response.City = body.City.String
		response.ContactPhone = body.ContactPhone.String
		response.CurrentAge = int(body.CurrentAge.Int64)
		response.Email = body.Email.String
		response.HomePhone = body.HomePhone.String
		response.NIS = body.NIS.String
		response.Name = body.Name.String
		response.Neighborhood = body.Neighborhood.String
		response.RG = body.RG.String
		response.School = body.School.String
		response.CurrentAge = int(body.CurrentAge.Int64)
		response.Active = body.Active.String

		if responsiblePerson.IDPerson.Valid {
			response.ResponsiblePerson.IDPerson = int(responsiblePerson.IDPerson.Int64)
			response.ResponsiblePerson.Name = responsiblePerson.Name.String
			response.ResponsiblePerson.CPF = responsiblePerson.CPF.String
			response.ResponsiblePerson.RG = responsiblePerson.RG.String
			response.ResponsiblePerson.CellPhone = responsiblePerson.CellPhone.String
			response.ResponsiblePerson.Relationship = responsiblePerson.Relationship.String
		} else {
			response.ResponsiblePerson = ResponseResponsiblePerson{}
		}

		bodies = append(bodies, response)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     "Erro ao iterar resultados: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"message":     "Consulta finalizada com sucesso",
		"data":        bodies,
	})
}
