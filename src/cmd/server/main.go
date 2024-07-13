package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"lacos.com/src/database/migrations"
	"lacos.com/src/handlers/user"
)

func main(){
	migrations.CreateTables()

	r := gin.Default()

	r.GET("/ping", user.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "ping"})
	})

	r.POST("/user/register", user.RegisterUser)
	r.POST("/user/login", user.LoginUser)
	r.PATCH("/user/changePassword", user.AuthMiddlewareAdmin(), user.ChangePassword)
	r.GET("/user/getUsers/:username", user.AuthMiddlewareAdmin(), user.GetUsers)
	r.DELETE("/user/deleteUser/:username", user.AuthMiddlewareAdmin(), user.DeleteUser)

	r.Run()
}