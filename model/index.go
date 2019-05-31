package model

import (
	"log"
	"os"

	"github.com/globalsign/mgo"
)

// Nombre de la base de datos

// Constantes de nombres de colecciones
const (
	CollectionNameDog  = "Dog"
	CollectionNameUser = "User"
)

// CreateIndex : Crea indices en las colecciones de mongo
func CreateIndex(session *mgo.Session) {

	// Creacion de indices
	err := session.DB(os.Getenv("DB_DB")).C(CollectionNameDog).EnsureIndexKey("age", "_id")
	if err != nil {
		log.Printf("Error al crear indice %s en : %s, %s", "TaggedData", "idData", err)
		panic("No se pudo crear indice")
	}

	log.Print("Creando indices")
}
