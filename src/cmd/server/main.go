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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "ping"})
	})
	r.POST("/register/user", user.RegisterUser)
	r.Run()
}