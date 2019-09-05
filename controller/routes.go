package controller

import (
	"github.com/citiaps/template-go-rest/middleware"
	"github.com/citiaps/template-go-rest/model"
	"github.com/gin-gonic/gin"
)

// Controllers
var authenticationController AuthenticationController
var dogController DogController

// Models

var dogModel model.Dog

func Routes(base *gin.RouterGroup) {
	// Middleware
	authNormal := middleware.LoadJWTAuth(middleware.LoginFunc)

	authenticationController.Routes(base, authNormal)
	dogController.Routes(base, authNormal)

}
