package main

import (
	"jwtapi/database"
	"jwtapi/routes"

	"github.com/gin-gonic/gin"
)

const portNum = ":8080"

func main() {
	database.DbConnection()
	r := gin.Default()
	routes.Paths(r)
	r.Run(portNum)
}
