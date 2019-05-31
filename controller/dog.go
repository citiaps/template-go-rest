package controller

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/citiaps/template-go-rest/middleware"
	"github.com/citiaps/template-go-rest/model"
	"github.com/citiaps/template-go-rest/util"
	"github.com/gin-gonic/gin"

	"github.com/globalsign/mgo/bson"
)

// DogController : Controlador de perro
type DogController struct {
}

// Routes : Define las rutas del controlador
func (dogController *DogController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) *gin.RouterGroup {

	// Dogs - Rutas
	dogRouter := base.Group("/dogs") //, middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc())
	{
		dogRouter.GET("", dogController.GetAll())
		// Al agregar asociar con usuario
		dogRouter.POST("", authNormal.MiddlewareFunc(), dogController.Create())
		dogRouter.GET("/:id", dogController.One())
		// Verificar en handler que el perro sea dueño de usuario
		dogRouter.PUT("/:id", authNormal.MiddlewareFunc(), dogController.Update())
		// Solo admin puede eliminar
		dogRouter.DELETE("/:id", middleware.SetRoles(RolAdmin), authNormal.MiddlewareFunc(), dogController.Delete())
	}
	return dogRouter
}

var dogModel model.Dog

// GetAll : Obtener todos los perros
func (dogController *DogController) GetAll() func(c *gin.Context) {
	return func(c *gin.Context) {

		groups, err := dogModel.Find(bson.M{})
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se pudo obtener la lista de grupos", err))
		}
		c.JSON(http.StatusOK, groups)
	}
}

// Create : Crear perro
func (dogController *DogController) Create() func(c *gin.Context) {
	return func(c *gin.Context) {

		// Traer Usuario
		user := userModel.LoadFromContext(c)
		var dog model.Dog
		err := c.Bind(&dog)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo decodificar json", err))
			return
		}
		// Asignar owner
		dog.Owner = user.ID
		err = dogModel.Create(&dog)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo insertar grupo", err))
			return
		}

		c.JSON(http.StatusOK, dog)
	}
}

// One : Obtener perro por _id
func (dogController *DogController) One() func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		group, err := dogModel.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontró grupo", err))
			return
		}
		c.JSON(http.StatusOK, group)
	}
}

// Update : Actualizar perro con _id
func (dogController *DogController) Update() func(c *gin.Context) {
	return func(c *gin.Context) {

		var dog model.Dog
		err := c.Bind(&dog)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo convertir collection json", err))
			return
		}
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, util.GetError("No se encuentra parametro :id", nil))
			return
		}

		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		// Update
		err = dogModel.Update(id, dog)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo actualizar grupo", err))
			return
		}

		c.String(http.StatusOK, "")
	}
}

// Delete : Eliminar perro por _id
func (dogController *DogController) Delete() func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		err := dogModel.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar grupo", err))
			return
		}
		c.String(http.StatusOK, "")
	}
}