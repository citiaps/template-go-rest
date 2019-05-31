package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo"
)

var (
	session *mgo.Session
)

// GetCollection : Obtener coleccion desde la bd
func GetCollection(collection string) (*mgo.Collection, *mgo.Session) {
	s := session.Copy()
	return s.DB(os.Getenv("DB_DB")).C(collection), s
}

// LoadDatabase : Carga la base de datos y devuelve la session correspondiente
func LoadDatabase() {
	info := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("DB_URL")},
		Timeout:  30 * time.Second,
		Database: os.Getenv("DB_DB"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
	}
	var err error
	session, err = mgo.DialWithInfo(info)
	if err != nil {
		panic(fmt.Sprintf("Conection to DB,  %s/%s, Error: %s", os.Getenv("DB_URL"), os.Getenv("DB_DB"), err))
	}
	log.Printf("Conected to DB, %s/%s", os.Getenv("DB_URL"), os.Getenv("DB_DB"))

	session.SetMode(mgo.Monotonic, true)

	//Crear indices
	CreateIndex(session)
}
