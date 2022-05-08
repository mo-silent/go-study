package operate

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"

	bson "go.mongodb.org/mongo-driver/bson"
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
		docs := bson.D{{Key: "name", Value: "bob"}}
		res = MongoDropDocument(coll, docs)
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

// MongoListAllDocument list mongodb document
//
// param coll *mongo.Collection "a handle for a Collection"
//
// return []bson.M
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

// MongoDropDocument drop mongodb document
//
// param coll *mongo.Collection "a handle for a Collection"
//
// param docs interface{} actually bson.D{}
//
// return string
func MongoDropDocument(coll *mongo.Collection, docs interface{}) string {
	res, err := coll.DeleteMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err)
		return fmt.Sprintf("err: %v\n", err)
	}
	return fmt.Sprintf("deleted %v documents\n", res.DeletedCount)
}

// MongoCreateDocument create mongodb document
//
// param coll *mongo.Collection "a handle for a Collection"
//
// param docs generics T actually []bson.D{} or bson.D{}
//
// return []interface{}
func MongoCreateDocument(coll *mongo.Collection, docs interface{}) interface{} {
	switch docs.(type) {
	case []interface{}:
		opts := options.InsertMany().SetOrdered(false)
		res, err := coll.InsertMany(context.TODO(), docs.([]interface{}), opts)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return res.InsertedIDs
	case interface{}:
		res, err := coll.InsertOne(context.TODO(), docs)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return res.InsertedID
	default:
		fmt.Println(reflect.TypeOf(docs))
		return fmt.Sprintln("Invalid param, Please input the correct param ([]interface{}, interface{})!")
	}
}
