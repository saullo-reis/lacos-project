package persons

import "database/sql"


type ResponsiblePerson struct {
    IDPerson     int  `json:"id_person"`
    Name         string `json:"name"`
    Relationship string `json:"relationship"`
    RG           string `json:"rg"`
    CPF          string `json:"cpf"`
    CellPhone    string `json:"cell_phone"`
}

type Body struct {
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
    ResponsiblePerson ResponsiblePerson `json:"responsible_person"`
}

type SearchFieldsResponsablePerson struct {
	IDPerson     sql.NullInt64  `json:"id_person"`
	Name         sql.NullString `json:"name"`
	Relationship sql.NullString `json:"relationship"`
	RG           sql.NullString `json:"rg"`
	CPF          sql.NullString `json:"cpf"`
	CellPhone    sql.NullString `json:"cell_phone"`
}

type SearchFieldsPerson struct {
	Name              sql.NullString            `json:"name"`
	BirthDate         sql.NullString            `json:"birth_date"`
	RG                sql.NullString            `json:"rg"`
	CPF               sql.NullString            `json:"cpf"`
	CadUnico          sql.NullString            `json:"cad_unico"`
	NIS               sql.NullString            `json:"nis"`
	School            sql.NullString            `json:"school"`
	Address           sql.NullString            `json:"address"`
	AddressNumber     sql.NullString            `json:"address_number"`
	BloodType         sql.NullString            `json:"blood_type"`
	Neighborhood      sql.NullString            `json:"neighborhood"`
	City              sql.NullString            `json:"city"`
	CEP               sql.NullString            `json:"cep"`
	HomePhone         sql.NullString            `json:"home_phone"`
	CellPhone         sql.NullString            `json:"cell_phone"`
	ContactPhone      sql.NullString            `json:"contact_phone"`
	Email             sql.NullString            `json:"email"`
	CurrentAge        sql.NullInt64             `json:"current_age"`
	Active            sql.NullString                    `json:"active"`
	ResponsiblePerson ResponseResponsiblePerson `json:"responsible_person"`
}

type Params struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	School string `json:"school"`
	RG     string `json:"rg"`
	Active string `json:"active"`
}

type ResponseResponsiblePerson struct {
	IDPerson     int    `json:"id_person"`
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	RG           string `json:"rg"`
	CPF          string `json:"cpf"`
	CellPhone    string `json:"cell_phone"`
}

type Response struct {
	Name              string                    `json:"name"`
	BirthDate         string                    `json:"birth_date"`
	RG                string                    `json:"rg"`
	CPF               string                    `json:"cpf"`
	CadUnico          string                    `json:"cad_unico"`
	NIS               string                    `json:"nis"`
	School            string                    `json:"school"`
	Address           string                    `json:"address"`
	AddressNumber     string                    `json:"address_number"`
	BloodType         string                    `json:"blood_type"`
	Neighborhood      string                    `json:"neighborhood"`
	City              string                    `json:"city"`
	CEP               string                    `json:"cep"`
	HomePhone         string                    `json:"home_phone"`
	CellPhone         string                    `json:"cell_phone"`
	ContactPhone      string                    `json:"contact_phone"`
	Email             string                    `json:"email"`
	CurrentAge        int                       `json:"current_age"`
	Active            string                    `json:"active"`
	ResponsiblePerson ResponseResponsiblePerson `json:"responsible_person"`
}