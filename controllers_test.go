package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBCreateReadDeleteUser(t *testing.T) {
	godotenv.Load()
	usersCollection := db().Database("test").Collection("users")
	person := User{
		FirstName: "test",
		LastName:  "test",
		UserName:  "test",
		Id:        1,
	}
	w := httptest.NewRecorder()

	// CREATE
	json_str, err := json.Marshal(person)
	assert.NoError(t, err)
	r, err := http.NewRequest(http.MethodPost, "test", strings.NewReader(string(json_str)))
	assert.NoError(t, err)
	CreateUser(w, r)
	findRes := usersCollection.FindOne(context.TODO(), bson.D{{"username", "test"}})
	var found User
	findRes.Decode(&found)
	assert.Equal(t, "test", found.UserName)

	// READ
	r, err = http.NewRequest(http.MethodGet, "localhost:8080", strings.NewReader(""))
	assert.NoError(t, err)
	q := r.URL.Query()
	q.Add("id", "1")
	r.URL.RawQuery = q.Encode()
	ReadUser(w, r)
	json.Unmarshal(w.Body.Bytes(), &found)
	assert.Equal(t, "test", found.UserName)

	// UPDATE
	type updateBody struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		UserName  string `json:"username"`
	}
	var body updateBody
	body.UserName = "test1"
	update_body, err := json.Marshal(body)
	assert.NoError(t, err)
	r, err = http.NewRequest(http.MethodPut, "localhost:8080", bytes.NewReader(update_body))
	assert.NoError(t, err)
	q = r.URL.Query()
	q.Add("id", "1")
	r.URL.RawQuery = q.Encode()
	UpdateUser(w, r)
	findRes = usersCollection.FindOne(context.TODO(), bson.D{{"id", 1}})
	findRes.Decode(&found)
	assert.Equal(t, "test1", found.UserName)

	// DELETE
	r, err = http.NewRequest(http.MethodDelete, "localhost:8080", strings.NewReader(""))
	assert.NoError(t, err)
	q = r.URL.Query()
	q.Add("id", "1")
	r.URL.RawQuery = q.Encode()
	DeleteUser(w, r)
	findRes = usersCollection.FindOne(context.TODO(), bson.E{"id", 1})
	assert.NoError(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, findRes.Err())
}
