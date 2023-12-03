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

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	input := getInput()

	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]

	fmt.Println(len(inputSlice))

	out := 0
	power := 0

	for _, v := range inputSlice {
		id, cubeSets := getGameIdAndCubeSets(v)
		fmt.Println(id)

		possible := checkPossibility(cubeSets)
		fmt.Println(possible)

		if possible {
			out += id
		}

		power += extractMinimumPower(cubeSets)
	}

	fmt.Printf("The Sum of the IDs of the Possible Game is %d", out)
	fmt.Printf("The Sum of the Power of these Sets is %d", power)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/2/input", nil)

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

func getGameIdAndCubeSets(line string) (int, string) {
	splitted := strings.Split(line, ": ")
	header := splitted[0]
	cubeSets := splitted[1]
	idStr := (strings.Split(header, " "))[1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	return id, cubeSets
}

func checkPossibility(line string) bool {
	redLimit, greenLimit, blueLimit := 12, 13, 14

	rounds := strings.Split(line, "; ")

	for _, v := range rounds {
		cubeSet := strings.Split(v, ", ")

		for _, cube := range cubeSet {
			splitted := strings.Split(cube, " ")
			n, err := strconv.Atoi(splitted[0])
			if err != nil {
				log.Fatal(err)
			}
			
			color := splitted[1]
			fmt.Println(splitted)

			switch color {
			case "red":
				if n > redLimit {
					return false
				}
			case "green":
				if n > greenLimit {
					return false
				}
			case "blue":
				if n > blueLimit {
					return false
				}
			}
		}

	}

	return true
}

func extractMinimumPower(line string) int {
	redMin, greenMin, blueMin := 0, 0, 0

	rounds := strings.Split(line, "; ")

	for _, v := range rounds {
		cubeSet := strings.Split(v, ", ")

		for _, cube := range cubeSet {
			splitted := strings.Split(cube, " ")
			n, err := strconv.Atoi(splitted[0])
			if err != nil {
				log.Fatal(err)
			}
			
			color := splitted[1]
			fmt.Println(splitted)

			switch color {
			case "red":
				if n > redMin {
					redMin = n
				}
			case "green":
				if n > greenMin {
					greenMin = n
				}
			case "blue":
				if n > blueMin {
					blueMin = n
				}
			}
		}

	}

	return redMin * greenMin * blueMin
}
