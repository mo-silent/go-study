// go mongodb example
//
// go mongodb 示例
//
// Author  mogd  2022-05-06 CST
//
// Update  mogd  2022-05-06 CST

package main

import (
	"context"
	"fmt"

	db "go-mongodb/operate"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// uri Connection URI
const uri = "mongodb://39.101.244.245:27017/?maxPoolSize=20&w=majority&connect=direct"

var (
	Mclient *mongo.Client
	Mdb     *mongo.Database
	Mcol    *mongo.Collection
)

func main() {
	Mclient = MongoClient()
	defer func() {
		if err := Mclient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	// if err := Mclient.Ping(context.TODO(), readpref.Primary()); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Successfully connected and pinged.")
	// ListMongoDatabaseName()

	Mdb = Mclient.Database("test")
	Mcol = Mdb.Collection("test")
	opt := "collection"
	action := "create"
	var result interface{}
	switch opt {
	case "db":
		result = db.MongoDbOperate(action, Mclient, Mdb)
		fmt.Println(result)
	case "collection":
		result = db.MongoCollectionOperate(action, Mdb, Mcol)
		fmt.Println(result)
	case "doc":
		result = db.MongoDocumentOperate(action, Mcol)
		fmt.Println(result)
	default:
		fmt.Println("nothing to do!")
	}
	// Col = "test"

}

// MongoClient Create a database connection
//
// return *mongo.Client
func MongoClient() *mongo.Client {
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      "mogd",
		Password:      "admin",
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		panic(err)
	}
	return client
}

func InsertDoc(col *mongo.Collection) {

}
