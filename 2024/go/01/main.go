package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	// "slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	input := getInput()

	rows := strings.Split(input, "\n")
	rows = rows[:len(rows)-1]

	group1, group2 := make([]int, 0, len(rows)), make([]int, 0, len(rows))

	for _, row := range rows {
		splitted := strings.Split(row, "   ")
		str1, str2 := splitted[0], splitted[1]

		num1, err := strconv.Atoi(str1)
		if err != err {
			log.Fatal("Fail to Convert String to Integer")
		}

		num2, err := strconv.Atoi(str2)
		if err != err {
			log.Fatal("Fail to Convert String to Integer")
		}

		group1 = append(group1, num1)
		group2 = append(group2, num2)
	}

	// slices.Sort(group1)
	// slices.Sort(group2)

	// sum := calcDiff(group1, group2)

	// fmt.Println(sum)

	nummap := make(map[int]bool)

	for _, num :=  range group1 {
		nummap[num] = true
	}

	sum2 := 0

	for _, num := range group2 {
		if nummap[num] {
			sum2 += num
		}
	}

	fmt.Println(sum2)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2024/day/1/input", nil)

	if err != nil {
		panic(err)
	}

	cookie := os.Getenv("AOC_COOKIE")

	req.Header.Set("cookie", cookie)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic("Fail to Retrieve Input")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
	}

	return string(body)
}

func calcDiff(group1, group2 []int) int {
	sum := 0

	for idx, num1 := range group1 {
		num2 := group2[idx]

		diff := float64(num2 - num1)
		diff = math.Abs(diff)

		sum += int(diff)
	}

	return sum
}
