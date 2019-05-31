package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/citiaps/template-go-rest/controller"
	"github.com/citiaps/template-go-rest/middleware"
	"github.com/citiaps/template-go-rest/model"
	"github.com/citiaps/template-go-rest/util"
	"github.com/gin-gonic/gin"
)

func main() {

	time.Local = time.UTC
	// Cargar Variables de entorno:
	util.LoadEnv()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Start template-go-rest")
	log.Printf("serverUp, %s", os.Getenv("ADDR"))

	// Cargar base de datos
	model.LoadDatabase()

	// Controllers
	var authenticationController controller.AuthenticationController
	var dogController controller.DogController
	//Raiz
	app := gin.Default()
	// CORS
	app.Use(middleware.CorsMiddleware())
	// Url Base
	base := app.Group("/api/v1/")
	authNormal := middleware.LoadJWTAuth(middleware.LoginFunc)

	// Collection - Rutas
	authenticationController.Routes(base, authNormal)
	dogController.Routes(base, authNormal)

	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	})

	http.ListenAndServe(os.Getenv("ADDR"), app)

}
