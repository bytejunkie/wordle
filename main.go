package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// how many guesses taken
	numberOfGuesses := 0

	rand.Seed(time.Now().UnixNano())

	answer := getAnswer()
	fmt.Println(answer)

	reader := bufio.NewReader(os.Stdin)

	for {
		if numberOfGuesses == 0 {
			fmt.Println("\nHello, try to guess my five letter word.")
		} else {
			fmt.Println("\nYou've made ", numberOfGuesses, "guess, keep trying!")
		}
		numberOfGuesses++

		guess, _ := reader.ReadString('\n')
		// convert CRLF to LF
		guess = strings.Replace(guess, "\n", "", -1)

		// TODO fix up this bit where it calls the same routine twice.
		if strings.Compare(answer, guess) == 0 {
			checkAnswer(answer, guess)
			fmt.Println("\nYou got it in", numberOfGuesses, "guesses! ")
			break
		} else {
			checkAnswer(answer, guess)
		}

	}

}

func getAnswer() string {
	// randomly pick a word from the csv file
	f, err := os.Open("words.csv")
	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for ", err)
	}

	word := records[rand.Intn(len(records))][0]

	return word
}

func checkAnswer(answer string, guess string) {
	// var Green = "\033[32m"
	var BGreen = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	var Reset = "\033[0m"
	var Yellow = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})

	if len(guess) != 5 {
		fmt.Println("your answer needs to be 5 chars long!")
		// break
	}
	fmt.Println("Checking your guess:")

	for i := 0; i < len(guess); i++ {
		if guess[i:i+1] == answer[i:i+1] {
			fmt.Printf(BGreen + guess[i:i+1] + Reset)
		} else {
			if strings.Contains(answer, guess[i:i+1]) {
				fmt.Printf(Yellow + guess[i:i+1] + Reset)
			} else {
				fmt.Printf(guess[i : i+1])
			}
		}
	}

}

//TODO - collect a list of all the letters used and display it.
//TODO - put some space into the results to make it more legible
