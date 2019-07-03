package model

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/citiaps/template-go-rest/util"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
)

// User : Usuario del sistema
type User struct {
	ID       bson.ObjectId `json:"id"              bson:"_id"`
	Email    string        `json:"email"           bson:"email"`
	Nombre   string        `json:"nombre"          bson:"nombre"`
	Rol      string        `json:"rol"             bson:"rol"`
	Hash     string        `json:"_hash"           bson:"_hash,omitempty"`
	Password string        `json:"password,omitempty"           bson:"password,omitempty"`
}

// LoadFromContext : Traer usuario desde contexto
func (userModel *User) LoadFromContext(c *gin.Context) *User {
	claims := jwt.ExtractClaims(c)
	var user User
	err := mapstructure.Decode(claims["user"], &user)
	if err != nil {
		panic(err)
	}
	user.ID = bson.ObjectIdHex(claims["user"].(map[string]interface{})["id"].(string))
	user.Hash = ""
	return &user
}

// Create : Traer usuario desde contexto
func (userModel *User) Create(user *User) error {

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()

	err := colUser.Insert(&user)

	return err
}

// GetUser : Se obtiene el usuario
func GetUser(c *gin.Context) {
	id := c.Param("id")

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var usuario User

	if err := colUser.FindId(bson.ObjectIdHex(id)).One(&usuario); err != nil {
		c.JSON(http.StatusNotFound, util.GetError("Usuario no encontrado", err))
	} else {
		c.JSON(http.StatusCreated, usuario)
	}
}

// CreateUsersBulk : Registrar usuario
func CreateUsersBulk(c *gin.Context) {
	var users []User
	e := c.BindJSON(&users)
	util.Check(e)

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	type Par struct {
		Usuario User
		Result  bson.M
	}

	type Respuesta struct {
		NoCreados []Par
		Creados   []Par
	}

	var resp Respuesta
	resp.Creados = make([]Par, 0)
	resp.NoCreados = make([]Par, 0)

	for _, u := range users {
		if u.Email == "" {
			var aux Par
			aux.Usuario = u
			aux.Result = bson.M{"mensaje": "No se especifico un email."}
			resp.NoCreados = append(resp.NoCreados, aux)
			continue
		}
		var temp User
		log.Printf("Buscando %s\n", u.Email)
		if err := colUser.Find(bson.M{"email": u.Email}).One(&temp); err != nil {
			//no existe, por lo que puedo crearlo

			u.ID = bson.NewObjectId()
			if err := colUser.Insert(&u); err != nil {
				var aux Par
				aux.Usuario = u
				aux.Result = bson.M{"mensaje": err.Error()}
				resp.NoCreados = append(resp.NoCreados, aux)
			} else {
				var aux Par
				aux.Usuario = u
				aux.Result = bson.M{"mensaje": "Usuario creado OK."}
				resp.Creados = append(resp.Creados, aux)
			}
		} else {
			var aux Par
			aux.Usuario = u
			aux.Result = bson.M{"mensaje": "Ese email ya esta siendo usado."}
			resp.NoCreados = append(resp.NoCreados, aux)
		}
	}
	c.JSON(http.StatusOK, resp)
}
