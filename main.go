package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "lacos.com/src/database/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"lacos.com/src/database/migrations"
	"lacos.com/src/handlers/persons"
	"lacos.com/src/handlers/user"
)

func init(){
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env: "+err.Error())
	}
}

func main(){
	fmt.Println(os.Getenv("USERPOSTGRES"))
	migrations.CreateTables()
	
	r := gin.Default()

	r.GET("/ping", user.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "ping"})
	})

	//CONNECTION HANDLERS
	r.POST("/user/login", user.LoginUser)

	//ADMIN HANDLERS
	r.PATCH("/user/changePassword", user.AuthMiddlewareAdmin(), user.ChangePassword)
	r.GET("/user/getUsers/:username", user.AuthMiddlewareAdmin(), user.GetUsers)
	r.POST("/user/register",user.AuthMiddlewareAdmin(), user.RegisterUser)
	r.DELETE("/user/deleteUser/:username", user.AuthMiddlewareAdmin(), user.DeleteUser)

	//NORMAL HANDLERS
	r.POST("/persons/register", user.AuthMiddleware(), persons.RegisterPersons)
	r.POST("/persons/search", user.AuthMiddleware(), persons.SearchPersons)
	r.DELETE("/persons/delete/:cpf", user.AuthMiddleware(), persons.DeletePerson)
	r.Run()
}