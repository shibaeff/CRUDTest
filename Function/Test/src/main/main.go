package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// this function is a helper function. i'm gonna use it
// in more complex cases. just preserving it for the boilerplate
func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI(
		"mongodb+srv://user:123@userscluster.whxir.mongodb.net/" +
			"test?retryWrites=true&w=majority")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Connected to MongoDB!")
	return client
}

var (
	usersCollection = db().Database("test").Collection("users")
)

// commands to run the deployment
// cd src && zip -r ../test-src.zip * && cd ..
// docker run -i openwhisk/actionloop-golang-v1.11 -compile main <test-src.zip >test-bin.zip
// ibmcloud fn action update Test test-bin.zip --native

func Main(args map[string]interface{}) map[string]interface{} {
	// time function exec, write it to response
	start := time.Now()
	dur := time.Now().Sub(start)
	res := make(map[string]interface{})
	res["dur"] = dur.Nanoseconds() / 1000
	return res
}
