package authentication

import (
	"github.com/citiaps/template-go-rest/models"
	"github.com/citiaps/template-go-rest/db"
	"github.com/citiaps/template-go-rest/utils"
	"reflect"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/appleboy/gin-jwt"
	"errors"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
)

var collectionNameUser = "user_model"

//Estructura que define el objeto recibido para el login
type login struct {
	Email string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//modificar esta funcion para hacer el login
func LoginFunc(c *gin.Context) (interface{}, error) {
	log.Print("LoginFunc\n")
  var loginVals login
  if err := c.BindJSON(&loginVals); err != nil {
    return "", jwt.ErrMissingLoginValues
  }

	session := db.MongoSession()
	defer session.Close()

	database := db.MongoDatabase(session)
	colUser := db.MongoCollection(collectionNameUser, database)

	var usuario models.User

	if err:= colUser.Find(bson.M{"email":loginVals.Email}).One(&usuario); err != nil {
		//return nil, jwt.ErrFailedAuthentication
		return nil, errors.New("Usuario y contraseña incorrectos")
	}else{
		if err:= ComparePasswords(usuario.Hash, loginVals.Password); err!=nil{
			//return nil, jwt.ErrFailedAuthentication
			return nil, errors.New("Usuario y contraseña incorrectos")
		}
		return usuario, nil
	}
}

//PARAMS:
//-storedHash: password almacenado en la BD
//-loginPass: el pasword ingresado por el usuario para hacer el login

//retorna:
//-true: si el password coincide
//-false: si el password no coincide
func ComparePasswords(storedHash string, loginPass string) error {
    byteHash := []byte(storedHash)
		loginHash := []byte(loginPass)
    err := bcrypt.CompareHashAndPassword(byteHash, loginHash)
    if err != nil {
        log.Println(err)
        return err
    }

    return nil
}

//funcion tipo middleware que define si el usuario esta autorizado a utilizar la siguiente funcion
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {
	log.Print("AuthorizatorFunc\n")

  if byteData, err:= json.Marshal(data); err!= nil {
		log.Print(err.Error())
		return false
	}else{
		var user models.User
		json.Unmarshal(byteData, &user)


		session := db.MongoSession()
		defer session.Close()

		database := db.MongoDatabase(session)
	  colUser := db.MongoCollection(collectionNameUser, database)

		var usuario models.User

	  if err := colUser.FindId(bson.ObjectId(user.Id)).One(&usuario); err != nil{
	    return false
	  } else {
			roles := c.MustGet("roles").(models.AuthRoles)
			log.Printf("%v\n",roles)
			if user.Rol == "ROL1" && usuario.Rol==user.Rol && roles.Rol1{
				log.Print("ERA ROL1 :D")
				return true
			}
			if user.Rol == "ROL2" && usuario.Rol==user.Rol && roles.Rol2{
				log.Print("ERA ROL2 :D")
				return true
			}
			if user.Rol == "ROL3" && usuario.Rol==user.Rol && roles.Rol3{
				log.Print("ERA ROL3 :D")
				return true
			}

			return false
	  }
	}
}

//funcion que se llama en caso de no estar autorizado a accesar al servicio
func UnauthorizedFunc(c *gin.Context, code int, message string) {
	log.Print("UnauthorizedFunc\n")
  c.JSON(code, gin.H{
    "message": message,
  })
}


//funcion que define lo que tendra el jwt que se enviara al realizarse el login
func PayLoad(data interface{}) jwt.MapClaims {
	log.Print("PayLoad\n")
	log.Print(reflect.TypeOf(data))
	log.Printf("%v\n",data)
  if v, ok := data.(models.User); ok {
		claim := jwt.MapClaims{
			"id":v.Id,
      "email": v.Email,
			"nombre": v.Nombre,
			"rol": v.Rol,
    }
		log.Printf("%v",claim)
    return claim
  }
  return jwt.MapClaims{}
}

func IdentityHandlerFunc(c *gin.Context) interface{} {
	log.Print("IdentityHandlerFunc\n")
	jwtClaims := jwt.ExtractClaims(c)
	log.Printf("%v",jwtClaims)
  //claims := jwt.ExtractClaims(c)
	//log.Printf("%v",claims)
  return &models.User{
		Id:  bson.ObjectIdHex(jwtClaims["id"].(string)),
		Email: jwtClaims["email"].(string),
		Nombre: jwtClaims["nombre"].(string),
		Rol: jwtClaims["rol"].(string),
  }
}


type loginFunc func(c *gin.Context) (interface{}, error)

func LoadJWTAuth( login loginFunc) *jwt.GinJWTMiddleware{
	log.Print("LoadJWTAuth\n")
	var key string
	var set bool
	key, set = os.LookupEnv("JWT_KEY")
	if !set{
		key = "string_largo_unico_por_proyecto"
	}

	log.Println("key: "+key)

  authMiddleware, err :=jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(key),
		//tiempo que define cuanto vence el jwt
		Timeout:     time.Hour*24*7,//una semana
		//tiempo maximo para poder refrescar el jwt token
		MaxRefresh:  time.Hour*24*7,

		PayloadFunc: PayLoad,
		IdentityHandler: IdentityHandlerFunc,
		Authenticator: login ,
		Authorizator: AuthorizatorFunc,
		Unauthorized: UnauthorizedFunc,
		//HTTPStatusMessageFunc: ResponseFunc,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	utils.Check(err)

	return authMiddleware

}
