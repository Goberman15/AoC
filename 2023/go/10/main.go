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

	pipes := make([][]string, 0, len(inputSlice))

	for _, e := range inputSlice {
		pipes = append(pipes, strings.Split(e, ""))
	}

	step, visited := findFarthestPoint(pipes)

	fmt.Println("step", step)

	internal := countNestField(pipes, visited)

	fmt.Println("Nested", internal)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/10/input", nil)

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

func findStartingPoint(pipes [][]string) (int, int) {
	for y, row := range pipes {
		for x, v := range row {
			if v == "S" {
				return x, y
			}
		}
	}
	return 0, 0
}

func findFarthestPoint(pipes [][]string) (int, map[string]bool) {
	startY, startX := findStartingPoint(pipes)
	visited := make(map[string]bool)
	y, x := startX, startY
	// pipe := []string{"-", "|", "F", "L", "J", "7"}
	step := 0

	for {
		currField := pipes[y][x]
		step++
		visited[fmt.Sprintf("%d-%d", y, x)] = true
		if currField == "S" {
			// check top
			if y > 0 {
				if pipes[y-1][x] == "|" || pipes[y-1][x] == "F" || pipes[y-1][x] == "7" {
					y--
					continue
				}
			}

			// check bottom
			if y < len(pipes)-1 {
				if pipes[y+1][x] == "|" || pipes[y+1][x] == "J" || pipes[y+1][x] == "L" {
					y++
					continue
				}
			}

			// check left
			if x > 0 {
				if pipes[y][x-1] == "-" || pipes[y][x-1] == "F" || pipes[y][x-1] == "L" {
					x--
					continue
				}
			}

			// check right
			if x < len(pipes[y])-1 {
				if pipes[y][x+1] == "-" || pipes[y][x+1] == "J" || pipes[y][x+1] == "7" {
					x++
					continue
				}
			}
		} else if currField == "-" {
			left := fmt.Sprintf("%d-%d", y, x-1)
			if visited[left] && pipes[y][x-1] != "S" {
				x++
			} else {
				x--
			}
		} else if currField == "|" {
			top := fmt.Sprintf("%d-%d", y-1, x)
			if visited[top] && pipes[y-1][x] != "S" {
				y++
			} else {
				y--
			}
		} else if currField == "L" {
			right := fmt.Sprintf("%d-%d", y, x+1)
			if visited[right] && pipes[y][x+1] != "s" {
				y--
			} else {
				x++
			}
		} else if currField == "J" {
			left := fmt.Sprintf("%d-%d", y, x-1)
			if visited[left] && pipes[y][x-1] != "S" {
				y--
			} else {
				x--
			}
		} else if currField == "7" {
			left := fmt.Sprintf("%d-%d", y, x-1)
			if visited[left] && pipes[y][x-1] != "S" {
				y++
			} else {
				x--
			}
		} else if currField == "F" {
			right := fmt.Sprintf("%d-%d", y, x+1)
			if visited[right] && pipes[y][x+1] != "S" {
				y++
			} else {
				x++
			}
		}

		if pipes[y][x] == "S" {
			break
		}
	}

	return step/2, visited
}

func countNestField(pipes [][]string, visited map[string]bool) int {
	internal := 0
	// choose all down or all up, because on one side must there must be odd number of pipe to consider a point as inside
	pipe := []string{"|", "F", "7"}

	for y, row := range pipes {
		// first and last row can't be part of internal field
		if y == 0 || y >= len(pipes) {
			continue
		}
		for x := range row {
			// first and last column can't be part of internal field
			if x == 0 || x >= len(row) {
				continue
			}

			currStr := fmt.Sprintf("%d-%d", y, x)
			rowFence := 0

			if visited[currStr] {
				continue
			}

			for i := x; i < len(row); i++ {
				if slices.Contains(pipe, row[i]) && visited[fmt.Sprintf("%d-%d", y, i)] {
					if row[i] != "S" {
						rowFence++
						continue
					}

					// if it S, check if S has connetion to down, if yes count it
					// if choose to check for up at line 179, check for up here, not down
					if row[i] == "S" && visited[fmt.Sprintf("%d-%d", y-1, i)] {
						rowFence++
					}
				}
			}

			if rowFence%2 != 0 {
				internal++
			}
		}
	}

	return internal
}
