package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL = "https://us-central1-pivotal-store-301811.cloudfunctions.net"
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

func meanVar(delta, squares, count int64) string {
	avg_delta := float64(delta) / float64(count)
	avg_square := float64(squares) / float64(count)
	variance := avg_square - avg_delta*avg_delta
	mean := avg_delta
	return fmt.Sprintf("mean %f deviation %f", mean, math.Sqrt(variance))
}

func main() {
	count := int64(10)
	delta := int64(0)
	square := int64(0)
	bar := pb.StartNew(int(count))
	var (
		create_s, create_sq, read_s, read_sq, upd_s, upd_sq, del_s, del_sq int64
	)
	for i := int64(0); i < int64(count); i++ {
		d1 := sendCreate("/create", i).Microseconds()
		d2 := sendRead("/read", i).Microseconds()
		d3 := sendUpd("/update", i).Microseconds()
		d4 := sendDelete("/delete", i).Microseconds()
		sq := d1 + d2 + d3 + d4
		delta += sq
		square += sq * sq

		create_s += d1
		create_sq += d1 * d1
		read_s += d2
		read_sq += d2 * d2
		upd_s += d3
		upd_sq += d3 * d3
		del_s += d4
		del_sq += d4 * d4

		bar.Increment()
		time.Sleep(time.Second)
	}
	bar.Finish()
	fmt.Println("General CRUD")
	fmt.Println(meanVar(delta, square, count))
	fmt.Println("CREATE")
	fmt.Println(meanVar(create_s, create_sq, count))
	fmt.Println("READ")
	fmt.Println(meanVar(read_s, read_sq, count))
	fmt.Println("UPDATE")
	fmt.Println(meanVar(upd_s, upd_sq, count))
	fmt.Println("DELETE")
	fmt.Println(meanVar(del_s, del_sq, count))
}
