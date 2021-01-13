package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/araddon/dateparse"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parse(path string) (mean, std float64) {
	mean, std = 0.0, 0.0
	lines, err := readLines(path)
	if err != nil {
		log.Fatal(err)
	}
	sum := int64(0)
	squares := int64(0)
	for i := 0; i < len(lines)/2; i++ {
		d1, err := dateparse.ParseAny(strings.Join(strings.Split(lines[2*i], " ")[4:], " "))
		if err != nil {
			log.Fatal(err)
		}

		d2, err := dateparse.ParseAny(strings.Join(strings.Split(lines[2*i+1], " ")[4:], " "))
		if err != nil {
			log.Fatal(err)
		}
		dur := d2.Sub(d1)
		sum += dur.Microseconds()
		squares += dur.Microseconds() * dur.Microseconds()
	}
	mean = float64(sum) / float64(len(lines)/2)
	avg_sq := float64(squares) / float64(len(lines)/2)
	std = math.Sqrt(avg_sq - mean*mean)
	return
}
func main() {
	fmt.Println(parse("./logs/create.log"))
	fmt.Println(parse("./logs/read.log"))
	fmt.Println(parse("./logs/upd.log"))
	fmt.Println(parse("./logs/del.log"))
}
