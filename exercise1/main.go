package main

import (
	"fmt"
	"encoding/csv"
	"os"
	"flag"
	"strings"
	"strconv"
)

const (
	defaultName = "./test_input.csv"
	fileSuffix = ".csv"
)

var (
	userInput = flag.Bool("f", false, "Enable user input file")
	filename string = defaultName
)

func readUserFile()  {
	fmt.Println("Enter the file you want to use: ")
	fmt.Scanln(&filename)

	if !strings.HasSuffix(filename, fileSuffix) {
		filename = defaultName // use the default
		fmt.Println("Filename does not end with .csv")
		fmt.Println("Using the default file:", filename)
	}
}


func main() {
	flag.Parse()	
	if *userInput == true {
		readUserFile()
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
	for _, line := range(records) {
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
			fmt.Println("CORRECT")
		}
	}

	fmt.Println(numCorrect, "out of", numQuestions)
}
