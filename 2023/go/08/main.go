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

	directions := strings.Split(inputSlice[0], "")
	mapSlice := inputSlice[2:]
	// start := strings.Split(mapSlice[0], " = ")[0]
	// curr := start
	// end := strings.Split(mapSlice[len(mapSlice)-1], " = ")[0]

	curr := "AAA"
	end := "ZZZ"

	nodeMap := buildNodeMap(mapSlice)

	step, idx := 0, 0

	for {
		dir := directions[idx]
		step++
		idx++

		if dir == "R" {
			curr = nodeMap[curr][1]
		} else if dir == "L" {
			curr = nodeMap[curr][0]
		}

		if curr == end {
			break
		}

		if idx == len(directions) {
			idx = 0
		}
	}

	fmt.Printf("Steps Required to Reach End is %d\n", step)

	startPoint := make([]string, 0)

	for k := range nodeMap {
		if string(k[2]) == "A" {
			startPoint = append(startPoint, k)
		}
	}

	steps := make([]int, 0, len(startPoint))

	for _, v := range startPoint {
		step2, idx2 := 0, 0
		nd := v

		for {
			dir := directions[idx2]
			step2++
			idx2++

			if dir == "R" {
				nd = nodeMap[nd][1]
			} else if dir == "L" {
				nd = nodeMap[nd][0]
			}

			if string(nd[2]) == "Z" {
				steps = append(steps, step2)
				break
			}

			if idx2 == len(directions) {
				idx2 = 0
			}
		}

	}

	lcm := steps[0]

	for x := 1; x < len(steps); x++ {
		lcm = (lcm * steps[x]) / gcd(lcm, steps[x])
	}

	fmt.Printf("Steps Required to Reach Only Nodes that End with Z is %d\n", lcm)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/8/input", nil)

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

func buildNodeMap(mapSlc []string) map[string][]string {
	nodes := make(map[string][]string)

	for _, v := range mapSlc {
		node := v[0:3]
		nodeDirs := []string{v[7:10], v[12:15]}

		nodes[node] = nodeDirs
	}

	return nodes
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}

	return gcd(b, a % b)
}
