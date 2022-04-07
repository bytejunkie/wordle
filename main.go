package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/inancgumus/screen"
)

type Stats struct {
	Stats []Stat `json:"stats"`
}

type Stat struct {
	TimesPlayed   int   `json:"timesPlayed"`
	LastPlayed    int   `json:"lastPlayed"`
	CurrentStreak int   `json:"currentStreak"`
	MaxStreak     int   `json:"maxStreak"`
	Tries         []int `json:"tries"`
	Tries0        int   `json:"tries0"`
	Tries1        int   `json:"tries1"`
	Tries2        int   `json:"tries2"`
	Tries3        int   `json:"tries3"`
	Tries4        int   `json:"tries4"`
	Tries5        int   `json:"tries5"`
	Tries6        int   `json:"tries6"`
}

func main() {
	// how many guesses taken
	numberOfGuesses := 0
	lettersUsed := make([]string, 0)
	// lettersUsed := ""
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

			updateStats(numberOfGuesses)
			break
		} else if numberOfGuesses == 6 {
			checkAnswer(answer, guess)
			updateStats(0)
			break
		} else {
			checkAnswer(answer, guess)
			lettersUsed = printLettersUsed(lettersUsed, strings.Split(guess, ""))
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
		return
	}
	screen.Clear()
	fmt.Println("Checking your guess:")

	for i := 0; i < len(guess); i++ {
		if guess[i:i+1] == answer[i:i+1] {
			fmt.Printf(BGreen + guess[i:i+1] + Reset + " ")
		} else {
			if strings.Contains(answer, guess[i:i+1]) {
				fmt.Printf(Yellow + guess[i:i+1] + Reset + " ")
			} else {
				fmt.Printf(guess[i:i+1] + " ")
			}
		}
	}
}

func updateStats(numberOfGuesses int) {

	_, err := os.Stat("stats.json")
	if err == nil {
		// fmt.Println("File Exists")
	}
	if errors.Is(err, os.ErrNotExist) {
		source, err := os.Open("_stats.json")

		destination, err := os.Create("stats.json")
		if err != nil {
			fmt.Println(err)
		}
		defer destination.Close()
		newStatsFile, err := io.Copy(destination, source)
		source.Close()
		fmt.Println(newStatsFile)
	}

	statsFile, err := os.Open("stats.json")
	byteValue, _ := ioutil.ReadAll(statsFile)

	var stats Stats
	json.Unmarshal(byteValue, &stats)

	// start updating the stats
	// number of times played in total
	stats.Stats[0].TimesPlayed++

	// did we get it right?
	if numberOfGuesses != 0 {
		fmt.Println("\nYou got it in", numberOfGuesses, "guesses! ")
		stats.Stats[0].CurrentStreak++
	} else {
		stats.Stats[0].CurrentStreak = 0
	}

	// are we currently on a max streak?
	if numberOfGuesses != 0 && stats.Stats[0].CurrentStreak > stats.Stats[0].MaxStreak {
		stats.Stats[0].MaxStreak++
	}

	stats.Stats[0].Tries[numberOfGuesses]++

	statsByte, err := json.Marshal(stats)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("stats.json", statsByte, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Statistics")
	fmt.Printf("Played %d games == %d %% Win == Current Streak %d == Max Streak %d \n",
		stats.Stats[0].TimesPlayed,
		((stats.Stats[0].TimesPlayed-stats.Stats[0].Tries0)/stats.Stats[0].TimesPlayed)*100,
		stats.Stats[0].CurrentStreak,
		stats.Stats[0].MaxStreak)

	fmt.Println("Win Distribution")
	for i := 1; i < 7; i++ {
		fmt.Printf("%d: %d\n", i, stats.Stats[0].Tries[i])
	}
}

func printLettersUsed(lettersUsed []string, guess []string) []string {
	for i := 0; i < len(guess); i++ {
		lettersUsed = append(lettersUsed, guess[i])
	}
	fmt.Println("\nLetters you have used:\n", lettersUsed)
	return lettersUsed
}
