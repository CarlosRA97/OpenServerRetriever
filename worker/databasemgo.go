package worker

import (
	"fmt"
	"upper.io/db.v2"
	"upper.io/db.v2/mongo"
)

var settings = mongo.ConnectionURL{
	Host:     config.Mongo.Host,
	User:     config.Mongo.User,
	Password: config.Mongo.Password,
	Database: config.Mongo.Database,
}

func connectTo(url mongo.ConnectionURL) db.Database {

	database, err := mongo.Open(url)
	Check(err)
	return database
}

func getCollection(name string, database db.Database) db.Collection {
	return database.Collection(name)
}

func GetData(sliceOfStructs interface{}, collec string) {
	database := connectTo(settings)
	collection := getCollection(collec, database)
	err := collection.Find().All(sliceOfStructs)
	Check(err)
}

func PutData(data interface{}, collec string) {
	database := connectTo(settings)
	collection := getCollection(collec, database)
	fmt.Println(collection.Insert(data))
}
