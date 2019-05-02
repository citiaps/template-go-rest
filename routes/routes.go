package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/citiaps/template-go-rest/models"
	"github.com/citiaps/template-go-rest/authentication"
)

func Setup(app *gin.Engine) {
	authNormal := authentication.LoadJWTAuth(authentication.LoginFunc)

	// Refresh time can be longer than token timeout
	app.GET("/refresh_token",
		SetRoles(models.AuthRoles{Rol1:true,Rol2:true,Rol3:true}),
		authNormal.MiddlewareFunc(),
		authNormal.RefreshHandler)

	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	})

	//funcion de login: recibe un objeto {email: , pass:}
	app.POST("/login",authNormal.LoginHandler)

	//creacion de usuarios
	app.POST("/user",
		SetRoles(models.AuthRoles{Rol1:true,Rol2:false,Rol3:false}),
		authNormal.MiddlewareFunc(), models.CreateUser)

}


//funcion tipo middleware que define los roles que pueden realizar la siguiente funcion
func SetRoles(roles models.AuthRoles ) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set example variable
		c.Set("roles", roles)
		// before request
		c.Next()
	}
}
