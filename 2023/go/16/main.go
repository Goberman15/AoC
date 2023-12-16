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

	tiles := make([][]string, 0, len(inputSlice))

	for _, v := range inputSlice {
		tiles = append(tiles, strings.Split(v, ""))
	}

	energyTiles := fillBeamEnergy(tiles, [2]int{0, 0}, ">")

	totalEnergy := countEnergy(energyTiles)

	fmt.Println("Total Starting Energized Tiles", totalEnergy)

	max := 0

	for i := 0; i < len(tiles); i++ {
		// top
		topStartTiles := fillBeamEnergy(tiles, [2]int{i, 0}, "v")
		topEnergy := countEnergy(topStartTiles)

		if topEnergy > max {
			max = topEnergy
		}

		// left
		leftStartTiles := fillBeamEnergy(tiles, [2]int{0, i}, ">")
		leftEnergy := countEnergy(leftStartTiles)

		if leftEnergy > max {
			max = leftEnergy
		}

		// bottom
		bottomStartTiles := fillBeamEnergy(tiles, [2]int{i, len(tiles) - 1}, "^")
		bottomEnergy := countEnergy(bottomStartTiles)

		if bottomEnergy > max {
			max = bottomEnergy
		}

		// right
		rightStartTiles := fillBeamEnergy(tiles, [2]int{len(tiles) - 1, i}, "<")
		rightEnergy := countEnergy(rightStartTiles)

		if rightEnergy > max {
			max = rightEnergy
		}
	}

	fmt.Println("Max Energized Tiles", max)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/16/input", nil)

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

func fillBeamEnergy(tiles [][]string, startCoord [2]int, startDir string) [][]string {
	energyTiles := make([][]string, len(tiles))
	pathTiles := make([][]string, len(tiles))
	for _, line := range tiles {
		for i := 0; i < len(line); i++ {
			energyTiles[i] = append(energyTiles[i], ".")
			pathTiles[i] = append(pathTiles[i], ".")
		}
	}
	coords := make([][2]int, 0)
	coords = append(coords, startCoord)
	dirs := []string{startDir}
	for {
		if len(coords) == 0 {
			break
		}

		deadIdx := make(map[int]bool)

		for idx, coord := range coords {
			x := coord[0]
			y := coord[1]

			if x < 0 || y < 0 || x >= len(tiles[0]) || y >= len(tiles) {
				deadIdx[idx] = true
				continue
			}

			dir := dirs[idx]
			tile := tiles[y][x]

			energyTiles[y][x] = "#"
			if tiles[y][x] == "." {
				pathTiles[y][x] = dir
			} else {
				pathTiles[y][x] = tiles[y][x]
			}

			switch tile {
			case ".":
				if dir == ">" {
					coords[idx] = [2]int{x + 1, y}
				} else if dir == "<" {
					coords[idx] = [2]int{x - 1, y}
				} else if dir == "^" {
					coords[idx] = [2]int{x, y - 1}
				} else if dir == "v" {
					coords[idx] = [2]int{x, y + 1}
				}
			case "/":
				if dir == ">" {
					if y > 0 && pathTiles[y-1][x] == "^" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x, y - 1}
					dirs[idx] = "^"
				} else if dir == "<" {
					if y < len(pathTiles)-1 && pathTiles[y+1][x] == "v" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x, y + 1}
					dirs[idx] = "v"
				} else if dir == "^" {
					if x < len(pathTiles[0])-1 && pathTiles[y][x+1] == ">" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x + 1, y}
					dirs[idx] = ">"
				} else if dir == "v" {
					if x > 0 && pathTiles[y][x-1] == "<" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x - 1, y}
					dirs[idx] = "<"
				}
			case "\\":
				if dir == ">" {
					if y < len(pathTiles)-1 && pathTiles[y+1][x] == "v" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x, y + 1}
					dirs[idx] = "v"
				} else if dir == "<" {
					if y > 0 && pathTiles[y-1][x] == "^" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x, y - 1}
					dirs[idx] = "^"
				} else if dir == "^" {
					if x > 0 && pathTiles[y][x-1] == "<" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x - 1, y}
					dirs[idx] = "<"
				} else if dir == "v" {
					if x < len(pathTiles[0])-1 && pathTiles[y][x+1] == ">" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x + 1, y}
					dirs[idx] = ">"
				}
			case "-":
				if dir == ">" {
					if x < len(pathTiles[0])-1 && pathTiles[y][x+1] == ">" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x + 1, y}
				} else if dir == "<" {
					if x > 0 && pathTiles[y][x-1] == "<" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x - 1, y}
				} else if dir == "^" || dir == "v" {
					if x < len(pathTiles[0])-1 && pathTiles[y][x+1] == ">" {
						deadIdx[idx] = true
						continue
					} else {
						coords[idx] = [2]int{x + 1, y}
						dirs[idx] = ">"
					}
					if x > 0 && pathTiles[y][x-1] == "<" {
						deadIdx[idx] = true
						continue
					} else {
						coords = append(coords, [2]int{x - 1, y})
						dirs = append(dirs, "<")
					}
				}
			case "|":
				if dir == ">" || dir == "<" {
					if y < len(pathTiles)-1 && pathTiles[y+1][x] == "v" {
						deadIdx[idx] = true
						continue
					} else {
						coords[idx] = [2]int{x, y + 1}
						dirs[idx] = "v"
					}
					if y > 0 && pathTiles[y-1][x] == "^" {
						deadIdx[idx] = true
						continue
					} else {
						coords = append(coords, [2]int{x, y - 1})
						dirs = append(dirs, "^")
					}
				} else if dir == "^" {
					if y > 0 && pathTiles[y-1][x] == "^" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x, y - 1}
				} else if dir == "v" {
					if y < len(pathTiles)-1 && pathTiles[y+1][x] == "v" {
						deadIdx[idx] = true
						continue
					}
					coords[idx] = [2]int{x, y + 1}
				}
			}
		}

		deadIdxSlice := make([]int, 0, len(deadIdx))
		for k := range deadIdx {
			deadIdxSlice = append(deadIdxSlice, k)
		}

		slices.Sort(deadIdxSlice)

		for i, k := range deadIdxSlice {
			dirs = splice(dirs, k-i)
			coords = splice(coords, k-i)
		}
	}

	return energyTiles
}

func splice[E comparable](slc []E, n int) []E {
	if n >= len(slc) || n < 0 {
		return slc
	}

	if n == 0 {
		return slc[1:]
	}

	res := make([]E, 0, len(slc)-1)
	res = append(res, slc[0:n]...)
	res = append(res, slc[n+1:]...)
	return res
}

func countEnergy(energyTiles [][]string) int {
	total := 0

	for _, line := range energyTiles {
		for _, v := range line {
			if v == "#" {
				total++
			}
		}
	}

	return total
}
