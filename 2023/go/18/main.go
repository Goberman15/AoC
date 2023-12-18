package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type DigPlanItem struct {
	dir  string
	dist int
	hex  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	input := getInput()
	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]

	digPlan := make([]DigPlanItem, 0, len(inputSlice))

	for _, line := range inputSlice {
		splitted := strings.Split(line, " ")

		hex := splitted[2][1 : len(splitted[2])-1]

		digPlan = append(digPlan, DigPlanItem{
			dir:  splitted[0],
			dist: stringToInt(splitted[1]),
			hex:  hex,
		})
	}

	digMap := dugPerimeter(digPlan)
	totalArea := countTotalArea(digMap)

	fmt.Printf("The Lagoon could Hold %d CubicMeters of Lava\n", totalArea)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/18/input", nil)

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

func dugPerimeter(digPlan []DigPlanItem) [][]string {
	digMap := make([][]string, 1)
	curX, curY := 0, 0
	maxLength := 1
	digMap[curY] = append(digMap[curY], "#")

	for _, d := range digPlan {
		dir := d.dir

		for x := 0; x < d.dist; x++ {
			if dir == "R" {
				curX++
				if curX == len(digMap[curY]) {
					digMap[curY] = append(digMap[curY], "#")
					if curX == maxLength {
						maxLength++
					}
				} else {
					digMap[curY][curX] = "#"
				}
			} else if dir == "L" {
				curX--
				if curX < 0 {
					temp := make([]string, curX*-1)
					for idx, line := range digMap {
						digMap[idx] = append(temp, line...)
					}
					curX = 0
					maxLength++
				}
				digMap[curY][curX] = "#"
			} else if dir == "D" {
				curY++
				if curY == len(digMap) {
					temp := make([]string, curX+1)
					digMap = append(digMap, temp)
				} else {
					if curX >= len(digMap[curY]) {
						temp := make([]string, curX-(len(digMap[curY])-1))
						digMap[curY] = append(digMap[curY], temp...)
					}
				}
				digMap[curY][curX] = "#"
			} else if dir == "U" {
				curY--

				if curY < 0 {
					temp := make([][]string, 0, len(digMap)+1)
					tempLine := make([]string, curX+1)
					temp = append(temp, tempLine)
					digMap = append(temp, digMap...)
					curY = 0
				} else {
					if curX >= len(digMap[curY]) {
						temp := make([]string, curX-(len(digMap[curY])-1))
						digMap[curY] = append(digMap[curY], temp...)
					}
				}
				digMap[curY][curX] = "#"
			}
		}
	}

	return digMap
}

func countTotalArea(digMap [][]string) int {
	totalArea := 0

	for i, line := range digMap {
		fences := 0
		for j, v := range line {
			if v == "#" {
				totalArea++
				if i == 0 || i == len(digMap)-1 {
					continue
				}

				if len(digMap[i-1]) > j && digMap[i-1][j] == "#" {
					fences++
				}
				continue
			}

			if fences%2 == 1 {
				totalArea++
			}
		}
	}

	return totalArea
}

func stringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func printDigMap(digMap [][]string) {
	for _, line := range digMap {
		for _, v := range line {
			if v == "" {
				fmt.Print(".")
			} else {
				fmt.Print(v)
			}
		}
		fmt.Println()
	}
}
