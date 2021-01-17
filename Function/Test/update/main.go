package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
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
	usersCollection *mongo.Collection
)

func init() {
	usersCollection = db().Database("test").Collection("users")
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Id        int64  `json:"id"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(user User) (Response, error) {

	var update bson.D
	if user.UserName != "" {
		update = bson.D{{"$set", bson.D{{"username", user.UserName}}}}
	}
	if user.FirstName != "" {
		update.Map()["firstname"] = user.FirstName
	}
	if user.LastName != "" {
		update.Map()["lastname"] = user.LastName
	}

	start := time.Now()
	_, err := usersCollection.UpdateOne(context.TODO(), bson.D{{"id", user.Id}}, update)
	dur := time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"dur": fmt.Sprintf("%d", dur.Nanoseconds()/1000),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)
	log.Println("Inserted user")
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "test-handler",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
