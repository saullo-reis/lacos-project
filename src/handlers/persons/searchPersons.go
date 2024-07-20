package persons

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "lacos.com/src/database/config"
)

type SearchFieldsResponsablePerson struct {
    IDPerson     sql.NullInt64  `json:"id_person"`
    Name         sql.NullString `json:"name"`
    Relationship sql.NullString `json:"relationship"`
    RG           sql.NullString `json:"rg"`
    CPF          sql.NullString `json:"cpf"`
    CellPhone    sql.NullString `json:"cell_phone"`
}

type SearchFieldsPerson struct {
	Name           sql.NullString   `json:"name"`
    BirthDate      sql.NullString  `json:"birth_date"`
    RG             sql.NullString   `json:"rg"`
    CPF            sql.NullString   `json:"cpf"`
    CadUnico       sql.NullString   `json:"cad_unico"`
    NIS            sql.NullString   `json:"nis"`
    School         sql.NullString   `json:"school"`
    Address        sql.NullString   `json:"address"`
    AddressNumber  sql.NullString   `json:"address_number"`
    BloodType      sql.NullString   `json:"blood_type"`
    Neighborhood   sql.NullString   `json:"neighborhood"`
    City           sql.NullString   `json:"city"`
    CEP            sql.NullString   `json:"cep"`
    HomePhone      sql.NullString   `json:"home_phone"`
    CellPhone      sql.NullString   `json:"cell_phone"`
    ContactPhone   sql.NullString   `json:"contact_phone"`
    Email          sql.NullString   `json:"email"`
    CurrentAge     sql.NullInt64  `json:"current_age"`
    ResponsiblePerson ResponseResponsiblePerson `json:"responsible_person"`
}

type Params struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	School string `json:"school"`
	RG     string `json:"rg"`
}

type ResponseResponsiblePerson struct {
    IDPerson     int  `json:"id_person"`
    Name         string `json:"name"`
    Relationship string `json:"relationship"`
    RG           string `json:"rg"`
    CPF          string `json:"cpf"`
    CellPhone    string `json:"cell_phone"`
}

type Response struct {
    Name           string   `json:"name"`
    BirthDate      string  `json:"birth_date"`
    RG             string   `json:"rg"`
    CPF            string   `json:"cpf"`
    CadUnico       string   `json:"cad_unico"`
    NIS            string   `json:"nis"`
    School         string   `json:"school"`
    Address        string   `json:"address"`
    AddressNumber  string   `json:"address_number"`
    BloodType      string   `json:"blood_type"`
    Neighborhood   string   `json:"neighborhood"`
    City           string   `json:"city"`
    CEP            string   `json:"cep"`
    HomePhone      string   `json:"home_phone"`
    CellPhone      string   `json:"cell_phone"`
    ContactPhone   string   `json:"contact_phone"`
    Email          string   `json:"email"`
    CurrentAge     int  `json:"current_age"`
    ResponsiblePerson ResponseResponsiblePerson `json:"responsible_person"`
}

func SearchPersons(c *gin.Context) {
	db, err := sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)
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
		rp.id_person as rp_id_person, rp.name as rp_name, rp.relationship, rp.rg as rp_rg, rp.cpf as rp_cpf, rp.cell_phone as rp_cell_phone
	FROM 
		persons p
	LEFT JOIN 
		responsible_person rp ON p.id_person = rp.id_person
	WHERE 1=1`

	var args []interface{}
	if params.Name != "" {
		query += " AND (p.name ILIKE $1)"
		args = append(args, "%"+params.Name+"%")
	}
	if params.CPF != "" {
		query += " AND (p.cpf ILIKE $2)"
		args = append(args, "%"+params.CPF+"%")
	}
	if params.School != "" {
		query += " AND (p.school ILIKE $3)"
		args = append(args, "%"+params.School+"%")
	}
	if params.RG != "" {
		query += " AND (p.rg ILIKE $4)"
		args = append(args, "%"+params.RG+"%")
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
            &responsiblePerson.IDPerson, &responsiblePerson.Name, &responsiblePerson.Relationship, &responsiblePerson.RG, &responsiblePerson.CPF, &responsiblePerson.CellPhone,
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

		if responsiblePerson.IDPerson.Valid {
			response.ResponsiblePerson.IDPerson = int(responsiblePerson.IDPerson.Int64)
			response.ResponsiblePerson.Name = responsiblePerson.Name.String
			response.ResponsiblePerson.CPF = responsiblePerson.CPF.String
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
