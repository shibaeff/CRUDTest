package main

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

const (
	baseURL = "https://function-ckn5zrqpjsuw4.azurewebsites.net/api/httptrigger"
)

func meanVar(delta, squares, count int64) string {
	avg_delta := float64(delta) / float64(count)
	avg_square := float64(squares) / float64(count)
	variance := avg_square - avg_delta*avg_delta
	mean := avg_delta
	return fmt.Sprintf("mean %f deviation %f", mean, math.Sqrt(variance))
}

func main() {
	s, sq := int64(0), int64(0)
	count := int64(10)
	for i := int64(0); i < count; i++ {
		start := time.Now()
		http.Get(baseURL)
		d := time.Now().Sub(start).Microseconds()
		s += d
		sq += d * d
	}
	fmt.Println("PING")
	fmt.Println(meanVar(s, sq, count))
}
