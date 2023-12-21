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

	startX, startY := 0, 0
	fieldMap := make([][]string, 0, len(inputSlice))

	for j, line := range inputSlice {
		slc := strings.Split(line, "")
		fieldMap = append(fieldMap, slc)
		for i, v := range slc {
			if v == "S" {
				startX, startY = i, j
			}
		}
	}

	count := countFarthest(fieldMap, startX, startY, 64)

	fmt.Println(count)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/21/input", nil)

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

func countFarthest(fieldMap [][]string, startX, startY, steps int) int {
	cur := make([][2]int, 0)
	cur = append(cur, [2]int{startX, startY})
	ops := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for x := 0; x < steps; x++ {
		fmt.Println(x)
		newReach := make([][2]int, 0)
		for _, field := range cur {
			for _, op := range ops {
				x := field[0] + op[0]
				y := field[1] + op[1]

				if x < 0 || y < 0 || x >= len(fieldMap[0]) || y >= len(fieldMap) || fieldMap[y][x] == "#" {
					continue
				}

				if fieldMap[y][x] != "O" {
					fieldMap[y][x] = "O"
					newReach = append(newReach, [2]int{x, y})
				}
			}

			fieldMap[field[1]][field[0]] = "."
		}

		cur = newReach
	}

	return len(cur)
}
