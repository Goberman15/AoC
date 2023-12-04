package main

import (
	"fmt"
	"io"
	"log"
	"math"
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

	totalCards := len(inputSlice)
	totalPoint := 0
	totalCardInstances := 0
	cardInstances := make([]int, totalCards)

	for i, v := range inputSlice {
		cardInstances[i]++
		cardValues := strings.Split(v, ": ")[1]
		point, totalMatch := getTotalPoint(cardValues)
		fmt.Println(point, totalMatch)
		totalPoint += point

		if totalMatch > 0 {
			for x := i + 1; x <= i+totalMatch; x++ {
				if x >= totalCards {
					break
				}
				cardInstances[x] += cardInstances[i]
			}
		}

		totalCardInstances += cardInstances[i]
	}

	fmt.Printf("Total Point from Cards is %d\n", totalPoint)
	fmt.Printf("Total Cards Instances is %d\n", totalCardInstances)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/4/input", nil)

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

func getTotalPoint(line string) (int, int) {
	splittedLine := strings.Split(line, " | ")
	totalFoundWinningNumber := 0

	winningNumber := extractWinningNumber(splittedLine[0])

	for _, v := range strings.Split(splittedLine[1], " ") {
		if v != "" {
			_, ok := winningNumber[v]

			if ok {
				totalFoundWinningNumber++
			}
		}
	}

	if totalFoundWinningNumber == 0 {
		return 0, 0
	}

	return int(math.Pow(2, float64(totalFoundWinningNumber-1))), totalFoundWinningNumber
}

func extractWinningNumber(cardNum string) map[string]bool {
	splittedNum := strings.Split(cardNum, " ")

	winningNumber := make(map[string]bool)

	for _, v := range splittedNum {
		if v != "" {
			winningNumber[v] = true
		}
	}

	return winningNumber
}
