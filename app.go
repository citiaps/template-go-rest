package main

import (
	"github.com/citiaps/template-go-rest/utils"
	"github.com/citiaps/template-go-rest/db"
	"github.com/citiaps/template-go-rest/routes"
	"github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
	"log"
  "os"
)

func main() {

  godotenv.Load()

	db.MongoSetup()
	log.Println("Mongo Setup ok")
	app := gin.Default()

	app.Use(utils.CorsMiddleware())
	log.Println("Cors Setup ok")
	routes.Setup(app)
	log.Println("Routes Setup ok")

	app.Run(":" + os.Getenv("PORT"))
}
