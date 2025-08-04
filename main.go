package main

import (
	"github.com/gin-gonic/gin"
	"kuhakuanime.com/db/sqlc"
	"kuhakuanime.com/routes"
)

func main() {
	db.Connect()
	server := gin.Default()
	routes.RegisterRoutes(server)

	server.Run(":8080")
}