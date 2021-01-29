package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultName          = "./test_input.csv"
	fileSuffix           = ".csv"
	defaultTimer float64 = 30.0
)

var (
	userInput           = flag.Bool("f", false, "Enable user input file")
	timer               = flag.Float64("t", defaultTimer, "Timer for the quiz (seconds)")
	shuffleFlag         = flag.Bool("s", false, "Shuffle the order of the parsed quiz")
	filename    string  = defaultName
	timerLength float64 = defaultTimer
)

func readUserFile() {
	fmt.Println("Enter the file you want to use: ")
	fmt.Scanln(&filename)

	if !strings.HasSuffix(filename, fileSuffix) {
		filename = defaultName // use the default
		fmt.Println("Filename does not end with .csv")
		fmt.Println("Using the default file:", filename)
	}
}

func shuffleQuestions(questions [][]string) [][]string {
	fmt.Println("Shuffling questions...")
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions
}

func main() {
	rand.Seed(time.Now().Unix())

	flag.Parse()
	if *userInput == true {
		readUserFile()
	}

	if *timer < 0 {
		fmt.Println("Timer flag was a non-positive value.")
		fmt.Println("Setting timer value to", defaultTimer, "seconds")

	} else {
		timerLength = defaultTimer
	}

	fmt.Println("Opening", filename, "...")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	numQuestions := len(records)
	numCorrect := 0

	if *shuffleFlag {
		records = shuffleQuestions(records)
	}

	processQuestions := func(ch chan<- struct{}) {
		for _, line := range records {
			question := line[0]
			answer, err := strconv.Atoi(line[1])
			if err != nil {
				fmt.Println("Error converting token to intger type")
			}
			var userAnswer int
			fmt.Printf("%v = ", question)
			fmt.Scanf("%d", &userAnswer)
			if userAnswer == answer {
				numCorrect++
			}
		}
		ch <- struct{}{}
	}

	signalChannel := make(chan struct{})
	defer close(signalChannel) // defer chan close

	delayTime := 30
	go processQuestions(signalChannel)

	timer := time.NewTimer(time.Duration(delayTime) * time.Second)
	defer timer.Stop() // handle this resource

	select {
	case <-signalChannel:
	case <-timer.C:
		fmt.Printf("\nTimeout after %ds\n", delayTime)
	}

	fmt.Println(numCorrect, "correct out of", numQuestions)

	// go run timer
	// select a timer message
	// otherwise process input
	// case timerMessage -> break out and print message

}
