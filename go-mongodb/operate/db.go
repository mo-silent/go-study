package operate

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDbOperate mongo database operate
//
// param action string "action performed"
//
// param client *mongo.Client "new mongodb Client "
//
// param db *mongo.Database "a handle for a database"
//
// return interface{}
func MongoDbOperate(action string, client *mongo.Client, db *mongo.Database) interface{} {
	var res interface{}
	switch action {
	case "list":
		res = ListMongoDatabase(client)
	case "drop":
		res = DropMongoDatabase(db)
	default:
		res = fmt.Sprintln("Invalid action, Please select the correct action (list and drop)!")
	}
	return res
}

// ListMongoDatabase list all non-empty database
//
// return mongo.ListDatabasesResult
func ListMongoDatabase(client *mongo.Client) mongo.ListDatabasesResult {
	// Use a filter to select databases.
	result, err := client.ListDatabases(
		context.TODO(),
		bson.D{
			{Key: "empty", Value: false},
		})
	if err != nil {
		panic(err)
	}
	// for _, db := range result.Databases {
	// 	fmt.Println(db.Name)
	// }
	return result
}

// DropMongoDatabase drop mongo database
//
// return string
func DropMongoDatabase(db *mongo.Database) string {
	err := db.Drop(context.TODO())
	if err != nil {
		return fmt.Sprintln("err: ", err)
	}
	return fmt.Sprintln("drop database success")
}
