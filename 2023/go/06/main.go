package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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
	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]

	fmt.Printf("%q\n", inputSlice)

	parsed := parseInput(inputSlice)
	time := parsed[0]
	distance := parsed[1]

	multiply := 1

	for i := 0; i < len(time); i++ {
		res := countWayToWin(time[i], distance[i])

		multiply *= res
	}

	fmt.Printf("Multiply of Way to Win is %d\n", multiply)

	parsed2 := parseInput2(inputSlice)
	time2, distance2 := parsed2[0], parsed2[1]

	realRes := countWayToWin(time2, distance2)
	fmt.Printf("Real Total Way to Win is %d\n", realRes)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/6/input", nil)

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

func parseInput(input []string) [][]int {
	parsed := make([][]int, 0)
	for _, v := range input {
		re := regexp.MustCompile(`\s+`)
		stripped := re.ReplaceAllString(v, " ")
		splitted := strings.Split(stripped, " ")

		temp := make([]int, 0, len(splitted)-1)

		for _, s := range splitted[1:] {
			num, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			temp = append(temp, num)
		}
		parsed = append(parsed, temp)
	}

	return parsed
}

func parseInput2(input []string) []int {
	parsed := make([]int, 0)
	for _, v := range input {
		re := regexp.MustCompile(`\s+`)
		stripped := re.ReplaceAllString(v, "")
		splitted := strings.Split(stripped, ":")

		num, err := strconv.Atoi(splitted[1])
		if err != nil {
			log.Fatal(err)
		}

		parsed = append(parsed, num)
	}

	return parsed
}

func countWayToWin(time int, dist int) int {
	total := 0

	for x := 1; x <= time; x++ {
		reach := x * (time - x)

		if reach > dist {
			total++
		}
	}

	return total
}
