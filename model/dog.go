package model

import "github.com/globalsign/mgo/bson"

// Dog : Perro de un usuario
type Dog struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name" bson:"name,omitempty"`
	Age   int           `json:"age" bson:"age,omitempty"`
	Owner bson.ObjectId `json:"owner" bson:"owner,omitempty"`
}

// Create : Crear perro por ID
func (dogModel *Dog) Create(dogDoc *Dog) error {
	col, session := GetCollection(CollectionNameDog)
	defer session.Close()
	dogDoc.ID = bson.NewObjectId()
	err := col.Insert(dogDoc)

	return err
}

// Get : Obtener perro por ID
func (dogModel *Dog) Get(id string) (*Dog, error) {
	col, session := GetCollection(CollectionNameDog)
	defer session.Close()
	var dogDoc Dog
	err := col.FindId(bson.ObjectIdHex(id)).One(&dogDoc)

	return &dogDoc, err
}

// Update : Actualizar perro por ID
func (dogModel *Dog) Update(id string, dogDoc Dog) error {

	col, session := GetCollection(CollectionNameDog)
	defer session.Close()
	err := col.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": dogDoc})
	return err
}

// Delete : Eliminar perro por ID
func (dogModel *Dog) Delete(id string) error {

	col, session := GetCollection(CollectionNameDog)
	defer session.Close()
	err := col.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

// Find : Obtener perro
func (dogModel *Dog) Find(query bson.M) ([]Dog, error) {

	col, session := GetCollection(CollectionNameDog)
	defer session.Close()
	dogs := []Dog{}

	err := col.Find(query).All(&dogs)
	return dogs, err
}
