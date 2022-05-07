package operate

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDocumentOperate mongo document operate
//
// param action string "action performed"
//
// param coll *mongo.Collection "a handle for a Collection"
//
// return interface{}
func MongoDocumentOperate(action string, coll *mongo.Collection) interface{} {
	var res interface{}
	switch action {
	case "list":
		res = MongoListAllDocument(coll)
	case "drop":
		res = MongoDropDocument(coll)
	case "create":
		docs := []interface{}{
			bson.D{{Key: "name", Value: "Alice"}},
			bson.D{{Key: "name", Value: "Bob"}},
		}
		res = MongoCreateDocument(coll, docs)
	default:
		res = fmt.Sprintln("Invalid action, Please select the correct action (list, drop, create)!")
	}
	return res
}

func MongoListAllDocument(coll *mongo.Collection) []bson.M {
	// opts := options.Find().SetSort(bson.D{{Key: "age", Value: 1}})
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return []bson.M{bson.M{"err": err}}
	}

	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		return []bson.M{bson.M{"err": err}}
	}
	return results
}

func MongoDropDocument(coll *mongo.Collection) string {
	return "drop document success!"
}

func MongoCreateDocument[T []interface{} | interface{}](coll *mongo.Collection, docs T) []interface{} {
	return nil
}
