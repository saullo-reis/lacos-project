package dbconfig

import (
	"fmt"
)

var User = "postgres"
var PostgresDriver = "postgres"
var Host = "localhost"
var Port = "5432"
var Password = "postgres"
var DbName = "lacos"

var DataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DbName)
