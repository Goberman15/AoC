package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	score := scoring(inputSlice)

	fmt.Printf("Total Score if everything goes exactly according to Strategy Guide: %d", score)
	
	score2 := scoring2(inputSlice)
	fmt.Printf("Total Score if everything goes exactly according to Strategy Guide: %d", score2)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2022/day/2/input", nil)

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

func scoring(strategies []string) int {
	totalScore := 0

	scoreMatrix := map[string]int{
		"A X": 4,
		"A Y": 8,
		"A Z": 3,
		"B X": 1,
		"B Y": 5,
		"B Z": 9,
		"C X": 7,
		"C Y": 2,
		"C Z": 6,
	}

	for _, v := range strategies {
		score, ok := scoreMatrix[v]

		if !ok {
			log.Fatalf("Wrong Startegy: %s", v)
		}

		totalScore += score
	}

	return totalScore
}

func scoring2(strategies []string) int {
	totalScore := 0

	scoreMatrix := map[string]int{
		"A X": 3,
		"A Y": 4,
		"A Z": 8,
		"B X": 1,
		"B Y": 5,
		"B Z": 9,
		"C X": 2,
		"C Y": 6,
		"C Z": 7,
	}

	for _, v := range strategies {
		score, ok := scoreMatrix[v]

		if !ok {
			log.Fatalf("Wrong Startegy: %s", v)
		}

		totalScore += score
	}

	return totalScore
}
