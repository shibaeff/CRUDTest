package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL = "http://localhost:8000/api"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Id        int64  `json:"id"`
}

type updateBody struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
}

var (
	client = http.Client{
		Timeout: time.Second,
	}
)

func sendCreate(postfix string, id int64) (d time.Duration) {
	url := baseURL + postfix
	person := User{
		FirstName: "test",
		LastName:  "test",
		UserName:  "test",
		Id:        id,
	}
	json_str, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(json_str)))
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	_, err = client.Do(r)
	d = time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func sendRead(postfix string, id int64) (d time.Duration) {
	url := baseURL + postfix
	r, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	q := r.URL.Query()
	q.Add("id", fmt.Sprintf("%d", id))
	r.URL.RawQuery = q.Encode()
	start := time.Now()
	_, err = client.Do(r)
	d = time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func sendUpd(postfix string, id int64) (d time.Duration) {
	url := baseURL + postfix
	var body updateBody
	body.UserName = "test1"
	update_body, err := json.Marshal(body)
	r, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(update_body))
	if err != nil {
		log.Fatal(err)
	}
	q := r.URL.Query()
	q.Add("id", fmt.Sprintf("%d", id))
	r.URL.RawQuery = q.Encode()
	start := time.Now()
	_, err = client.Do(r)
	d = time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func sendDelete(postfix string, id int64) (d time.Duration) {
	url := baseURL + postfix
	r, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	q := r.URL.Query()
	q.Add("id", fmt.Sprintf("%d", id))
	r.URL.RawQuery = q.Encode()
	start := time.Now()
	client.Do(r)
	d = time.Now().Sub(start)
	return
}

func main() {
	count := int64(12)
	delta := int64(0)
	square := int64(0)
	for i := int64(0); i < int64(count); i++ {
		d1 := sendCreate("/create", i)
		d2 := sendRead("/read", i)
		d3 := sendUpd("/update", i)
		d4 := sendDelete("/delete", i)
		sq := d1.Microseconds() + d2.Microseconds() + d3.Microseconds() + d4.Microseconds()
		delta += sq
		square += sq * sq
	}
	avg_delta := float64(delta) / float64(count)
	avg_square := float64(square) / float64(count)
	fmt.Printf("Mean %v\n", avg_delta)
	fmt.Printf("Var %v\n", math.Sqrt(avg_square-avg_delta*avg_delta))
}
