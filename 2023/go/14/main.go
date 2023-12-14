package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
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

	rockMap := parseRockMap(inputSlice)

	tiltedMap := tilt(rockMap)

	totalLoad := 0

	for idx, line := range tiltedMap {
		multiplier := len(tiltedMap) - idx
		load := calculateLoad(line)

		totalLoad += multiplier * load
	}
	fmt.Println("Total Load", totalLoad)

	cycle := 1_000_000_000
	cycleStart := 0
	cycleCount := 0
	isCycleFound := false
	loadList := make([]int, 0)
	cycleLoad := 0

	for x := 1; x <= cycle; x++ {
		for c := 0; c < 4; c++ {
			rockMap = tilt(rockMap)
			rockMap = rotate(rockMap)
		}

		total2 := 0
		for idx, line := range rockMap {
			multiplier := len(rockMap) - idx
			load := calculateLoad(line)

			total2 += multiplier * load
		}

		if x == 0 {
			loadList = append(loadList, total2)
			continue
		}

		if isCycleFound {
			cycleCount++
			if (cycleCount + 1) > len(loadList) {
				isCycleFound = false
				cycleStart = 0
				cycleCount = 0
				loadList = nil
				loadList = append(loadList, total2)
				continue
			}

			if total2 == loadList[cycleCount] {
				if cycleCount == len(loadList)-1 {
					y := (cycle - cycleStart) % len(loadList)
					cycleLoad = loadList[y]
					break
				}
			} else {
				isCycleFound = false
				cycleStart = 0
				cycleCount = 0
				loadList = nil
				loadList = append(loadList, total2)
				continue
			}
		} else {
			if slices.Contains(loadList, total2) {
				if len(loadList) > 1 {
					if loadList[0] == total2 {
						isCycleFound = true
						cycleStart = x
					} else {
						loadList = nil
						loadList = append(loadList, total2)
					}
				}
			} else {
				loadList = append(loadList, total2)
			}
		}
	}

	fmt.Println("Total Load after 1000000000 Cycles", cycleLoad)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/14/input", nil)

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

func parseRockMap(rockMap []string) [][]string {
	rockMapSlice := make([][]string, len(rockMap))

	for i, line := range rockMap {
		rockMapSlice[i] = strings.Split(line, "")
	}

	return rockMapSlice
}

func tilt(rockMap [][]string) [][]string {
	rockMapClone := make([][]string, len(rockMap))
	copy(rockMapClone, rockMap)

	for x := 0; x < len(rockMapClone[0]); x++ {
		roundRock := 0
		lastCubeRockIdx := -1
		for y := 0; y < len(rockMapClone); y++ {
			rock := rockMapClone[y][x]

			if rock == "#" {
				lastCubeRockIdx = y
				roundRock = 0
			} else if rock == "O" {
				roundRock++
				rockMapClone[y][x] = "."
				rockMapClone[lastCubeRockIdx+roundRock][x] = "O"
			}
		}
	}

	return rockMapClone
}

func rotate(rockMap [][]string) [][]string {
	rotated := make([][]string, 0, len(rockMap[0]))

	for x := 0; x < len(rockMap[0]); x++ {
		line := make([]string, 0, len(rockMap))
		for y := len(rockMap) - 1; y >= 0; y-- {
			line = append(line, rockMap[y][x])
		}
		rotated = append(rotated, line)
	}

	return rotated
}

func calculateLoad(line []string) int {
	load := 0

	for _, rock := range line {
		if rock == "O" {
			load++
		}
	}

	return load
}
