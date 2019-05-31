package middleware

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/citiaps/template-go-rest/model"
	"github.com/citiaps/template-go-rest/util"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

//Estructura que define el objeto recibido para el login
type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//LoginFunc : modificar esta funcion para hacer el login
func LoginFunc(c *gin.Context) (interface{}, error) {

	log.Print("LoginFunc\n")
	var loginVals login
	if err := c.BindJSON(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	colUser, session := model.GetCollection(model.CollectionNameUser)
	defer session.Close()
	var usuario model.User

	if err := colUser.Find(bson.M{"email": loginVals.Email}).One(&usuario); err != nil {
		//return nil, jwt.ErrFailedAuthentication
		return nil, errors.New("Usuario y contraseña incorrectos")
	} else {
		if err := ComparePasswords(usuario.Hash, loginVals.Password); err != nil {
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

func GeneratePassword(password string) string {
	binpwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(binpwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

// AuthorizatorFunc : funcion tipo middleware que define si el usuario esta autorizado a utilizar la siguiente funcion
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {
	user := data.(map[string]interface{})
	colUser, session := model.GetCollection(model.CollectionNameUser)
	defer session.Close()
	var usuario model.User

	if err := colUser.FindId(bson.ObjectIdHex(user["id"].(string))).One(&usuario); err != nil {
		return false
	}
	roles, exists := c.Get("roles")
	if !exists {
		return true
	}
	for _, r := range roles.([]string) {
		if usuario.Rol == r {
			return true
		}
	}
	return false

}

//UnauthorizedFunc : funcion que se llama en caso de no estar autorizado a accesar al servicio
func UnauthorizedFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

//PayLoad : funcion que define lo que tendra el jwt que se enviara al realizarse el login
func PayLoad(data interface{}) jwt.MapClaims {
	if v, ok := data.(model.User); ok {
		claim := jwt.MapClaims{
			"user": data.(model.User),
			"rol":  v.Rol,
		}
		log.Printf("%v", claim)
		return claim
	}
	return jwt.MapClaims{}
}

func IdentityHandlerFunc(c *gin.Context) interface{} {
	jwtClaims := jwt.ExtractClaims(c)
	return jwtClaims["user"]
}

type loginFunc func(c *gin.Context) (interface{}, error)

func LoadJWTAuth(login loginFunc) *jwt.GinJWTMiddleware {
	log.Print("LoadJWTAuth\n")
	var key string
	var set bool
	key, set = os.LookupEnv("JWT_KEY")
	if !set {
		key = "string_largo_unico_por_proyecto"
	}

	log.Println("key: " + key)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "test zone",
		Key:   []byte(key),
		//tiempo que define cuanto vence el jwt
		Timeout: time.Hour * 24 * 7, //una semana
		//tiempo maximo para poder refrescar el jwt token
		MaxRefresh: time.Hour * 24 * 7,

		PayloadFunc:     PayLoad,
		IdentityHandler: IdentityHandlerFunc,
		Authenticator:   login,
		Authorizator:    AuthorizatorFunc,
		Unauthorized:    UnauthorizedFunc,
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

	util.Check(err)

	return authMiddleware

}

//SetRoles : funcion tipo middleware que define los roles que pueden realizar la siguiente funcion
func SetRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set example variable
		c.Set("roles", roles)
		// before request
		c.Next()
	}
}
