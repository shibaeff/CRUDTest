package main

import (
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	baseURL = ""
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
		Timeout: 3 * time.Second,
	}
)

func perFormRequest(err error, r *http.Request) (int64, int64) {
	start := time.Now()
	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	d := time.Now().Sub(start)
	type duration struct {
		Dur int64 `json:"dur"`
	}
	var dur duration
	if err = json.Unmarshal(bodyBytes, &dur); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	// i, err := strconv.Atoi(dur.Dur)
	if err != nil {
		log.Fatal(err)
	}
	return d.Microseconds(), dur.Dur
}

func sendCreate(postfix string, id int64) (d int64, i int64) {
	url := baseURL + postfix
	person := User{
		//FirstName: "test",
		//LastName:  "test",
		//UserName:  "test",
		Id: id,
	}
	json_str, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.NewRequest(http.MethodGet, url, strings.NewReader(string(json_str)))
	if err != nil {
		log.Fatal(err)
	}
	d, i = perFormRequest(err, r)
	return
}

func sendRead(postfix string, id int64) (d int64, i int64) {
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
	r, err := http.NewRequest(http.MethodGet, url, strings.NewReader(string(json_str)))
	if err != nil {
		log.Fatal(err)
	}
	d, i = perFormRequest(err, r)
	return
}

func sendUpd(postfix string, id int64) (d int64, i int64) {
	url := baseURL + postfix
	person := User{
		//FirstName: "test1",
		//LastName:  "test1",
		UserName: "test1",
		Id:       id,
	}
	json_str, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.NewRequest(http.MethodGet, url, strings.NewReader(string(json_str)))
	d, i = perFormRequest(err, r)
	return
}

func sendDelete(postfix string, id int64) (d int64, i int64) {
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
	r, err := http.NewRequest(http.MethodGet, url, strings.NewReader(string(json_str)))
	d, i = perFormRequest(err, r)
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
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	baseURL = os.Getenv("BASE")
	count := int64(50)
	delta := int64(0)
	square := int64(0)
	bar := pb.StartNew(int(count))
	var (
		create_s, create_sq, read_s, read_sq, upd_s, upd_sq, del_s, del_sq int64
	)
	// var i1, i2, i3, i4 int
	var (
		i1_s, i2_s, i3_s, i4_s, i1_k, i2_k, i3_k, i4_k int64
	)
	for i := int64(0); i < int64(count); i++ {
		d1, i1 := sendCreate("/create", i)
		i1_s += i1
		i1_k += i1 * i1
		d2, i2 := sendRead("/read", i)
		i2_s += i2
		i2_k += i2 * i2
		d3, i3 := sendCreate("/update", i)
		i3_s += i3
		i3_k += i3 * i3
		d4, i4 := sendDelete("/delete", i)
		i4_s += i4
		i4_k += i4 * i4
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
		// time.Sleep(time.Second)
	}
	bar.Finish()
	fmt.Println("General CRUD")
	fmt.Println(meanVar(delta, square, count))

	fmt.Println("CREATE")
	fmt.Println(meanVar(create_s, create_sq, count))

	fmt.Println("CREATE Internal")
	fmt.Println(meanVar(i1_s, i1_k, count))

	fmt.Println("READ")
	fmt.Println(meanVar(read_s, read_sq, count))

	fmt.Println("READ Internal")
	fmt.Println(meanVar(i2_s, i2_k, count))

	fmt.Println("UPDATE")
	fmt.Println(meanVar(upd_s, upd_sq, count))

	fmt.Println("UPDATE Internal")
	fmt.Println(meanVar(i3_s, i3_k, count))

	fmt.Println("DELETE")
	fmt.Println(meanVar(del_s, del_sq, count))

	fmt.Println("DELETE Internal")
	fmt.Println(meanVar(i4_s, i4_k, count))
}
