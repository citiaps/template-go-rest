package models

import (
  "gopkg.in/mgo.v2/bson"
  "net/http"
	"github.com/gin-gonic/gin"
  "log"
  "encoding/json"
	"github.com/citiaps/template-go-rest/db"
	"github.com/citiaps/template-go-rest/utils"
)

var collectionNameUser = "user_model"

type User struct {
	Id              bson.ObjectId  `json:"id"              bson:"_id"`
	Email           string         `json:"email"           bson:"email"`
  Nombre          string         `json:"nombre"          bson:"nombre"`
  Rol             string         `json:"rol"             bson:"rol"`	
  Hash            string         `json:"_hash"           bson:"_hash,omitempty"`
}

//Estructura que define el objeto para distintos roles
type AuthRoles struct{
	Rol1 bool
	Rol2 bool
	Rol3 bool
}




//funcion tipo middleware que define si el usuario esta autorizado a utilizar la siguiente funcion
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {
	log.Print("AuthorizatorFunc\n")

  if byteData, err:= json.Marshal(data); err!= nil {
		log.Print(err.Error())
		return false
	}else{
		var user User
		json.Unmarshal(byteData, &user)


		session := db.MongoSession()
		defer session.Close()

		database := db.MongoDatabase(session)
	  colUser := db.MongoCollection(collectionNameUser, database)

		var usuario User

	  if err := colUser.FindId(bson.ObjectId(user.Id)).One(&usuario); err != nil{
	    return false
	  } else {
			roles := c.MustGet("roles").(AuthRoles)
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


func GetUser(c *gin.Context){
  id := c.Param("id")

	session := db.MongoSession()
	defer session.Close()

	database := db.MongoDatabase(session)
  colUser := db.MongoCollection(collectionNameUser, database)

  var usuario User

  if err := colUser.FindId(bson.ObjectIdHex(id)).One(&usuario); err != nil{
    c.JSON(http.StatusNotFound, utils.FormatError("Usuario no encontrado",err.Error()))
  } else {
    c.JSON(http.StatusCreated, usuario)
  }
}


func CreateUser(c *gin.Context){
  var user User
  e := c.BindJSON(&user)
  utils.Check(e)
  user.Id = bson.NewObjectId()


  session := db.MongoSession()
	defer session.Close()

  database := db.MongoDatabase(session)
  colUser := db.MongoCollection(collectionNameUser, database)

  if err := colUser.Insert(&user); err != nil {
    log.Println("FALLO CREAR USUARIO")
    c.JSON(http.StatusInternalServerError, utils.FormatError( "Fallo al crear el usuario" ,err.Error()))
  } else{
    c.JSON(http.StatusCreated, user)
  }
}


func CreateUsersBulk(c *gin.Context){
	var users []User
  e := c.BindJSON(&users)
  utils.Check(e)

	session := db.MongoSession()
	defer session.Close()

  database := db.MongoDatabase(session)
  colUser := db.MongoCollection(collectionNameUser, database)

	type Par struct {
		Usuario User
		Result bson.M
	}

	type Respuesta struct{
		NoCreados []Par
		Creados []Par
	}


	var resp Respuesta
	resp.Creados = make([]Par,0)
	resp.NoCreados = make([]Par,0)

	for _, u := range users {
		if(u.Email==""){
			var aux Par
			aux.Usuario = u
			aux.Result = bson.M{"mensaje":"No se especifico un email."}
			resp.NoCreados = append(resp.NoCreados, aux)
			continue
		}
		var temp User
		log.Printf("Buscando %s\n",u.Email)
		if err := colUser.Find(bson.M{"email":u.Email}).One(&temp); err != nil{
	    //no existe, por lo que puedo crearlo

			u.Id = bson.NewObjectId()
			if err := colUser.Insert(&u); err != nil {
				var aux Par
				aux.Usuario = u
				aux.Result = bson.M{"mensaje":err.Error()}
		    resp.NoCreados = append(resp.NoCreados, aux)
		  } else{
				var aux Par
				aux.Usuario = u
				aux.Result = bson.M{"mensaje":"Usuario creado OK."}
				resp.Creados = append(resp.Creados, aux)
			}
	  } else {
			var aux Par
			aux.Usuario = u
			aux.Result = bson.M{"mensaje":"Ese email ya esta siendo usado."}
			resp.NoCreados = append(resp.NoCreados, aux)
	  }
	}
	c.JSON(http.StatusOK, resp)
}
