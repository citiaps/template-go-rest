package controller

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/citiaps/template-go-rest/middleware"
	"github.com/citiaps/template-go-rest/model"
	"github.com/citiaps/template-go-rest/util"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// Roles en el sistema
const (
	RolAdmin = "Admin"
	RolUser  = "User"
)

// AuthenticationController : Estructura controladora de las colecciones
type AuthenticationController struct {
}

// Routes : Define las rutas del controlador
func (authenticationController *AuthenticationController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) {

	// Refresh time can be longer than token timeout
	base.GET("/refresh_token",
		middleware.SetRoles(RolAdmin, RolUser),
		authNormal.MiddlewareFunc(),
		authNormal.RefreshHandler)

	//funcion de login: recibe un objeto {email: , pass:}
	base.POST("/login", authNormal.LoginHandler)

	//creacion de usuarios
	base.POST("/user",
		CreateUser)
}

var userModel model.User

// CreateUser : Registrar usuario
func CreateUser(c *gin.Context) {
	var user model.User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo registrar al usuario", e))
		return
	}
	user.ID = bson.NewObjectId()
	user.Hash = middleware.GeneratePassword(user.Password)
	user.Password = ""
	user.Rol = RolUser
	if err := userModel.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, util.GetError("Fallo al crear el usuario", err))
		return
	}
	c.String(http.StatusCreated, "")
}
