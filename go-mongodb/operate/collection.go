package operate

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCollectionOperate mongo collection operate
//
// param action string "action performed"
//
// param db *mongo.Database "a handle for a database"
//
// param col *mongo.Collection "a handle for a collection"
//
// return interface{}
func MongoCollectionOperate(action string, db *mongo.Database, col *mongo.Collection) interface{} {
	var res interface{}
	switch action {
	case "list":
		res = MongoListCollection(db)
	case "drop":
		res = MongoDropCollection(col)
	case "create":
		res = MongoCreateCollection(db)
	default:
		res = fmt.Sprintln("Invalid action, Please select the correct action (list, drop, create)!")
	}
	return res
}

func MongoListCollection(db *mongo.Database) []bson.M {
	result, err := db.ListCollections(
		context.TODO(),
		bson.D{},
	)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err := result.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(results)
	return results
}

func MongoDropCollection(col *mongo.Collection) string {
	err := col.Drop(context.TODO())
	if err != nil {
		return fmt.Sprintln("err: ", err)
	}
	return "drop collection success"
}

func MongoCreateCollection(db *mongo.Database) string {
	opts := options.CreateCollection().SetValidator(bson.M{})
	err := db.CreateCollection(context.TODO(), "users", opts)
	if err != nil {
		return fmt.Sprintln("err: ", err)
	}
	return "create collection success"
}
