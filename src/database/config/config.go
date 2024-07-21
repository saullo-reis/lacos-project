package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	User           string
	PostgresDriver string
	Host           string
	Port           string
	Password       string
	DbName         string
	DataSourceName string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: " + err.Error())
	}

	User = os.Getenv("USERPOSTGRES")
	PostgresDriver = os.Getenv("POSTGRESDRIVER")
	Host = os.Getenv("HOSTPOSTGRES")
	Port = os.Getenv("PORTPOSTGRES")
	Password = os.Getenv("PASSWORDPOSTGRES")
	DbName = os.Getenv("DBNAME")

	DataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DbName)

	// Print variables to debug
	log.Printf("USERPOSTGRES: %s", User)
	log.Printf("POSTGRESDRIVER: %s", PostgresDriver)
	log.Printf("HOSTPOSTGRES: %s", Host)
	log.Printf("PORTPOSTGRES: %s", Port)
	log.Printf("PASSWORDPOSTGRES: %s", Password)
	log.Printf("DBNAME: %s", DbName)
}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open(PostgresDriver, DataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
