package db

import (
	"gopkg.in/mgo.v2"
	"github.com/citiaps/template-go-rest/utils"
  "os"
)

var (
	mongoHost     string
	mongoDatabase string
	mongoUsername string
	mongoPass string
	session *mgo.Session
)

func MongoSetup() {
	mongoHost = os.Getenv("DB_URL")
	mongoDatabase = os.Getenv("DB_DB")
	mongoUsername = os.Getenv("DB_USER")
	mongoPass = os.Getenv("DB_PASS")

	info := &mgo.DialInfo{
		Addrs:    []string{mongoHost},

		Database: mongoDatabase,
		Username: mongoUsername,
		Password: mongoPass,
	}
	var err error
	session, err = mgo.DialWithInfo(info)
	utils.Check(err)
}


func MongoSession() *mgo.Session {

	return session.Copy()
}

func MongoDatabase(session *mgo.Session) *mgo.Database {
	return session.DB(mongoDatabase)
}

func MongoCollection(collection string, db *mgo.Database) *mgo.Collection {
	return db.C(collection)
}
