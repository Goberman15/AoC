package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type cardEvaluation struct {
	hand      string
	bid       int
	handType  string
	handScore int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	input := getInput()
	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]

	cardEvals := make([]*cardEvaluation, 0, len(inputSlice))

	for _, v := range inputSlice {
		splitted := strings.Split(v, " ")
		hand, bid := splitted[0], splitted[1]

		card := evaluateCard(hand, bid)

		cardEvals = append(cardEvals, card)
	}

	cardVal1 := map[string]int{
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"T": 10,
		"J": 11,
		"Q": 12,
		"K": 13,
		"A": 14,
	}

	sortHelper := sortHelperBuilder(cardVal1)

	slices.SortStableFunc(cardEvals, sortHelper)

	totalScore := 0

	for i, c := range cardEvals {
		score := (i + 1) * c.bid
		fmt.Printf("rank %d. %+v, total winnings: %d\n", i+1, c, score)
		totalScore += score
	}

	fmt.Printf("The Total Winnings are %d\n", totalScore)

	//--------Part2
	cardEvals2 := make([]*cardEvaluation, 0, len(inputSlice))

	for _, v := range inputSlice {
		splitted := strings.Split(v, " ")
		hand, bid := splitted[0], splitted[1]

		card := evaluateCard2(hand, bid)

		cardEvals2 = append(cardEvals2, card)
	}

	cardVal2 := map[string]int{
		"J": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"T": 10,
		"Q": 12,
		"K": 13,
		"A": 14,
	}

	sortHelper2 := sortHelperBuilder(cardVal2)

	slices.SortStableFunc(cardEvals2, sortHelper2)

	totalScore2 := 0

	for i, c := range cardEvals2 {
		score := (i + 1) * c.bid
		fmt.Printf("rank %d. %+v, total winnings: %d\n", i+1, c, score)
		totalScore2 += score
	}

	fmt.Printf("The Total Winnings of Part 2 are %d\n", totalScore2)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/7/input", nil)

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

func sortHelperBuilder(cardVal map[string]int) func(a, b *cardEvaluation) int {
	return func(a, b *cardEvaluation) int {
		comp := cmp.Compare(a.handScore, b.handScore)

		if comp == 0 {
			for x := 0; x < len(a.hand); x++ {
				valA := cardVal[string(a.hand[x])]
				valB := cardVal[string(b.hand[x])]

				if valA > valB {
					return 1
				} else if valB > valA {
					return -1
				}
			}
			return 0
		} else {
			return comp
		}
	}

}

func evaluateCard(hand string, bid string) *cardEvaluation {
	bidInt, err := strconv.Atoi(bid)
	if err != nil {
		log.Fatal(err)
	}

	cardEval := &cardEvaluation{
		hand: hand,
		bid:  bidInt,
	}

	handMap := make(map[string]int)

	for _, c := range strings.Split(hand, "") {
		_, ok := handMap[c]
		if !ok {
			handMap[c] = 1
		} else {
			handMap[c]++
		}
	}

	switch len(handMap) {
	case 1:
		cardEval.handType = "Five of A Kind"
		cardEval.handScore = 7
	case 2:
		for _, v := range handMap {
			if v == 4 || v == 1 {
				cardEval.handType = "Four of A Kind"
				cardEval.handScore = 6
				break
			} else if v == 3 || v == 2 {
				cardEval.handType = "Full House"
				cardEval.handScore = 5
				break
			}
		}
	case 3:
		for _, v := range handMap {
			if v == 3 {
				cardEval.handType = "Three of A Kind"
				cardEval.handScore = 4
				break
			} else if v == 2 {
				cardEval.handType = "Two Pair"
				cardEval.handScore = 3
				break
			}
		}
	case 4:
		cardEval.handType = "One Pair"
		cardEval.handScore = 2
	case 5:
		cardEval.handType = "High Card"
		cardEval.handScore = 1
	}

	return cardEval
}

func evaluateCard2(hand string, bid string) *cardEvaluation {
	jokerCount := 0

	bidInt, err := strconv.Atoi(bid)
	if err != nil {
		log.Fatal(err)
	}

	cardEval := &cardEvaluation{
		hand: hand,
		bid:  bidInt,
	}

	handMap := make(map[string]int)

	for _, c := range strings.Split(hand, "") {
		_, ok := handMap[c]
		if !ok {
			handMap[c] = 1
		} else {
			handMap[c]++
		}

		if c == "J" {
			jokerCount++
		}
	}

	switch len(handMap) {
	case 1:
		cardEval.handType = "Five of A Kind"
		cardEval.handScore = 7
	case 2:
		for _, v := range handMap {
			if jokerCount > 0 {
				cardEval.handType = "Five of A Kind"
				cardEval.handScore = 7
				break
			}

			if v == 4 || v == 1 {
				cardEval.handType = "Four of A Kind"
				cardEval.handScore = 6
				break
			} else if v == 3 || v == 2 {
				cardEval.handType = "Full House"
				cardEval.handScore = 5
				break
			}
		}
	case 3:
		for _, v := range handMap {
			if v == 3 {
				if jokerCount > 0 {
					cardEval.handType = "Four of A Kind"
					cardEval.handScore = 6
					break
				}
				cardEval.handType = "Three of A Kind"
				cardEval.handScore = 4
				break
			} else if v == 2 {
				if jokerCount == 2 {
					cardEval.handType = "Four of A Kind"
					cardEval.handScore = 6
					break
				} else if jokerCount == 1 {
					cardEval.handType = "Full House"
					cardEval.handScore = 5
					break
				}
				cardEval.handType = "Two Pair"
				cardEval.handScore = 3
				break
			}
		}
	case 4:
		if jokerCount > 0 {
			cardEval.handType = "Three of A Kind"
			cardEval.handScore = 4
			break
		}
		cardEval.handType = "One Pair"
		cardEval.handScore = 2
	case 5:
		if jokerCount > 0 {
			cardEval.handType = "One Pair"
			cardEval.handScore = 2
			break
		}
		cardEval.handType = "High Card"
		cardEval.handScore = 1
	}

	return cardEval
}
