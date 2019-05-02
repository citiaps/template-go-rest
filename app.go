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
func chequearVariables() bool{
  var set bool
  _, set = os.LookupEnv("GO_REST_ENV")
  if !set{
    return false
  }
  _, set = os.LookupEnv("DB_USER")
  if !set{
    return false
  }
  _, set = os.LookupEnv("DB_PASS")
  if !set{
    return false
  }
  _, set = os.LookupEnv("DB_DB")
  if !set{
    return false
  }
  _, set = os.LookupEnv("DB_URL")
  if !set{
    return false
  }
  _, set = os.LookupEnv("PORT")
  if !set{
    return false
  }
  _, set = os.LookupEnv("JWT_KEY")
  if !set{
    return false
  }
  return true
}

func main() {

  env := os.Getenv("TEMPLATE_ENV")
  if "" == env {
    env = "development"
  }

  godotenv.Load(".env." + env + ".local")
  if "test" != env {
    godotenv.Load(".env.local")
  }
  godotenv.Load(".env." + env)
  godotenv.Load() // The Original .env
  if !chequearVariables(){
    log.Println("ERROR: Variables de entorno necesarias no definidas")
    return
  }
  log.Println("Variables de entorno ok")  
	db.MongoSetup()
	log.Println("Mongo Setup ok")
	app := gin.Default()

	app.Use(utils.CorsMiddleware())
	log.Println("Cors Setup ok")
	routes.Setup(app)
	log.Println("Routes Setup ok")

	app.Run(":" + os.Getenv("PORT"))
}
