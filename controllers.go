package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	createLog = "./logs/create.log"
	readLog   = "./logs/read.log"
	updLog    = "./logs/upd.log"
	delLog    = "./logs/del.log"

	CREATE = "CREATE"
	READ   = "READ"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
)

// struct for storing data
type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Id        int64  `json:"id"`
}

type Return struct {
	Dur int64 `json:"dur"`
}

var (
	usersCollection = db().Database("test").Collection("users")
)

func Test(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("hello")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Print(err)
	}
	_, err = usersCollection.InsertOne(context.TODO(), user)
	dur := time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted user")
	var ret Return
	ret.Dur = dur.Microseconds()
	json.NewEncoder(w).Encode(ret)
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	_id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	var result User
	err = usersCollection.FindOne(context.TODO(), bson.D{{"id", _id}}).Decode(&result)
	dur := time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	var ret Return
	ret.Dur = dur.Microseconds()
	json.NewEncoder(w).Encode(ret)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	_id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	type updateBody struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		UserName  string `json:"username"`
	}
	var body updateBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{"id", _id}}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	var update bson.D
	if body.UserName != "" {
		update = bson.D{{"$set", bson.D{{"username", body.UserName}}}}
	}
	if body.FirstName != "" {
		update.Map()["firstname"] = body.FirstName
	}
	if body.LastName != "" {
		update.Map()["lastname"] = body.LastName
	}
	updateResult := usersCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)
	dur := time.Now().Sub(start)
	var ret Return
	ret.Dur = dur.Microseconds()
	json.NewEncoder(w).Encode(ret)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	_id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	opts := options.Delete().SetCollation(&options.Collation{})
	res, err := usersCollection.DeleteOne(context.TODO(), bson.D{{"id", _id}}, opts)
	dur := time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("deleted %v documents\n", res.DeletedCount)
	var ret Return
	ret.Dur = dur.Microseconds()
	json.NewEncoder(w).Encode(ret)
}

//
//func getAllUsers(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	var results []primitive.M                                    //slice for multiple documents
//	cur, err := usersCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
//	if err != nil {
//
//		fmt.Println(err)
//
//	}
//	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor
//
//		var elem primitive.M
//		err := cur.Decode(&elem)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		results = append(results, elem) // appending document pointed by Next()
//	}
//	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
//	json.NewEncoder(w).Encode(results)
//}
